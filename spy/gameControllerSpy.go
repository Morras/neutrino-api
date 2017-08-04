package spy

import "github.com/morras/go-neutrino/game"

type GameControllerSpy struct {
	PlayGameGame   *game.Game
	GameReturn     *game.Game
	MakeMoveMove   game.Move
	MakeMoveReturn game.State
	MakeMoveErr    error
}

func (spy *GameControllerSpy) PlayGame(game *game.Game) {
	spy.PlayGameGame = game
}

func (spy *GameControllerSpy) Game() *game.Game {
	return spy.GameReturn
}

func (spy *GameControllerSpy) MakeMove(m game.Move) (game.State, error) {
	spy.MakeMoveMove = m
	return spy.MakeMoveReturn, spy.MakeMoveErr
}
