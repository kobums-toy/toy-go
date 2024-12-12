package global

import (
	"fmt"
	"time"
	"toysgo/config"
	"toysgo/models"

	"github.com/golang-jwt/jwt/v5"
)

// AuthTokenClaims 구조체 정의
type AuthTokenClaims struct {
	User models.User `json:"user"`
	jwt.RegisteredClaims
}

// GenerateAuthToken JWT 토큰 생성
func GenerateAuthToken(user *models.User) (string, error) {
	claims := AuthTokenClaims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)), // 유효기간 6시간
		},
	}

	// JWT 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedToken, err := token.SignedString([]byte(config.SecretCode))
	if err != nil {
		fmt.Printf("Error signing JWT: %v\n", err)
		return "", err
	}
	return signedToken, nil
}

func GetDate(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
