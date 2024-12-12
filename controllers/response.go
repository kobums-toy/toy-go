package controllers

import "github.com/gofiber/fiber/v2"

// 공통 성공 응답 관리 함수
func SendResponse(c *fiber.Ctx, code string, message string, data map[string]interface{}) error {
	response := fiber.Map{
		"code":    code,
		"message": message,
	}

	// 추가 데이터가 있는 경우 병합
	for key, value := range data {
		response[key] = value
	}

	return c.JSON(response)
}

// 공통 에러 응답 관리 함수
func SendError(c *fiber.Ctx, message string) error {
	return c.JSON(fiber.Map{
		"code":    "error",
		"message": message,
	})
}
