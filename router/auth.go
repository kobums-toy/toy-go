package router

import (
	"errors"
	"log"
	"net/http"
	"time"
	"toysgo/config"
	"toysgo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenClaims struct {
	User models.User `json:"user"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JwtAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		if (c.Method() == fiber.MethodPost || c.Method() == fiber.MethodGet) && c.Path() == "/api/user" {
			return c.Next()
		}

		if (c.Method() == fiber.MethodPost || c.Method() == fiber.MethodGet) && c.Path() == "/api/oauth/token" {
			return c.Next()
		}

		values := c.Get("Authorization")
		if values != "" {
			str := values

			if len(str) > 7 && str[:7] == "Bearer " {
				token = str[7:]

				claims := AuthTokenClaims{}
				key := func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("Unexpected Signing Method")
					}
					return []byte(config.SecretCode), nil
				}

				tok, err := jwt.ParseWithClaims(token, &claims, key)
				if err == nil {
					c.Locals("jwt", tok)
					c.Locals("user", &(claims.User))
					return c.Next()
				}
			} else {
				log.Println("Jwt header is broken")
			}
		} else {
			log.Println("Jwt header not found")
		}

		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    "error",
			"message": "not auth",
		})
	}
}

func JwtAuth(email string, passwd string) fiber.Map {
	conn := models.NewConnection()

	manager := models.NewUserManager(conn)
	user := manager.GetByEmail(email)

	if user == nil {
		return fiber.Map{
			"code":    "error",
			"message": "user not found",
		}
	}

	if user.Passwd != passwd {
		return fiber.Map{
			"code":    "error",
			"message": "wrong password",
		}
	}

	at := AuthTokenClaims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	// authManager := models.NewAuthManager(conn)
	// auth := authManager.GetByUser(user.Id)

	// rt := RefreshTokenClaims{
	// 	UserId: user.Id,
	// 	Email:  user.Email,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	// 		// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
	// 	},
	// }

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, _ := atoken.SignedString([]byte(config.SecretCode))

	// rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &rt)
	// signedRefreshToken, _ := rtoken.SignedString([]byte(config.SecretCode))

	// refreshTokenItem := models.Auth{
	// 	User:  user.Id,
	// 	Token: signedRefreshToken,
	// }

	// if auth == nil {
	// 	authManager.Insert(&refreshTokenItem)
	// } else {
	// 	refreshTokenItem.Id = auth.Id
	// 	authManager.Update(&refreshTokenItem)
	// }

	user.Passwd = ""
	return fiber.Map{
		"code":        "ok",
		"accessToken": signedAuthToken,
		"user":        user,
		// "refreshToken": signedRefreshToken,
	}
}

func JwtToken(refreshToken string) fiber.Map {
	values := refreshToken
	if values != "" {
		str := values

		if len(str) > 7 && str[:7] == "Bearer " {
			refreshToken = str[7:]

			claims := RefreshTokenClaims{}
			key := func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Unexpected Signing Method")
				}
				return []byte(config.SecretCode), nil
			}

			_, err := jwt.ParseWithClaims(refreshToken, &claims, key)
			if err == nil {
				conn := models.NewConnection()

				manager := models.NewUserManager(conn)
				user := manager.GetByEmail((claims.Email))

				authManager := models.NewAuthManager((conn))
				auth := authManager.GetByUser((claims.UserId))

				if auth == nil {
					return fiber.Map{
						"code":    "error",
						"message": "token not found",
					}
				}

				if auth.Token != refreshToken {
					return nil
				}

				if user == nil {
					return fiber.Map{
						"code":    "error",
						"message": "user not found",
					}
				}

				at := AuthTokenClaims{
					User: *user,
					RegisteredClaims: jwt.RegisteredClaims{
						// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
					},
				}

				atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
				signedAuthToken, _ := atoken.SignedString([]byte(config.SecretCode))

				user.Passwd = ""
				return fiber.Map{
					"code":        "ok",
					"accessToken": signedAuthToken,
				}

			}
		} else {
			log.Println("Jwt header is broken")
		}
	} else {
		log.Println("Jwt header not found")
	}

	return fiber.Map{
		"code":    "error",
		"message": "not auth",
	}
}

func JwtMe(token string) fiber.Map {
	values := token
	if values != "" {
		str := values

		if len(str) > 7 && str[:7] == "Bearer " {
			token = str[7:]

			claims := AuthTokenClaims{}
			key := func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Unexpected Signing Method")
				}
				return []byte(config.SecretCode), nil
			}

			_, err := jwt.ParseWithClaims(token, &claims, key)
			if err == nil {
				return fiber.Map{
					"code":     "ok",
					"id":       claims.User.Id,
					"name":     claims.User.Name,
					"email":    claims.User.Email,
					"imageUrl": "/logo/codefactory_logo.png",
				}
			}
		} else {
			log.Println("Jwt header is broken")
		}
	} else {
		log.Println("Jwt header not found")
	}

	return fiber.Map{
		"code":    "error",
		"message": "not auth",
	}
}

// GenerateAuthResponse JWT 토큰 생성 및 응답 설정
func GenerateAuthResponse(c *fiber.Ctx, user *models.User, message string) error {
	// JWT 클레임 생성
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
		log.Printf("Error signing JWT: %v\n", err)
		return c.JSON(fiber.Map{
			"code":    "error",
			"message": "Failed to generate JWT",
		})
	}

	// 성공 응답 반환
	return c.JSON(fiber.Map{
		"code":        "success",
		"message":     message,
		"accessToken": signedToken,
		"user":        user,
	})
}
