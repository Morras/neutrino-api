package neutrinoapi

import "net/http"

type NewGameEndpoint struct {
	rp RequestParser
	ds GameDataStore
}

func NewNewGameEndpoint(rp RequestParser, ds GameDataStore) *NewGameEndpoint {
	return &NewGameEndpoint{rp: rp, ds: ds}
}

func (ne *NewGameEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := ne.rp.GetUserID(r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if eligible, statusCode := ne.isEligibleForNewGame(userID); !eligible {
		w.WriteHeader(statusCode)
		return
	}

	gameID, err := ne.joinExistingGame(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if gameID != "" {
		w.Write([]byte(gameID))
		return
	}

	gameID, err = ne.ds.StartNewGame(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(gameID))
}

func (ne *NewGameEndpoint) isEligibleForNewGame(userID string) (bool, int) {
	numberOfGames, err := ne.ds.NumberOfActiveGames(userID)
	if err != nil {
		return false, http.StatusInternalServerError
	}
	if numberOfGames >= MAX_ACTIVE_GAMES {
		return false, http.StatusBadRequest
	}

	return true, 0
}

func (ne *NewGameEndpoint) joinExistingGame(userID string) (string, error) {
	// So this is not at all thread safe. It is possible that two players join the same game,
	// where the latter one then overrides the first one. I should do something about that if
	// I ever actually get anyone to play this.
	activeGame, err := ne.ds.GameWaitingForPlayers()

	if err != nil {
		return "", err
	}

	if activeGame != nil {
		if err = ne.ds.JoinGame(userID, activeGame.GameID); err != nil {
			return "", err
		}
		return activeGame.GameID, nil
	}

	return "", nil
}
