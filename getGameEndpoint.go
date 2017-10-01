package neutrinoapi

import (
	"net/http"
)

type GetGameEndpoint struct {
	ds GameDataStore
}

func NewGetGameEndpoint(ds GameDataStore) *GetGameEndpoint {
	return &GetGameEndpoint{ds: ds}
}

func (ge *GetGameEndpoint) PerformAction(userID string, gameID string, includeInactive bool) ([]*Game, int) { //TODO Should change this away from http response codes I think
	if gameID != "" {
		return ge.getSingleGameFromDataStoreAndReturn(gameID, userID)
	} else if includeInactive{
		return ge.getAllGamesFromDataStoreAndReturn(gameID, userID)
	} else {
		return ge.getActiveGamesFromDataStoreAndReturn(gameID, userID)
	}
}

func (ge *GetGameEndpoint) getActiveGamesFromDataStoreAndReturn(gameID string, userID string) ([]*Game, int) {
	games, err := ge.ds.ActiveGames(userID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return games, http.StatusOK
}

func (ge *GetGameEndpoint) getAllGamesFromDataStoreAndReturn(gameID string, userID string) ([]*Game, int) {
	games, err := ge.ds.Games(userID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return games, http.StatusOK
}

func (ge *GetGameEndpoint) getSingleGameFromDataStoreAndReturn(gameID string, userID string) ([]*Game, int) {
	game, err := ge.ds.Game(gameID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	if game == nil {
		return nil, http.StatusNotFound
	}
	if game.PlayerOneID != userID && game.PlayerTwoID != userID {
		return nil, http.StatusForbidden
	}
	games := []*Game{game}
	return games, http.StatusOK
}