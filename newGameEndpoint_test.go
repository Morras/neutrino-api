package neutrinoapi_test

import (
	"errors"
	api "github.com/Morras/neutrinoapi"
	"github.com/Morras/neutrinoapi/spy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var _ = Describe("newGameEndpoint", func() {

	testUserID := "TestUserId"

	var request *http.Request
	var response *httptest.ResponseRecorder

	BeforeEach(func() {
		request = httptest.NewRequest(http.MethodGet, "/", nil)
		response = httptest.NewRecorder()
	})

	Context("performAction method", func() {

		It("Should attempt to get the user from request", func() {
			requestParserSpy := &spy.RequestParserSpy{}
			endpoint := api.NewNewGameEndpoint(requestParserSpy, &spy.GameDataStoreSpy{})
			endpoint.performAction(response, request)

			Expect(requestParserSpy.Request).To(BeIdenticalTo(request))
		})

		Context("Given the user is not logged in", func() {
			It("Should return forbidden http response", func() {
				requestParserSpy := &spy.RequestParserSpy{Err: api.ErrInvalidJWT}
				endpoint := api.NewNewGameEndpoint(requestParserSpy, &spy.GameDataStoreSpy{})
				endpoint.performAction(response, request)

				Expect(response.Code).To(BeIdenticalTo(http.StatusForbidden))
			})
		})

		Context("Given the user is logged in", func() {
			requestParserSpy := &spy.RequestParserSpy{UserID: testUserID}
			var gameDataStoreSpy *spy.GameDataStoreSpy
			var endpoint *api.NewGameEndpoint

			BeforeEach(func() {
				gameDataStoreSpy = &spy.GameDataStoreSpy{}
				endpoint = api.NewNewGameEndpoint(requestParserSpy, gameDataStoreSpy)
			})

			It("Should ask the datastore for the users games", func() {
				endpoint.performAction(response, request)

				Expect(gameDataStoreSpy.NumberOfActiveGamesUserID).To(BeIdenticalTo(testUserID))
			})

			Context("And an error occurs while getting the users games", func() {
				It("Should return an server error", func() {
					gameDataStoreSpy.NumberOfActiveGamesErr = errors.New("Test error")
					endpoint.performAction(response, request)

					Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
				})
			})

			Context("And the user already has "+strconv.Itoa(api.MAX_ACTIVE_GAMES)+" active games", func() {
				BeforeEach(func() {
					gameDataStoreSpy.NumberOfActiveGamesReturn = api.MAX_ACTIVE_GAMES
				})

				It("Should return an client error", func() {
					endpoint.performAction(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusBadRequest))
				})

				It("Should not try to get games waiting for players", func() {
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.GameWaitingForPlayersCalled).To(BeFalse())
				})

				It("Should not try to join a game", func() {
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(""))
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(""))
				})

				It("Should not try to create a new game", func() {
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(""))
				})
			})

			Context("And the user has less than "+strconv.Itoa(api.MAX_ACTIVE_GAMES)+" active games", func() {
				BeforeEach(func() {
					gameDataStoreSpy.NumberOfActiveGamesReturn = 0
				})

				It("Should ask for a vacant game to join", func() {
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.GameWaitingForPlayersCalled).To(BeTrue())
				})

				It("Should join a vacant game if one exists", func() {
					id := "vacant game id"
					gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: id}
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(id))
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should not attempt to create a new game if a vacant one exist", func() {
					id := "vacant game id second test"
					gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: id}
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(""))
				})

				It("Should not attempt join a vacant game if none exists", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.JoinGameGameID).To(BeIdenticalTo(""))
					Expect(gameDataStoreSpy.JoinGameUserID).To(BeIdenticalTo(""))
				})

				It("Should create a new game if no vacant game exists", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					endpoint.performAction(response, request)
					Expect(gameDataStoreSpy.StartNewGameUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should return OK if no errors occurred", func() {
					gameDataStoreSpy.GameWaitingForPlayersReturn = nil
					gameDataStoreSpy.StartNewGameReturn = "new game id"
					endpoint.performAction(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
					Expect(response.Body.String()).To(BeIdenticalTo("new game id"))
				})

				Context("If an error occurs while calling the data store", func() {
					It("Should return an internal server error if the datastore cannot lookup vacant games", func() {
						gameDataStoreSpy.GameWaitingForPlayersErr = errors.New("Error getting vacant games")
						endpoint.performAction(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})

					It("Should return an internal server error if the datastore cannot join an existing game", func() {
						gameDataStoreSpy.GameWaitingForPlayersReturn = &api.Game{GameID: "game id"}
						gameDataStoreSpy.JoinGameErr = errors.New("Error joining a game")
						endpoint.performAction(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})

					It("Should return an internal server error if the datastore cannot create a new game", func() {
						gameDataStoreSpy.StartNewGameErr = errors.New("Error creating new game")
						endpoint.performAction(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})
				})
			})
		})
	})
})
