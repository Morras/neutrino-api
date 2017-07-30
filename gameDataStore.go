package firebaseJwtValidator

type GameDataStore interface {
	GetActiveGames(userID string) ([]*Game, error)
	StartNewGame(userID string) error
	JoinGame(userID string, gameID string) error
}
