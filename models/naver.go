package models

type NaverAuthRequest struct {
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	GrantType          string `json:"grant_type"`
	RedirectURI        string `json:"redirect_uri"`
	AuthorizationCode  string `json:"code"`
	AuthorizationState string `json:"state"`
}

type NaverTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        string `json:"expires_in"`
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type NaverResponse struct {
	Response struct {
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
		Email        string `json:"email"`
		Name         string `json:"name"`
	} `json:"response"`
}
