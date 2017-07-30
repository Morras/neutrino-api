package firebaseJwtValidator

type State int8

const (
	INITIALIZING State = iota
	PLAYING
	DONE
)

type WinningCondition int8

const (
	BACK_LINE WinningCondition = iota
	TRAP
	DEFAULT
)

type Game struct {
	gameID, playerOneID, playerTwoID string
	state State
	winningCondition WinningCondition
	serializedGame uint64
}
