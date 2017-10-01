package neutrinoapi

import (
	"github.com/Morras/go-neutrino/game"
	"net/http"
)

type MakeMoveRequest struct {
	GameID                                                 string
	NeutrinoFromX, NeutrinoToX, NeutrinoFromY, NeutrinoToY byte
	PieceFromX, PieceToX, PieceFromY, PieceToY             byte
}

type MakeMoveEndpoint struct {
	ds             GameDataStore
	gameController game.GameController
}

func NewMakeMoveEndpoint(ds GameDataStore) *MakeMoveEndpoint {
	return &MakeMoveEndpoint{ds: ds}
}

func (mme *MakeMoveEndpoint) PerformAction(userID string, makeMoveReq *MakeMoveRequest, gameController game.GameController) int {
	dsGame, err := mme.ds.Game(makeMoveReq.GameID)
	if err != nil {
		return http.StatusInternalServerError
	}

	actualGame := game.UInt64ToGame(dsGame.SerializedGame)

	if playersTurn := isPlayersTurn(userID, dsGame, actualGame); !playersTurn {
		return http.StatusForbidden
	}

	mme.gameController.PlayGame(actualGame)

	if err = mme.makeMoves(makeMoveReq); err != nil {
		return http.StatusBadRequest
	}

	dsGame.SerializedGame = game.GameToUInt64(mme.gameController.Game())

	if err = mme.ds.UpdateGame(dsGame); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
func isPlayersTurn(userID string, datastoreGame *Game, actualGame *game.Game) bool {
	if userID == datastoreGame.PlayerOneID &&
		(actualGame.State == game.Player1NeutrinoMove || actualGame.State == game.Player1Move) {
		return true
	} else if userID == datastoreGame.PlayerTwoID &&
		(actualGame.State == game.Player2NeutrinoMove || actualGame.State == game.Player2Move) {
		return true
	}
	return false
}

func (mme *MakeMoveEndpoint) makeMoves(makeMoveReq *MakeMoveRequest) error {
	neutrinoMove := game.NewMove(makeMoveReq.NeutrinoFromX, makeMoveReq.NeutrinoFromY, makeMoveReq.NeutrinoToX, makeMoveReq.NeutrinoToY)
	if _, err := mme.gameController.MakeMove(neutrinoMove); err != nil {
		return err
	}

	pieceMove := game.NewMove(makeMoveReq.PieceFromX, makeMoveReq.PieceFromY, makeMoveReq.PieceToX, makeMoveReq.PieceToY)
	_, err := mme.gameController.MakeMove(pieceMove)
	return err
}
