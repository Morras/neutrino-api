package spy

import (
	"net/http"
)

type RequestParserSpy struct {
	UserID  string
	Err     error
	Request *http.Request
}

func (rp *RequestParserSpy) GetUserID(r *http.Request) (string, error) {
	rp.Request = r
	return rp.UserID, rp.Err
}
