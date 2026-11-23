package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// UploadHandler manages file uploads with injected configuration
type UploadHandler struct {
	baseURL       string
	uploadDir     string
	maxFileSizeMB int
	userRepo      *repository.UserRepo
}

// NewUploadHandler creates a new upload handler with configuration
func NewUploadHandler(baseURL, uploadDir string, maxFileSizeMB int) *UploadHandler {
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5 // default 5MB strictly enforced
	}
	return &UploadHandler{
		baseURL:       baseURL,
		uploadDir:     uploadDir,
		maxFileSizeMB: maxFileSizeMB,
	}
}

func (h *UploadHandler) WithUserRepo(repo *repository.UserRepo) *UploadHandler {
	h.userRepo = repo
	return h
}

// AppBaseURL is kept for backward compatibility with routes that reference it directly.
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

func detectImageHeader(header []byte) (mime, ext string, err error) {
	for _, m := range magicBytes {
		if bytes.HasPrefix(header, m.prefix) {
			if m.ext == ".webp" {
				if len(header) >= 12 && string(header[8:12]) != "WEBP" {
					return "", "", fmt.Errorf("invalid WebP file")
				}
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
		fileMB := float64(file.Size) / (1024 * 1024)
		return response.BadRequest(c, fmt.Sprintf("Image size is %.1f MB. Maximum allowed size is %d MB. Tip: Take a screenshot of the photo or use the crop tool before uploading — screenshots maintain crisp quality and are usually under 1 MB.", fileMB, h.maxFileSizeMB))
	}

	_, ext, err := detectImageType(file)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Use 0700 for upload directory
	if err := os.MkdirAll(h.uploadDir, 0700); err != nil { //nolint:gosec
		return response.InternalError(c, fmt.Errorf("could not create upload directory"))
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to generate filename"))
	}
	filename := hex.EncodeToString(b) + ext
	dest := filepath.Join(h.uploadDir, filename)

	// Containment check: ensure the resolved dest is inside uploadDir (CWE-22 path traversal).
	cleanUpload := filepath.Clean(h.uploadDir)
	if !strings.HasPrefix(filepath.Clean(dest), cleanUpload+string(filepath.Separator)) { //nolint:gosec
		return response.InternalError(c, fmt.Errorf("invalid upload path"))
	}

	src, err := file.Open()
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to open upload"))
	}
	defer src.Close()
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) //nolint:gosec
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

// ── SSRF Defense Helpers ──────────────────────────────────────────────────────

var privateIPBlocks []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918 private
		"172.16.0.0/12",  // RFC1918 private
		"192.168.0.0/16", // RFC1918 private
		"169.254.0.0/16", // RFC3927 link-local / AWS metadata
		"0.0.0.0/8",      // RFC1122 "this" network
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}
	for _, cidr := range cidrs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

