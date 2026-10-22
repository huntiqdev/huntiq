package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// truncate caps a string to maxLen runes to prevent oversized DB queries.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}

// mustUserID extracts the authenticated user's ID from Fiber locals.
// Returns (0, error response) if the local is missing or wrong type — callers must return the error immediately.
func mustUserID(c *fiber.Ctx) (int, error) {
	id, ok := c.Locals("userID").(int)
	if !ok || id == 0 {
		return 0, response.Unauthorized(c, "unauthorized")
	}
	return id, nil
}

// parseIDParam parses a named URL param as a positive integer.
func parseIDParam(c *fiber.Ctx, name string) (int, error) {
	id, err := strconv.Atoi(c.Params(name))
	if err != nil || id <= 0 {
		return 0, response.BadRequest(c, "invalid "+name)
	}
	return id, nil
}

// parseQueryInt parses a query parameter as an integer with a default value.
// Always returns a valid integer (never errors) — defaults if missing or invalid.
func parseQueryInt(c *fiber.Ctx, key string, defaultVal int) int {
	v, _ := strconv.Atoi(c.Query(key, strconv.Itoa(defaultVal)))
	return v
}

// parsePagination extracts page and limit from query params with defaults and bounds.
// Returns (page, limit) with page >= 1 and limit capped between 1 and maxLimit.
func parsePagination(c *fiber.Ctx, defaultLimit, maxLimit int) (page, limit int) {
	page = parseQueryInt(c, "page", 1)
	limit = parseQueryInt(c, "limit", defaultLimit)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > maxLimit {
		limit = defaultLimit
	}
	return page, limit
}
