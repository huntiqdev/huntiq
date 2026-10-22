package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// SearchHandler exposes search analytics for the admin panel.
type SearchHandler struct {
	repo *repository.SearchRepo
}

func NewSearchHandler(repo *repository.SearchRepo) *SearchHandler {
	return &SearchHandler{repo: repo}
}

// GET /api/admin/search/popular?limit=20
func (h *SearchHandler) PopularSearches(c *fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 20)
	if limit < 1 || limit > 100 {
		limit = 20
	}
	data, err := h.repo.GetPopularSearches(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	if data == nil {
		data = []map[string]interface{}{}
	}
	return response.OKList(c, data)
}

// GET /api/admin/search/zero-results?limit=50
func (h *SearchHandler) ZeroResultQueries(c *fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 50)
	if limit < 1 || limit > 200 {
		limit = 50
	}
	data, err := h.repo.GetZeroResultQueries(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	if data == nil {
		data = []map[string]interface{}{}
	}
	return response.OKList(c, data)
}