// isDisallowedIP checks if an IP is private, loopback, link-local, or cloud metadata.
func isDisallowedIP(ip net.IP) bool {
	if ip == nil || ip.IsUnspecified() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// validateHostForSSRF resolves a hostname and ensures no resolved IPs are disallowed.
func validateHostForSSRF(hostname string) error {
	trimmed := strings.TrimSpace(hostname)
	if trimmed == "" || strings.EqualFold(trimmed, "localhost") {
		return errors.New("disallowed destination host")
	}

	// If direct IP was supplied
	if directIP := net.ParseIP(trimmed); directIP != nil {
		if isDisallowedIP(directIP) {
			return fmt.Errorf("access to internal or private IP %s is prohibited", trimmed)
		}
		return nil
	}

	// Resolve host IPs
	ips, err := net.LookupIP(trimmed)
	if err != nil {
		return fmt.Errorf("failed to resolve host %s: %w", trimmed, err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("no IP addresses resolved for host %s", trimmed)
	}

	for _, ip := range ips {
		if isDisallowedIP(ip) {
			return fmt.Errorf("resolved destination IP %s is private/internal and prohibited", ip.String())
		}
	}
	return nil
}

// ImportImageFromURL imports an external image via URL with strict SSRF defense.
// POST /api/upload/image-url  (requires Gold/Premium seller auth)
func (h *UploadHandler) ImportImageFromURL(c *fiber.Ctx) error {
	// 1. Authentication and role check
	callerID, ok := c.Locals("userID").(int)
	if !ok || callerID <= 0 {
		return response.Unauthorized(c, "authentication required")
	}

	role, _ := c.Locals("role").(string)

	// 2. Gold/Premium tier check (admins bypass)
	if role != "admin" {
		isPrivileged := false
		if h.userRepo != nil {
			user, err := h.userRepo.GetByID(c.Context(), callerID)
			if err == nil && user != nil {
				tier := strings.ToLower(strings.TrimSpace(user.Tier))
				buyerTier := strings.ToLower(strings.TrimSpace(user.BuyerTier))
				if tier == "gold" || tier == "premium" || buyerTier == "gold" || buyerTier == "premium" {
					isPrivileged = true
				}
			}
		}

		if !isPrivileged {
			return response.Forbidden(c, "Importing listing images from external URLs is exclusively available to Gold and Premium members. Upgrade your membership to unlock this capability.")
		}
	}

	// 3. Parse input
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		return response.BadRequest(c, "valid image 'url' is required")
	}

	rawURL := strings.TrimSpace(req.URL)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return response.BadRequest(c, "invalid URL format")
	}

	// 4. Scheme validation
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return response.BadRequest(c, "only http and https schemes are permitted")
	}

	// 5. Hostname SSRF validation
	hostname := parsedURL.Hostname()
	if err := validateHostForSSRF(hostname); err != nil {
		return response.BadRequest(c, fmt.Sprintf("SSRF validation failed: %v", err))
	}

	// 6. Safe HTTP client with socket-level Control check to prevent DNS rebinding
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
		Control: func(network, address string, c syscall.RawConn) error {
			host, _, splitErr := net.SplitHostPort(address)
			if splitErr == nil {
				if ip := net.ParseIP(host); ip != nil {
					if isDisallowedIP(ip) {
						return fmt.Errorf("connection to prohibited IP %s blocked", ip.String())
					}
				}
			}
			return nil
		},
	}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConns:        5,
		IdleConnTimeout:     10 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("stopped after 5 redirects")
			}
			// Validate redirect destination
			if err := validateHostForSSRF(req.URL.Hostname()); err != nil {
				return fmt.Errorf("redirect blocked by SSRF policy: %w", err)
			}
			return nil
		},
	}

	httpReq, err := http.NewRequestWithContext(c.Context(), http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return response.BadRequest(c, "failed to build image fetch request")
	}
	httpReq.Header.Set("User-Agent", "HUNTIQ-Image-Importer/1.0")

	resp, err := client.Do(httpReq)
	if err != nil {
		return response.BadRequest(c, fmt.Sprintf("failed to fetch image: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.BadRequest(c, fmt.Sprintf("remote server returned HTTP %d", resp.StatusCode))
	}

	// 7. Stream and enforce 10MB size limit
	const maxLimitBytes = 10 * 1024 * 1024 // 10MB
	limitReader := io.LimitReader(resp.Body, maxLimitBytes+1)

	headerBuf := make([]byte, 512)
	n, err := io.ReadAtLeast(limitReader, headerBuf, 12)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return response.BadRequest(c, "unable to read image header")
	}
	headerBuf = headerBuf[:n]

	_, ext, err := detectImageHeader(headerBuf)
	if err != nil {
		return response.BadRequest(c, "fetched file is not a valid JPEG, PNG, or WebP image")
	}

	// 8. Generate safe random filename
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to generate filename"))
	}
	filename := hex.EncodeToString(b) + ext
	dest := filepath.Join(h.uploadDir, filename)

	cleanUpload := filepath.Clean(h.uploadDir)
	if !strings.HasPrefix(filepath.Clean(dest), cleanUpload+string(filepath.Separator)) {
		return response.InternalError(c, fmt.Errorf("invalid upload destination"))
	}

	if err := os.MkdirAll(h.uploadDir, 0700); err != nil {
		return response.InternalError(c, fmt.Errorf("could not create upload directory"))
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to create destination file"))
	}
	defer out.Close()

	// Write header buffer first
	if _, err := out.Write(headerBuf); err != nil {
		return response.InternalError(c, fmt.Errorf("failed to write image data"))
	}

	// Copy remaining stream
	written, err := io.Copy(out, limitReader)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("failed to save complete image"))
	}

	totalSize := int64(len(headerBuf)) + written
	if totalSize > maxLimitBytes {
		_ = os.Remove(dest)
		return response.BadRequest(c, "remote image exceeds 10MB limit")
	}

	baseURL := h.baseURL
	if baseURL == "" {
		return response.OK(c, fiber.Map{
			"url":      "/uploads/" + filename,
			"filename": filename,
			"size":     totalSize,
		})
	}
	return response.OK(c, fiber.Map{
		"url":      baseURL + "/uploads/" + filename,
		"filename": filename,
		"size":     totalSize,
	})
}

// UploadImage is the legacy global function for backward compatibility
// Deprecated: Use UploadHandler.UploadImage instead
func UploadImage(c *fiber.Ctx) error {
	handler := NewUploadHandler(AppBaseURL, UploadDir, 20)
	return handler.UploadImage(c)
}
