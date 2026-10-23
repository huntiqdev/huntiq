// Package response provides HTTP response helpers for the API.
package response

import (
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	apperrors "github.com/your-org/used-tools-marketplace/pkg/errors"
)

type envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"` // Machine-readable error code for frontend translation
}

// OK sends a single item response.
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(envelope{Success: true, Data: data})
}

// Created sends a 201 Created response.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(envelope{Success: true, Data: data})
}

// OKList sends a list response, guaranteeing [] instead of null for nil slices.
func OKList(c *fiber.Ctx, data interface{}) error {
	if data == nil {
		return c.Status(fiber.StatusOK).JSON(envelope{Success: true, Data: []struct{}{}})
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		return c.Status(fiber.StatusOK).JSON(envelope{Success: true, Data: []struct{}{}})
	}
	return c.Status(fiber.StatusOK).JSON(envelope{Success: true, Data: data})
}

// AppError sends a structured AppError response with both code and message.
// The frontend can use the code for precise translation, falling back to message.
func AppError(c *fiber.Ctx, err *apperrors.AppError) error {
	status := err.HTTPStatus
	if status == 0 {
		status = fiber.StatusBadRequest
	}
	return c.Status(status).JSON(envelope{
		Success: false,
		Error:   err.Message,
		Code:    string(err.Code),
	})
}

func NotFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusNotFound).JSON(envelope{Success: false, Error: msg})
}

func BadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(envelope{Success: false, Error: msg})
}

// BadRequestCode sends a bad request with both a human-readable message and a machine-readable code.
func BadRequestCode(c *fiber.Ctx, msg string, code apperrors.ErrorCode) error {
	return c.Status(fiber.StatusBadRequest).JSON(envelope{
		Success: false,
		Error:   msg,
		Code:    string(code),
	})
}

// UnauthorizedCode sends an unauthorized response with both message and code.
func UnauthorizedCode(c *fiber.Ctx, msg string, code apperrors.ErrorCode) error {
	return c.Status(fiber.StatusUnauthorized).JSON(envelope{
		Success: false,
		Error:   msg,
		Code:    string(code),
	})
}

// ForbiddenCode sends a forbidden response with both message and code.
func ForbiddenCode(c *fiber.Ctx, msg string, code apperrors.ErrorCode) error {
	return c.Status(fiber.StatusForbidden).JSON(envelope{
		Success: false,
		Error:   msg,
		Code:    string(code),
	})
}

// InternalError logs the real error server-side and returns a generic message to the client.
// Never expose raw DB/internal errors to clients.
func InternalError(c *fiber.Ctx, err error) error {
	log.Printf("[ERROR] %s %s: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(envelope{
		Success: false,
		Error:   "an internal error occurred",
		Code:    string(apperrors.ErrSrvInternal),
	})
}

func Unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(envelope{Success: false, Error: msg})
}

func Forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(envelope{Success: false, Error: msg})
}
