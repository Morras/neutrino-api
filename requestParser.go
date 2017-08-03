package neutrinoapi

import (
	fjv "github.com/morras/firebaseJwtValidator"
	"net/http"
	"strings"
)

type RequestParser interface {
	GetUserID(r *http.Request) (string, error)
}

type FirebaseTokenRequestParser struct {
	validator fjv.TokenValidator
}

func NewRequestParser(validator fjv.TokenValidator) RequestParser {
	return &FirebaseTokenRequestParser{validator: validator}
}

func (rp *FirebaseTokenRequestParser) GetUserID(r *http.Request) (string, error) {

	jwt := r.Header.Get(JWT_HEADER_KEY)

	if jwt == "" {
		return "", ErrMissingJWT
	}

	// Ignoring error as it should have been logged by the library
	valid, _ := rp.validator.Validate(jwt)

	if !valid {
		return "", ErrInvalidJWT
	}

	// We know the format is correct because validation succeeded
	rawClaims := strings.Split(jwt, ".")[1]

	_, claims := fjv.DecodeRawClaims(rawClaims)
	return claims.Sub, nil
}
