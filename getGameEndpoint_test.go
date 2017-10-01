package neutrinoapi_test

import (
	api "github.com/Morras/neutrinoapi"
	"github.com/Morras/neutrinoapi/spy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"errors"
	"net/http"
)

var _ = Describe("getGameEndpoint", func() {

	var dataStoreSpy *spy.GameDataStoreSpy
	var endpoint *api.GetGameEndpoint

	testUserID := "test user id"
	testGame := &api.Game{GameID: "testGameID", SerializedGame: 1234, PlayerOneID: testUserID, PlayerTwoID: "other id"}
	testGame2 := &api.Game{GameID: "testGameID 2", SerializedGame: 4321, PlayerOneID: "not test id", PlayerTwoID: "other id"}

	BeforeEach(func() {
		dataStoreSpy = &spy.GameDataStoreSpy{}
		endpoint = api.NewGetGameEndpoint(dataStoreSpy)
	})

	Context("performAction method", func() {

		Context("Given a specific game is requested", func() {
			const userID = "testUserID"
			const gameID = "testGameID"
			const includeInactive = false

			It("should choose specific game over inactive games option", func() {
				dataStoreSpy.GameReturn = testGame
				endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(dataStoreSpy.ActiveGamesUserID).To(BeEmpty())
				Expect(dataStoreSpy.GamesUserID).To(BeEmpty())
			})

			It("Should query the datastore for the game", func() {
				dataStoreSpy.GameReturn = testGame
				endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(dataStoreSpy.GameGameID).To(BeIdenticalTo(gameID))
			})

			It("Should return internal server error if there is a problem talking with the datastore", func() {
				dataStoreSpy.GameErr = errors.New("Error getting a specific game")
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusInternalServerError))
				Expect(games).To(BeEmpty())
			})

			It("Should return 404 if the game does not exist", func() {
				dataStoreSpy.GameReturn = nil
				endpoint.PerformAction(userID, gameID, includeInactive)
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusNotFound))
				Expect(games).To(BeEmpty())
			})

			It("Should return the request game if it exists", func() {
				dataStoreSpy.GameReturn = testGame
				endpoint.PerformAction(userID, gameID, includeInactive)
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusOK))
				Expect(len(games)).To(BeIdenticalTo(1))
				Expect(games[0]).To(BeIdenticalTo(testGame))
			})

			It("Should return forbidden if the player is not part of the requested game", func() {
				dataStoreSpy.GameReturn = testGame2
				endpoint.PerformAction(userID, gameID, includeInactive)
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusForbidden))
				Expect(games).To(BeEmpty())
			})
		})

		Context("and inactive games are requested", func() {

			const userID = "testUserID"
			const gameID = ""
			const includeInactive = true

			It("Should query the data store for all games for the player", func() {
				dataStoreSpy.GamesReturn = []*api.Game{}
				endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(dataStoreSpy.GamesUserID).To(BeIdenticalTo(testUserID))
			})

			It("Should return a list of all games for the player", func() {
				dataStoreSpy.GamesReturn = []*api.Game{testGame, testGame2}
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusOK))
				Expect(len(games)).To(BeIdenticalTo(1))
				Expect(games[0]).To(BeIdenticalTo(testGame))
			})

			It("Should return internal server error if there is a problem talking with the datastore", func() {
				dataStoreSpy.GamesErr = errors.New("Error getting inactive games")
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusInternalServerError))
				Expect(games).To(BeEmpty())
			})
		})

		Context("and neither specific or inactive games are requested", func() {

			const userID = "testUserID"
			const gameID = ""
			const includeInactive = false

			It("Should query the data store for active games for the player", func() {
				dataStoreSpy.ActiveGamesReturn = []*api.Game{}
				endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(dataStoreSpy.ActiveGamesUserID).To(BeIdenticalTo(testUserID))
			})

			It("Should return a list of all active games for the player", func() {
				dataStoreSpy.ActiveGamesReturn = []*api.Game{testGame, testGame2}
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusOK))
				Expect(len(games)).To(BeIdenticalTo(2))
				Expect(games[0]).To(BeIdenticalTo(testGame))
				Expect(games[1]).To(BeIdenticalTo(testGame2))
			})

			It("Should return internal server error if there is a problem talking with the datastore", func() {
				dataStoreSpy.ActiveGamesErr = errors.New("Error getting active games")
				games, code := endpoint.PerformAction(userID, gameID, includeInactive)
				Expect(code).To(BeIdenticalTo(http.StatusInternalServerError))
				Expect(games).To(BeEmpty())
			})
		})
	})
})
