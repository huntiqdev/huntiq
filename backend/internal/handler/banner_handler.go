package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type BannerHandler struct{ repo repository.BannerRepository }

func NewBannerHandler(repo repository.BannerRepository) *BannerHandler {
	return &BannerHandler{repo: repo}
}

// GET /api/banners?type=buyer|seller  — public
func (h *BannerHandler) ListPublic(c *fiber.Ctx) error {
	t := c.Query("type")
	if t != "buyer" && t != "seller" {
		return response.BadRequest(c, "type must be buyer or seller")
	}
	banners, err := h.repo.ListActive(c.Context(), t)
	if err != nil {
		// Table may not exist yet — return empty gracefully
		return response.OKList(c, []interface{}{})
	}
	if banners == nil {
		banners = []models.HeroBanner{}
	}
	return response.OKList(c, banners)
}

// GET /api/admin/banners
func (h *BannerHandler) ListAll(c *fiber.Ctx) error {
	banners, err := h.repo.ListAll(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	if banners == nil {
		banners = []models.HeroBanner{}
	}
	return response.OKList(c, banners)
}

// POST /api/admin/banners
func (h *BannerHandler) Create(c *fiber.Ctx) error {
	var b models.HeroBanner
	if err := c.BodyParser(&b); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if b.MediaURL == "" {
		return response.BadRequest(c, "media_url is required")
	}
	if b.BannerType != "buyer" && b.BannerType != "seller" {
		return response.BadRequest(c, "banner_type must be buyer or seller")
	}
	if b.MediaType != "image" && b.MediaType != "video" {
		b.MediaType = "image"
	}
	b.IsActive = true
	id, err := h.repo.Create(c.Context(), b)
	if err != nil {
		msg := err.Error()
		if len(msg) > 120 {
			msg = msg[:120]
		}
		return response.BadRequest(c, msg)
	}
	b.ID = id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": b})
}

// PUT /api/admin/banners/:id
func (h *BannerHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var b models.HeroBanner
	if err := c.BodyParser(&b); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if b.MediaType != "image" && b.MediaType != "video" {
		b.MediaType = "image"
	}
	b.ID = id
	if err := h.repo.Update(c.Context(), b); err != nil {
		msg := err.Error()
		if len(msg) > 120 {
			msg = msg[:120]
		}
		return response.BadRequest(c, msg)
	}
	return response.OK(c, fiber.Map{"message": "banner updated"})
}

// DELETE /api/admin/banners/:id
func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.repo.Delete(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "banner deleted"})
}
