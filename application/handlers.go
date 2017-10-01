package main

import (
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	api "github.com/Morras/neutrinoapi"
	"github.com/Morras/neutrinoapi/spy"

	fjv "github.com/Morras/firebaseJwtValidator"
	"net/http"
	"errors"
	"strconv"
	"github.com/Morras/go-neutrino/game"
)

var eventParser EventParser
var gameDataStore api.GameDataStore

var getGameEndpoint *api.GetGameEndpoint
var newGameEndpoint *api.NewGameEndpoint
var makeMoveEndpoint *api.MakeMoveEndpoint

const projectID = api.FIREBASE_PROJECT_ID

func init() {
	eventParser = NewEventParser(fjv.NewDefaultTokenValidator(projectID))
	gameDataStore = &spy.GameDataStoreSpy{} //TODO substitute datastore
	getGameEndpoint = api.NewGetGameEndpoint(gameDataStore)
	newGameEndpoint = api.NewNewGameEndpoint(gameDataStore)
	makeMoveEndpoint = api.NewMakeMoveEndpoint(gameDataStore)
}

func GetGameHandler(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	userID, err := eventParser.GetUserID(evt)

	if err != nil {
		return nil, prefixErrorMessageInStatusCode(err, http.StatusForbidden)
	}

	gameID := evt.QueryStringParameters[api.QUERY_GET_GAME_GAME_ID]
	// Do not care about errors as parse errors return false anyway
	includeInactive, _ := strconv.ParseBool(evt.QueryStringParameters[api.QUERY_GET_GAME_INCLUDE_INACTIVE])

	games, statusCode := getGameEndpoint.PerformAction(userID, gameID, includeInactive)
	if statusCode != http.StatusOK {
		return games, wrapStatusCodeInError(statusCode)
	}
	return games, nil
}

func NewGameHandler(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	userID, err := eventParser.GetUserID(evt)

	if err != nil {
		return nil, prefixErrorMessageInStatusCode(err, http.StatusForbidden)
	}

	gameID, statusCode := newGameEndpoint.PerformAction(userID)
	if statusCode != http.StatusOK {
		return "", wrapStatusCodeInError(statusCode)
	}
	return gameID, nil
}

func MakeMoveHandler(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	userID, err := eventParser.GetUserID(evt)

	if err != nil {
		return nil, prefixErrorMessageInStatusCode(err, http.StatusForbidden)
	}

	makeMoveReq, err := eventParser.ExtractMakeMoveRequest(evt)
	if err != nil {
		return nil, prefixErrorMessageInStatusCode(err, http.StatusBadRequest)
	}

	statusCode := makeMoveEndpoint.PerformAction(userID, makeMoveReq, &game.Controller{})
	if statusCode != http.StatusOK {
		return "", wrapStatusCodeInError(statusCode)
	}
	return nil, nil
}

func wrapStatusCodeInError(statusCode int) error {
	return errors.New("[" + strconv.Itoa(statusCode) + "]")
}

func prefixErrorMessageInStatusCode(err error, httpCode int) error {
	return errors.New("[" + strconv.Itoa(httpCode) + "]" + err.Error())
}
