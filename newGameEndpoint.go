package neutrinoapi

import "net/http"

type NewGameEndpoint struct {
	ds GameDataStore
}

func NewNewGameEndpoint(ds GameDataStore) *NewGameEndpoint {
	return &NewGameEndpoint{ds: ds}
}

func (ne *NewGameEndpoint) PerformAction(userID string) (string, int){

	if eligible, statusCode := ne.isEligibleForNewGame(userID); !eligible {
		return "", statusCode
	}

	gameID, err := ne.joinExistingGame(userID)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	if gameID != "" {
		return gameID, http.StatusOK
	}

	gameID, err = ne.ds.StartNewGame(userID)
	if err != nil {
		return "", http.StatusInternalServerError
	}

	return gameID, http.StatusOK
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
	// where the latter one then overrides the first one. TODO I should do something about that if
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
