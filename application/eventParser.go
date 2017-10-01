package main

import (
	fjv "github.com/Morras/firebaseJwtValidator"
	api "github.com/Morras/neutrinoapi" // TODO move the dependencies away from api
	"strings"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"encoding/json"
	"fmt"
	"errors"
)

type EventParser interface {
	GetUserID(evt *apigatewayproxyevt.Event) (string, error)
	ExtractMakeMoveRequest(evt *apigatewayproxyevt.Event) (*api.MakeMoveRequest, error)
}

type FirebaseTokenEventParser struct {
	validator fjv.TokenValidator
}

func NewEventParser(validator fjv.TokenValidator) EventParser {
	return &FirebaseTokenEventParser{validator: validator}
}

func (parser *FirebaseTokenEventParser) GetUserID(evt *apigatewayproxyevt.Event) (string, error) {
	jwt := evt.Headers[api.JWT_HEADER_KEY]

	if jwt == "" {
		return "", api.ErrMissingJWT
	}

	// Ignoring error as it should have been logged by the library
	valid, _ := parser.validator.Validate(jwt)

	if !valid {
		return "", api.ErrInvalidJWT
	}

	// We know the format is correct because validation succeeded
	rawClaims := strings.Split(jwt, ".")[1]

	_, claims := fjv.DecodeRawClaims(rawClaims)
	return claims.Sub, nil
}

func (parser *FirebaseTokenEventParser) ExtractMakeMoveRequest(evt *apigatewayproxyevt.Event) (*api.MakeMoveRequest, error) {
	bodyContent := []byte(evt.Body)

	mmReq := &api.MakeMoveRequest{}
	if err := json.Unmarshal(bodyContent, mmReq); err != nil {
		fmt.Printf("Error unmarshalling %v", err)
		return nil, err
	}

	if mmReq.GameID == "" {
		return nil, errors.New("Missing game id in make move request body")
	}

	return mmReq, nil
}
