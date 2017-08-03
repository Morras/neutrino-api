package neutrinoapi

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

// TODO figure out if these fields should be private or public. I've made GameID public for now to create a test
type Game struct {
	GameID, playerOneID, playerTwoID string
	state                            State
	winningCondition                 WinningCondition
	serializedGame                   uint64
}
