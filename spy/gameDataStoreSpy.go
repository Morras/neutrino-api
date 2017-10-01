package spy

import api "github.com/Morras/neutrinoapi"

type GameDataStoreSpy struct {
	NumberOfActiveGamesUserID string
	NumberOfActiveGamesReturn int
	NumberOfActiveGamesErr    error

	ActiveGamesUserID string
	ActiveGamesReturn []*api.Game
	ActiveGamesErr    error

	GameWaitingForPlayersCalled bool
	GameWaitingForPlayersReturn *api.Game
	GameWaitingForPlayersErr    error

	StartNewGameUserID string
	StartNewGameReturn string
	StartNewGameErr    error

	JoinGameUserID, JoinGameGameID string
	JoinGameErr                    error

	GameGameID string
	GameReturn *api.Game
	GameErr    error

	GamesUserID string
	GamesReturn []*api.Game
	GamesErr    error

	UpdateGameGame *api.Game
	UpdateGameErr  error
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

func (ds *GameDataStoreSpy) Game(gameID string) (*api.Game, error) {
	ds.GameGameID = gameID
	return ds.GameReturn, ds.GameErr
}

func (ds *GameDataStoreSpy) Games(userID string) ([]*api.Game, error) {
	ds.GamesUserID = userID
	return ds.GamesReturn, ds.GamesErr
}

func (ds *GameDataStoreSpy) UpdateGame(game *api.Game) error {
	ds.UpdateGameGame = game
	return ds.UpdateGameErr
}
