package neutrinoapi

import "errors"

const FIREBASE_PROJECT_ID = "neutrino-1151"
const DEFAULT_PORT = "5000"

const JWT_HEADER_KEY = "neutrino-user"

var ErrInvalidJWT = errors.New("Invalid JWT supplied.")
var ErrMissingJWT = errors.New("No JWT supplied.")
