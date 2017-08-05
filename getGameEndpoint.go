package neutrinoapi

import (
	"net/http"
	"strconv"
	"encoding/json"
)

type GetGameEndpoint struct {
	rp RequestParser
	ds GameDataStore
}

func NewGetGameEndpoint(rp RequestParser, ds GameDataStore) *GetGameEndpoint {
	return &GetGameEndpoint{rp: rp, ds: ds}
}

func (ge *GetGameEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := ge.rp.GetUserID(r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	gameID := r.URL.Query().Get(QUERY_GET_GAME_GAME_ID)
	includeInactive := r.URL.Query().Get(QUERY_GET_GAME_INCLUDE_INACTIVE)

	if gameID != "" {
		ge.getSingleGameFromDataStoreAndReturn(gameID, userID, w)
		return
	} else if includeInactive, _ := strconv.ParseBool(includeInactive); includeInactive{
		ge.getAllGamesFromDataStoreAndReturn(gameID, userID, w)
		return
	} else {
		ge.getActiveGamesFromDataStoreAndReturn(gameID, userID, w)
		return
	}
}

func (ge *GetGameEndpoint) getActiveGamesFromDataStoreAndReturn(gameID string, userID string, w http.ResponseWriter) {
	games, err := ge.ds.ActiveGames(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshalJsonAndWriteResponse(games, w)
}

func (ge *GetGameEndpoint) getAllGamesFromDataStoreAndReturn(gameID string, userID string, w http.ResponseWriter) {
	games, err := ge.ds.Games(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshalJsonAndWriteResponse(games, w)
}

func (ge *GetGameEndpoint) getSingleGameFromDataStoreAndReturn(gameID string, userID string, w http.ResponseWriter) {
	game, err := ge.ds.Game(gameID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if game == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if game.PlayerOneID != userID && game.PlayerTwoID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	marshalJsonAndWriteResponse(game, w)
}

func marshalJsonAndWriteResponse(input interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}