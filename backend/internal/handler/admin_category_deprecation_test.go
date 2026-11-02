package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestDeprecatedLegacyCategoryMutations(t *testing.T) {
	app := fiber.New()
	handler := &AdminHandler{}

	app.Post("/api/admin/categories", handler.CreateCategory)
	app.Put("/api/admin/categories/:id", handler.UpdateCategory)
	app.Delete("/api/admin/categories/:id", handler.DeleteCategory)

	tests := []struct {
		name       string
		method     string
		target     string
		body       string
		wantStatus int
	}{
		{
			name:       "Create category legacy write should return 409",
			method:     http.MethodPost,
			target:     "/api/admin/categories",
			body:       `{"name_en":"Test","slug":"test"}`,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Update category legacy write should return 409",
			method:     http.MethodPut,
			target:     "/api/admin/categories/123",
			body:       `{"name_en":"Updated Test"}`,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Delete category legacy write should return 409",
			method:     http.MethodDelete,
			target:     "/api/admin/categories/123",
			body:       "",
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			}
			req := httptest.NewRequest(tt.method, tt.target, bodyReader)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 1000)
			if err != nil {
				t.Fatalf("Failed to execute test request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var resBody map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
				t.Fatalf("Failed to decode response JSON: %v", err)
			}

			if resBody["error"] != "DEPRECATED_MUTATION_ENDPOINT" {
				t.Errorf("Expected error DEPRECATED_MUTATION_ENDPOINT, got %v", resBody["error"])
			}

			if resBody["canonical_endpoint"] != "/api/admin/category-hierarchy" {
				t.Errorf("Expected canonical_endpoint /api/admin/category-hierarchy, got %v", resBody["canonical_endpoint"])
			}
		})
	}
}
