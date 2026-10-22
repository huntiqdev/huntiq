package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// UploadHandler manages file uploads with injected configuration
type UploadHandler struct {
	baseURL       string
	uploadDir     string
	maxFileSizeMB int
}

// NewUploadHandler creates a new upload handler with configuration
func NewUploadHandler(baseURL, uploadDir string, maxFileSizeMB int) *UploadHandler {
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 20 // default 20MB
	}
	return &UploadHandler{
		baseURL:       baseURL,
		uploadDir:     uploadDir,
		maxFileSizeMB: maxFileSizeMB,
	}
}

// AppBaseURL is kept for backward compatibility with routes that reference it directly.
// Prefer injecting via NewUploadHandler.
var AppBaseURL = ""

// UploadDir resolves the upload directory from env or defaults to ./uploads.
var UploadDir = func() string {
	if d := os.Getenv("UPLOAD_DIR"); d != "" {
		return filepath.Clean(d)
	}
	return "./uploads"
}()

// magicBytes maps the first bytes of a file to its real MIME type.
var magicBytes = []struct {
	mime   string
	ext    string
	prefix []byte
}{
	{"image/jpeg", ".jpg", []byte{0xFF, 0xD8, 0xFF}},
	{"image/png", ".png", []byte{0x89, 0x50, 0x4E, 0x47}},
	{"image/webp", ".webp", []byte{0x52, 0x49, 0x46, 0x46}},
}

func detectImageType(f *multipart.FileHeader) (mime, ext string, err error) {
	src, err := f.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	header := make([]byte, 12)
	if _, err := io.ReadFull(src, header); err != nil {
		return "", "", fmt.Errorf("could not read file header")
	}
	for _, m := range magicBytes {
		if bytes.HasPrefix(header, m.prefix) {
			if m.ext == ".webp" && string(header[8:12]) != "WEBP" {
				return "", "", fmt.Errorf("invalid WebP file")
			}
			return m.mime, m.ext, nil
		}
	}
	return "", "", fmt.Errorf("only JPEG, PNG and WebP images are allowed")
}

// UploadImage handles image uploads to local disk.
// POST /api/upload/image  (requires seller auth)
func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "image file is required")
	}
	
	// Validate file size
	if file.Size <= 0 {
		return response.BadRequest(c, "image file is empty or invalid")
	}
	maxSize := int64(h.maxFileSizeMB) * 1024 * 1024
	if file.Size > maxSize {
		return response.BadRequest(c, fmt.Sprintf("image must be smaller than %d MB", h.maxFileSizeMB))
	}

	_, ext, err := detectImageType(file)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Use 0700 for upload directory (owner only: read/write/execute, no group/world access)
	if err := os.MkdirAll(h.uploadDir, 0700); err != nil { //nolint:gosec // G301: 0700 is intentional — upload dir is owner-only
		return response.InternalError(c, fmt.Errorf("could not create upload directory"))
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to generate filename"))
	}
	filename := hex.EncodeToString(b) + ext
	dest := filepath.Join(h.uploadDir, filename) //nolint:gosec // G filepath.Join: filename is hex-random, containment check below prevents traversal

	// Containment check: ensure the resolved dest is inside uploadDir (CWE-22 path traversal).
	// filepath.Join + Clean can still escape if uploadDir itself is manipulated.
	cleanUpload := filepath.Clean(h.uploadDir)
	if !strings.HasPrefix(filepath.Clean(dest), cleanUpload+string(filepath.Separator)) { //nolint:gosec // G304: dest verified to be inside uploadDir
		return response.InternalError(c, fmt.Errorf("invalid upload path"))
	}

	// SaveFile with explicit 0640 permissions — owner rw, group r, no world access.
	// c.SaveFile uses os.Create which defaults to 0666 (world-readable), so we
	// write manually to enforce secure file permissions (CWE-276).
	src, err := file.Open()
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to open upload"))
	}
	defer src.Close()
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) //nolint:gosec // G304: dest verified above; 0600 = owner-only rw, no group/world access (CWE-276)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to create file"))
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to save file"))
	}

	baseURL := h.baseURL
	if baseURL == "" {
		log.Printf("WARN: APP_BASE_URL not configured — returning relative image URL.")
		return response.OK(c, fiber.Map{"url": "/uploads/" + filename})
	}
	return response.OK(c, fiber.Map{"url": baseURL + "/uploads/" + filename})
}

// UploadImage is the legacy global function for backward compatibility
// Deprecated: Use UploadHandler.UploadImage instead
func UploadImage(c *fiber.Ctx) error {
	// Use deprecated globals
	handler := NewUploadHandler(AppBaseURL, UploadDir, 20)
	return handler.UploadImage(c)
}
