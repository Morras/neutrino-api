package neutrinoapi_test

import (
	"bytes"
	"errors"
	api "github.com/morras/neutrinoapi"
	g "github.com/morras/go-neutrino/game"
	"github.com/morras/neutrinoapi/spy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("newGameEndpoint", func() {

	testUserID := "TestUserId"

	var request *http.Request
	var response *httptest.ResponseRecorder
	var requestParserSpy *spy.RequestParserSpy
	var dataStoreSpy *spy.GameDataStoreSpy
	var gameControllerSpy *spy.GameControllerSpy
	var endpoint *api.MakeMoveEndpoint

	const validBodyJSON = `{
		"GameID": "TestGameID",
		"NeutrinoFromX": 1,
		"NeutrinoToX": 2,
		"NeutrinoFromY": 3,
		"NeutrinoToY": 4,
		"PieceFromX": 1,
		"PieceToX": 2,
		"PieceFromY": 3,
		"PieceToY": 4
	}`

	var validBody *bytes.Reader

	// Invalid due to missing GameID
	const invalidBodyJSON = `{
		"NeutrinoFromX": 1,
		"NeutrinoToX": 2,
		"NeutrinoFromY": 3,
		"NeutrinoToY": 4,
		"PieceFromX": 1,
		"PieceToX": 2,
		"PieceFromY": 3,
		"PieceToY": 4
	}`

	var invalidBody *bytes.Reader

	BeforeEach(func() {
		invalidBody = bytes.NewReader([]byte(invalidBodyJSON))
		validBody = bytes.NewReader([]byte(validBodyJSON))
		request = httptest.NewRequest(http.MethodPost, "/", nil)
		response = httptest.NewRecorder()
		requestParserSpy = &spy.RequestParserSpy{}
		dataStoreSpy = &spy.GameDataStoreSpy{}
		gameControllerSpy = &spy.GameControllerSpy{}
		endpoint = api.NewMakeMoveEndpoint(requestParserSpy, dataStoreSpy, gameControllerSpy)
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

			Context("and the body is missing", func() {
				It("Should return a bad request", func() {
					// Using http instead of httptest to force a nil body
					request, _ = http.NewRequest(http.MethodPost, "/", nil)
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusBadRequest))
				})
			})

			Context("and the body does not contain a valid request", func() {
				It("Should return a bad request", func() {
					request = httptest.NewRequest(http.MethodPost, "/", invalidBody)
					endpoint.ServeHTTP(response, request)
					Expect(response.Code).To(BeIdenticalTo(http.StatusBadRequest))
				})
			})

			Context("and the body contains a valid request", func() {
				BeforeEach(func() {
					request = httptest.NewRequest(http.MethodPost, "/", validBody)
				})

				It("Should try and get the game", func() {
					dataStoreSpy.GameErr = errors.New("error getting game")
					endpoint.ServeHTTP(response, request)
					Expect(dataStoreSpy.GameGameID).To(BeIdenticalTo("TestGameID"))
				})

				Context("and there was an error getting the game", func() {
					It("Should return an internal server error", func() {
						dataStoreSpy.GameErr = errors.New("error getting game")
						endpoint.ServeHTTP(response, request)
						Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
					})
				})

				Context("and the data store returns a game", func() {
					var game *api.Game
					BeforeEach(func() {
						game = &api.Game{}
						game.SerializedGame = g.GameToUInt64(g.NewStandardGame())
						dataStoreSpy.GameReturn = game
						gameControllerSpy.GameReturn = g.NewStandardGame()
					})

					Context("and it is not the players turn", func() {
						It("Should return forbidden", func() {
							// We are returning standard game so we know that its player ones turn
							game.PlayerOneID = "someoneElse"
							game.PlayerTwoID = testUserID
							endpoint.ServeHTTP(response, request)
							Expect(response.Code).To(BeIdenticalTo(http.StatusForbidden))
						})
					})

					Context("and the move is not valid", func() {
						It("Should return a bad request", func() {
							gameControllerSpy.MakeMoveErr = errors.New("Invalid move")
							endpoint.ServeHTTP(response, request)
							Expect(response.Code).To(BeIdenticalTo(http.StatusBadRequest))
						})
					})

					Context("and the move is valid", func() {
						It("Should attempt to save the game", func() {
							endpoint.ServeHTTP(response, request)
							Expect(dataStoreSpy.UpdateGameGame).ToNot(BeNil())
						})

						Context("but there was an error saving the game", func() {
							It("Should return an internal server error", func() {
								dataStoreSpy.UpdateGameErr = errors.New("error updating game")
								endpoint.ServeHTTP(response, request)
								Expect(response.Code).To(BeIdenticalTo(http.StatusInternalServerError))
							})
						})

						Context("and the game was successfully saved", func() {
							It("Should return status ok", func() {
								endpoint.ServeHTTP(response, request)
								Expect(response.Code).To(BeIdenticalTo(http.StatusOK))
							})
						})
					})

				})
			})
		})
	})
})
