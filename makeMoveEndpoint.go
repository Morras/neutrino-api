package neutrinoapi

import (
	"encoding/json"
	"errors"
	"github.com/morras/go-neutrino/game"
	"net/http"
	"io/ioutil"
	"fmt"
)

type MakeMoveRequest struct {
	GameID                                                 string
	NeutrinoFromX, NeutrinoToX, NeutrinoFromY, NeutrinoToY byte
	PieceFromX, PieceToX, PieceFromY, PieceToY             byte
}

type MakeMoveEndpoint struct {
	rp             RequestParser
	ds             GameDataStore
	gameController game.GameController
}

func NewMakeMoveEndpoint(rp RequestParser, ds GameDataStore, gameController game.GameController) *MakeMoveEndpoint {
	return &MakeMoveEndpoint{rp: rp, ds: ds, gameController: gameController}
}

func (mme *MakeMoveEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := mme.rp.GetUserID(r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	makeMoveReq, err := extractMakeMoveRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dsGame, err := mme.ds.Game(makeMoveReq.GameID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actualGame := game.UInt64ToGame(dsGame.SerializedGame)

	if playersTurn := isPlayersTurn(userID, dsGame, actualGame); ! playersTurn{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mme.gameController.PlayGame(actualGame)

	if err = mme.makeMoves(makeMoveReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dsGame.SerializedGame = game.GameToUInt64(mme.gameController.Game())

	if err = mme.ds.UpdateGame(dsGame); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func isPlayersTurn(userID string, datastoreGame *Game, actualGame *game.Game) bool {
	if userID == datastoreGame.PlayerOneID &&
		(actualGame.State == game.Player1NeutrinoMove || actualGame.State == game.Player1Move){
		return true
	} else if userID == datastoreGame.PlayerTwoID &&
		(actualGame.State == game.Player2NeutrinoMove || actualGame.State == game.Player2Move){
		return true
	}
	return false
}

func extractMakeMoveRequest(r *http.Request) (*MakeMoveRequest, error) {
	if r.Body == nil {
		return nil, errors.New("make move endpoint called with empty body")
	}

	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	mmReq := &MakeMoveRequest{}
	if err := json.Unmarshal(bodyContent, mmReq); err != nil {
		fmt.Printf("Error unmarshalling %v", err)
		return nil, err
	}

	if mmReq.GameID == "" {
		return nil, errors.New("Missing game id in make move request body")
	}

	return mmReq, nil
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