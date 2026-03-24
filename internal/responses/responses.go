package responses

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationMeta struct {
	TotalItems int   `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	PageSize   int64 `json:"page_size"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data,omitempty"`
	Meta PaginationMeta `json:"meta"`
}

func NewSuccessResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

func NewSuccessResponseWithData(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func NewErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}

func NewPaginatedResponse(c *fiber.Ctx, status int, data interface{}, meta interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta":    *meta.(*PaginationMeta),
	})
}
