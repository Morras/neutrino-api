package neutrinoapi

type GameDataStore interface {
	ActiveGames(userID string) ([]*Game, error)
	NumberOfActiveGames(userID string) (int, error)
	GameWaitingForPlayers() (*Game, error)
	StartNewGame(userID string) (string, error)
	JoinGame(userID string, gameID string) error
}
