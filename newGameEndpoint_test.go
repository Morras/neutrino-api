package neutrinoapi_test

import (
	"errors"
	api "github.com/morras/neutrinoapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strconv"
)

type AcceptingRequestParser struct {
	request *http.Request
}

func (rp *AcceptingRequestParser) GetUserID(r *http.Request) (string, error) {
	rp.request = r
	return "TestUserId", nil
}

type RejectingRequestParser struct {
}

func (rp *RejectingRequestParser) GetUserID(r *http.Request) (string, error) {
	return "", api.ErrInvalidJWT
}

type GameDataStoreSpy struct {
	NumberOfActiveGamesReturn                                                                        int
	ActiveGamesReturn                                                                                []*api.Game
	GameWaitingForPlayersReturn                                                                      *api.Game
	StartNewGameReturn                                                                               string
	GameWaitingForPlayersCalled                                                                      bool
	ActiveGamesErr, StartNewGameErr, JoinGameErr, NumberOfActiveGamesErr, GameWaitingForPlayersErr   error
	ActiveGamesUserID, NumberOfActiveGamesUserID, StartNewGameUserID, JoinGameUserID, JoinGameGameID string
}

func (ds *GameDataStoreSpy) ActiveGames(userID string) ([]*api.Game, error) {
	ds.ActiveGamesUserID = userID
	return ds.ActiveGamesReturn, ds.ActiveGamesErr
}

func (ds *GameDataStoreSpy) GameWaitingForPlayers() (*api.Game, error) {
	ds.GameWaitingForPlayersCalled = true
	return ds.GameWaitingForPlayersReturn, ds.GameWaitingForPlayersErr
}

func (ds *GameDataStoreSpy) NumberOfActiveGames(userID string) (int, error) {
	ds.NumberOfActiveGamesUserID = userID
	return ds.NumberOfActiveGamesReturn, ds.NumberOfActiveGamesErr
}

func (ds *GameDataStoreSpy) StartNewGame(userID string) (string, error) {
	ds.StartNewGameUserID = userID
	return ds.StartNewGameReturn, ds.StartNewGameErr
}
func (ds *GameDataStoreSpy) JoinGame(userID string, gameID string) error {
	ds.JoinGameUserID = userID
	ds.JoinGameGameID = gameID
	return ds.JoinGameErr
}

var _ = Describe("newGameEndpoint", func() {

	testUserID := "TestUserId"

	var request *http.Request
	var response *httptest.ResponseRecorder

	BeforeEach(func() {
		request, _ = http.NewRequest("GET", "", nil)
		response = httptest.NewRecorder()
	})

	Context("ServeHTTP method", func() {

		It("Should attempt to get the user from request", func() {
			requestParserSpy := &AcceptingRequestParser{}
			endpoint := api.NewNewGameEndpoint(requestParserSpy, &GameDataStoreSpy{})
			endpoint.ServeHTTP(response, request)

			Expect(requestParserSpy.request).To(BeIdenticalTo(request))
		})

		Context("Given the user is not logged in", func() {
			It("Should return forbidden http response", func() {
				endpoint := api.NewNewGameEndpoint(&RejectingRequestParser{}, &GameDataStoreSpy{})
				endpoint.ServeHTTP(response, request)

				Expect(response.Code).To(BeIdenticalTo(http.StatusForbidden))
			})
		})

		Context("Given the user is logged in", func() {
			requestParserSpy := &AcceptingRequestParser{}
			var gameDataStoreSpy *GameDataStoreSpy
			var endpoint *api.NewGameEndpoint

			BeforeEach(func() {
				gameDataStoreSpy = &GameDataStoreSpy{}
				endpoint = api.NewNewGameEndpoint(requestParserSpy, gameDataStoreSpy)
			})

			It("Should ask the datastore for the users games", func() {
				endpoint.ServeHTTP(response, request)

				Expect(gameDataStoreSpy.NumberOfActiveGamesUserID).To(BeIdenticalTo(testUserID))
			})

			Context("And an error occurs while getting the users games", func() {
				It("Should return an server error", func() {
					gameDataStoreSpy.NumberOfActiveGamesErr = errors.New("Test error")
					endpoint.ServeHTTP(response, request)

					Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
				})
			})

			Context("And the user already has "+strconv.Itoa(api.MAX_ACTIVE_GAMES)+" active games", func() {
				BeforeEach(func() {
					gameDataStoreSpy.NumberOfActiveGamesReturn = api.MAX_ACTIVE_GAMES
				})

				It("Should return an client error", func() {
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusBadRequest))
				})

				It("Should not try to get games waiting for players", func() {
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.GameWaitingForPlayersCalled).To(BeFalse())
				})

				It("Should not try to join a game", func() {
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(""))
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(""))
				})

				It("Should not try to create a new game", func() {
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(""))
				})
			})

			Context("And the user has less than "+strconv.Itoa(api.MAX_ACTIVE_GAMES)+" active games", func() {
				BeforeEach(func() {
					gameDataStoreSpy.NumberOfActiveGamesReturn = 0
				})

				It("Should ask for a vacant game to join", func() {
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.GameWaitingForPlayersCalled).To(BeTrue())
				})

				It("Should join a vacant game if one exists", func() {
					id := "vacant game id"
					gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: id}
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(id))
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should not attempt to create a new game if a vacant one exist", func() {
					id := "vacant game id second test"
					gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: id}
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(""))
				})

				It("Should not attempt join a vacant game if none exists", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(""))
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(""))
				})

				It("Should create a new game if no vacant game exists", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					endpoint.ServeHTTP(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should return OK if no errors occurred", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					gameDataStoreSpy.StartNewGameReturn = "new game id"
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
					Expect(response.Body.String()).To(BeIdenticalTo("new game id"))
				})

				Context("If an error occurs while calling the data store", func() {
					It("Should return an internal server error if the datastore cannot lookup vacant games", func() {
						gameDataStoreSpy.GameWaitingForPlayersErr = errors.New("Error getting vacant games")
						endpoint.ServeHTTP(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})

					It("Should return an internal server error if the datastore cannot join an existing game", func() {
						gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: "game id"}
						gameDataStoreSpy.JoinGameErr = errors.New("Error joining a game")
						endpoint.ServeHTTP(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})

					It("Should return an internal server error if the datastore cannot create a new game", func() {
						gameDataStoreSpy.StartNewGameErr = errors.New("Error creating new game")
						endpoint.ServeHTTP(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})
				})
			})
		})
	})
})
