package models

type KakaoAuthRequest struct {
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	GrantType         string `json:"grant_type"`
	RedirectURI       string `json:"redirect_uri"`
	AuthorizationCode string `json:"code"`
}

type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type KakaoResponse struct {
	Properties struct {
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
	} `json:"properties"`
	KakaoAccount struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"kakao_account"`
}
