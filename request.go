package applesignin

import (
	"net/url"
)

// ValidateTokenRequest see https://developer.apple.com/documentation/sign_in_with_apple/generate_and_validate_tokens.
type ValidateTokenRequest struct {
	// The identifier (App ID or Services ID) for your app.
	// The identifier must not include your Team ID, to help mitigate sensitive data exposure to the end user.
	// This parameter is required for both authorization code and refresh token validation requests.
	ClientID string

	// A secret JSON Web Token, generated by the developer,
	// that uses the Sign in with Apple private key associated with your developer account.
	// This parameter is required for both authorization code and refresh token validation requests.
	ClientSecret string

	// The authorization code received in an authorization response sent to your app.
	// The code is single-use only and valid for five minutes.
	// This parameter is required for authorization code validation requests.
	Code string

	// The grant type determines how the client app interacts with the validation server.
	// This parameter is required for both authorization code and refresh token validation requests.
	// For authorization code validation, use authorization_code.
	// For refresh token validation requests, use refresh_token.
	GrantType string

	// The destination URI provided in the authorization request when authorizing a user with your app, if applicable.
	// The URI must use the HTTPS protocol, include a domain name, and cannot contain an IP address or localhost.
	// This parameter is required for authorization code validation requests.
	RedirectURI string
}

func (r *ValidateTokenRequest) Encode() string {
	v := url.Values{
		"client_id":     []string{r.ClientID},
		"client_secret": []string{r.ClientSecret},
		"code":          []string{r.Code},
		"grant_type":    []string{r.GrantType},
	}
	if r.RedirectURI != "" {
		v.Set("redirect_uri", r.RedirectURI)
	}
	return v.Encode()
}

// RefreshTokenRequest see https://developer.apple.com/documentation/sign_in_with_apple/generate_and_validate_tokens.
type RefreshTokenRequest struct {
	// The identifier (App ID or Services ID) for your app.
	// The identifier must not include your Team ID, to help mitigate sensitive data exposure to the end user.
	// This parameter is required for both authorization code and refresh token validation requests.
	ClientID string

	// A secret JSON Web Token, generated by the developer,
	// that uses the Sign in with Apple private key associated with your developer account.
	// This parameter is required for both authorization code and refresh token validation requests.
	ClientSecret string

	// The grant type determines how the client app interacts with the validation server.
	// This parameter is required for both authorization code and refresh token validation requests.
	// For authorization code validation, use authorization_code.
	// For refresh token validation requests, use refresh_token.
	GrantType string

	// The refresh token received from the validation server during a authorization request.
	// This parameter is required for refresh token validation requests.
	RefreshToken string
}

func (r *RefreshTokenRequest) Encode() string {
	v := url.Values{
		"client_id":     []string{r.ClientID},
		"client_secret": []string{r.ClientSecret},
		"grant_type":    []string{r.GrantType},
		"refresh_token": []string{r.RefreshToken},
	}
	return v.Encode()
}
