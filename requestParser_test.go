package neutrinoapi_test

import (
	"errors"
	api "github.com/morras/neutrinoapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

// This JWT will not validate using the Firebase validator, but it does not have to as we are mocking that
// out anyway. This token however do include a "sub" claim which cannot be mocked out and which is needed
// for testing the GetUserID returns the correct ID.
// The "sub" in the JWT claims is 1234567890
const MINIMAL_JWT = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.Dt2W1GtOLnnqf4-PUr5Ns_9BuLLmwpWO5zBwN4fokX4"

type acceptingJWTValidator struct {
	inputToken string
}

func (a *acceptingJWTValidator) Validate(token string) (bool, error) {
	a.inputToken = token
	return true, nil
}

type rejectingJWTValidator struct {
}

func (a *rejectingJWTValidator) Validate(token string) (bool, error) {
	return false, errors.New("Validator error")
}

var _ = Describe("RequestParser", func() {

	Context("Given the JWT is not present", func() {
		req, _ := http.NewRequest("GET", "/", nil)
		Context("GetUserID", func() {
			requestParser := api.NewRequestParser(&acceptingJWTValidator{})
			It("Should return an empty id and an error", func() {
				subID, err := requestParser.GetUserID(req)
				Expect(subID).To(BeEmpty())
				Expect(err).To(BeIdenticalTo(api.ErrMissingJWT))
			})
		})
	})

	Context("Given the JWT is present", func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add(api.JWT_HEADER_KEY, MINIMAL_JWT)

		Context("GetUserID", func() {
			validatorSpy := &acceptingJWTValidator{}
			requestParser := api.NewRequestParser(validatorSpy)
			It("Should try an validate the correct request parameter", func() {
				requestParser.GetUserID(req)
				Expect(validatorSpy.inputToken).To(BeIdenticalTo(MINIMAL_JWT))
			})
		})

		Context("Given the JWT validates", func() {
			requestParser := api.NewRequestParser(&acceptingJWTValidator{})
			Context("GetUserID", func() {
				It("Should return the correct ID an no errors", func() {
					subID, err := requestParser.GetUserID(req)
					Expect(subID).To(BeIdenticalTo("1234567890"))
					Expect(err).To(BeNil())
				})
			})
		})

		Context("Given the JWT does not validate", func() {
			requestParser := api.NewRequestParser(&rejectingJWTValidator{})
			Context("GetUserID", func() {
				It("Should return an empty id and an error", func() {
					subID, err := requestParser.GetUserID(req)
					Expect(subID).To(BeEmpty())
					Expect(err).To(BeIdenticalTo(api.ErrInvalidJWT))
				})
			})
		})
	})

})
