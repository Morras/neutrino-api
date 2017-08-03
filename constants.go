package neutrinoapi

import "errors"

// Platform config
const FIREBASE_PROJECT_ID = "neutrino-1151"
const DEFAULT_PORT = "5000"
const JWT_HEADER_KEY = "neutrino-user"

// Gameplay config
const MAX_ACTIVE_GAMES = 5

// Errors
var ErrInvalidJWT = errors.New("Invalid JWT supplied.")
var ErrMissingJWT = errors.New("No JWT supplied.")
