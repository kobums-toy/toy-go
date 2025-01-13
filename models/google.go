package models

type GoogleAuthRequest struct {
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	GrantType         string `json:"grant_type"`
	RedirectURI       string `json:"redirect_uri"`
	AuthorizationCode string `json:"code"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
}

type GoogleResponse struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Name           string `json:"name"`
	GivenName      string `json:"family_name"`
	FamilyName     string `json:"given_name"`
	ProfileImage   string `json:"picture"`
	Locale         string `json:"locale"`
}
