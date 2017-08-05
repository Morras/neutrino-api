package neutrinoapi_test

import (
	api "github.com/morras/neutrinoapi"
	"github.com/morras/neutrinoapi/spy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"errors"
)

var _ = Describe("getGameEndpoint", func() {

	var request *http.Request
	var response *httptest.ResponseRecorder
	var requestParserSpy *spy.RequestParserSpy
	var dataStoreSpy *spy.GameDataStoreSpy
	var endpoint *api.GetGameEndpoint


	testUserID := "test user id"
	testGame := &api.Game{GameID: "testGameID", SerializedGame: 1234, PlayerOneID: testUserID, PlayerTwoID: "other id"}
	testGame2 := &api.Game{GameID: "testGameID 2", SerializedGame: 4321, PlayerOneID: "not test id", PlayerTwoID: "other id"}

	BeforeEach(func() {
		request = httptest.NewRequest(http.MethodPost, "/", nil)
		response = httptest.NewRecorder()
		requestParserSpy = &spy.RequestParserSpy{}
		dataStoreSpy = &spy.GameDataStoreSpy{}
		endpoint = api.NewGetGameEndpoint(requestParserSpy, dataStoreSpy)
	})

	Context("ServeHTTP method", func() {
		It("Should attempt to get the user from request", func() {
			endpoint.ServeHTTP(response, request)
			Expect(requestParserSpy.Request).To(BeIdenticalTo(request))
		})

		Context("Given the user is not logged in", func() {
			It("Should return forbidden http response", func() {
				requestParserSpy.Err = api.ErrInvalidJWT
				endpoint.ServeHTTP(response, request)
				Expect(response.Code).To(BeIdenticalTo(http.StatusForbidden))
			})
		})

		Context("Given that the user is logged in", func() {
			BeforeEach(func(){
				requestParserSpy.UserID = testUserID
			})

			Context("and specific game is requested", func(){
				const testGameID = "testGameID"
				BeforeEach(func(){
					request = httptest.NewRequest(http.MethodPost, "/?gameID=" + testGameID, nil)
				})

				It("should choose specific game over inactive games option", func(){
					dataStoreSpy.GameReturn = testGame
					endpoint.ServeHTTP(response, request)
					Expect(dataStoreSpy.ActiveGamesUserID).To(BeEmpty())
					Expect(dataStoreSpy.GamesUserID).To(BeEmpty())
				})

				It("Should query the datastore for the game", func(){
					dataStoreSpy.GameReturn = testGame
					endpoint.ServeHTTP(response, request)
					Expect(dataStoreSpy.GameGameID).To(BeIdenticalTo(testGameID))
				})

				It("Should return internal server error if there is a problem talking with the datastore", func(){
					dataStoreSpy.GameErr = errors.New("Error getting a specific game")
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
				})

				It("Should return 404 if the game does not exist", func(){
					dataStoreSpy.GameReturn = nil
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusNotFound))
				})

				It("Should return the request game if it exists", func(){
					dataStoreSpy.GameReturn = testGame
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
					//TODO Verify body
				})

				It("Should return forbidden if the player is not part of the requested game", func(){
					dataStoreSpy.GameReturn = testGame2
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusForbidden))
				})
			})

			Context("and inactive games are requested", func(){
				BeforeEach(func(){
					request = httptest.NewRequest(http.MethodPost, "/?includeInactive=true", nil)
				})

				It("Should query the data store for all games for the player", func(){
					dataStoreSpy.GamesReturn = []*api.Game{}
					endpoint.ServeHTTP(response, request)
					Expect(dataStoreSpy.GamesUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should return a list of all games for the player", func(){
					dataStoreSpy.GamesReturn = []*api.Game{testGame, testGame2}
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
					//TODO Verify body
				})

				It("Should return internal server error if there is a problem talking with the datastore", func(){
					dataStoreSpy.GamesErr = errors.New("Error getting inactive games")
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
				})
			})

			Context("and neither specific or inactive games are requested", func(){
				It("Should query the data store for active games for the player", func(){
					dataStoreSpy.ActiveGamesReturn = []*api.Game{}
					endpoint.ServeHTTP(response, request)
					Expect(dataStoreSpy.ActiveGamesUserID).To(BeIdenticalTo(testUserID))
				})

				It("Should return a list of all active games for the player", func(){
					dataStoreSpy.ActiveGamesReturn = []*api.Game{testGame, testGame2}
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
					//TODO Verify body
				})

				It("Should return internal server error if there is a problem talking with the datastore", func(){
					dataStoreSpy.ActiveGamesErr = errors.New("Error getting active games")
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
				})
			})
		})
	})
})
