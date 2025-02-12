package rest

import (
	"encoding/json"
	"log"
	"toysgo/controllers"
	"toysgo/global"
	"toysgo/models"

	"github.com/go-resty/resty/v2"
)

type KakaoController struct {
	controllers.Controller
}

func (c *KakaoController) Index() {
	ClientID := c.Query("client_id")
	ClientSecret := c.Query("client_secret")
	GrantType := c.Query("grant_type")
	RedirectURI := c.Query("redirect_uri")
	AuthorizationCode := c.Query("code")

	log.Println(ClientID, ClientSecret, GrantType, RedirectURI, AuthorizationCode)

	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"client_id":     ClientID,
			"client_secret": ClientSecret,
			"grant_type":    GrantType,
			"redirect_uri":  RedirectURI,
			"code":          AuthorizationCode,
		}).
		Post("https://kauth.kakao.com/oauth/token")

	if err != nil {
		log.Println("Error communicating with Kakao API:", err)
	}
	var tokenResponse models.KakaoTokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResponse); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}

	c.KakaoUserApi(tokenResponse.AccessToken)
}

func (c *KakaoController) KakaoUserApi(AccessToken string) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+AccessToken).
		Post("https://kapi.kakao.com/v2/user/me")

	if err != nil {
		log.Println("Error communicating with Kakao API:", err)
	}

	var kakaoResp models.KakaoResponse
	if err := json.Unmarshal(resp.Body(), &kakaoResp); err != nil {
		log.Fatalf("JSON Unmarshal Error: %v", err)
	}

	conn := c.NewConnection() // DB 연결
	manager := models.NewUserManager(conn)

	// 카카오에서 반환된 이메일로 사용자 조회
	existingUser := manager.GetByEmail(kakaoResp.KakaoAccount.Email)

	if existingUser != nil {
		// JWT 토큰 생성
		signedAuthToken, err := global.GenerateAuthToken(existingUser)
		if err != nil {
			controllers.SendError(c.Context, "Failed to generate JWT")
			return
		}

		// 성공 응답 반환
		controllers.SendResponse(c.Context, "success", "User already exists", map[string]interface{}{
			"accessToken": signedAuthToken,
			"user":        existingUser,
		})
		return
	}

	// 사용자가 없는 경우 새 사용자 생성
	newUser := &models.User{
		Passwd: "qwer1234!!",
		Name:   kakaoResp.KakaoAccount.Name,
		Email:  kakaoResp.KakaoAccount.Email,
	}

	if err := manager.Insert(newUser); err != nil {
		log.Printf("Error inserting new user: %v\n", err)
		controllers.SendError(c.Context, "Failed to create user")
		return
	}

	// 새로 생성된 사용자 ID 할당
	newUser.Id = manager.GetIdentity()

	// 사용자 생성 성공 응답
	controllers.SendResponse(c.Context, "success", "User created successfully", map[string]interface{}{
		"user": newUser,
	})
}
