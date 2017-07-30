package neutrinoapi

import "net/http"

type NewGameEndpoint struct {
	rp *RequestParser
}

func NewNewGameEndpoint(rp *RequestParser) *NewGameEndpoint {
	return &NewGameEndpoint{rp: rp}
}

func (n *NewGameEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := n.rp.GetUserID(r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}
}
