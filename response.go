package applesignin

import (
	"fmt"
)

//nolint:lll
// tokenResponseRaw see https://developer.apple.com/documentation/sign_in_with_apple/tokenresponse.
type tokenResponseRaw struct {
	AccessToken  string `json:"access_token"`  // (Reserved for future use) A token used to access allowed data. Currently, no data set has been defined for access.
	ExpiresIn    int    `json:"expires_in"`    // The amount of time, in seconds, before the access token expires.
	IDToken      string `json:"id_token"`      // A JSON Web Token that contains the user’s identity information.
	RefreshToken string `json:"refresh_token"` // The refresh token used to regenerate new access tokens. Store this token securely on your server.
	TokenType    string `json:"token_type"`    // The type of access token. It will always be bearer.
}

type TokenResponse struct {
	ExpiresIn    int            `json:"expires_in"`
	IDToken      *IDTokenClaims `json:"id_token"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	TokenType    string         `json:"token_type"`
}

const (
	ReasonInvalidRequest       = "invalid_request"
	ReasonInvalidClient        = "invalid_client"
	ReasonInvalidGrant         = "invalid_grant"
	ReasonUnauthorizedClient   = "unauthorized_client"
	ReasonUnsupportedGrantType = "unsupported_grant_type"
	ReasonInvalidScope         = "invalid_scope"
)

// ErrorResponse see https://developer.apple.com/documentation/sign_in_with_apple/errorresponse.
type ErrorResponse struct {
	Reason string `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("apple error response: %s", r.Reason)
}
