package pocket

import "errors"

var (
	// RedirectServerListen defines the redirect server listen addr, should match the redirect url
	RedirectServerListen = "127.0.0.1:38269"
	// RedirectURL stores the redirect url used by the cli
	RedirectURL = "http://127.0.0.1:38269/complete"

	// ErrLoginCanceled Error
	ErrLoginCanceled = errors.New("Login canceled")
	// ErrInvalidOAuth2State Error
	ErrInvalidOAuth2State = errors.New("Invalid OAuth2 state")
)
