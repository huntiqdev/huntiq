package validations

import (
	"net/url"
	"path/filepath"
	"strings"
)

// Allowed domains for remote image assets
var allowedRemoteHosts = map[string]bool{
	"images.unsplash.com":       true,
	"lh3.googleusercontent.com": true,
	"res.cloudinary.com":        true,
	"picsum.photos":             true,
	"huntiq.et":                 true,
	"api.huntiq.et":             true,
	"localhost:8080":            true,
	"127.0.0.1:8080":            true,
	"localhost:3000":            true,
	"127.0.0.1:3000":            true,
}

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var allowedDocExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".pdf":  true,
}

// ValidateImageURL checks if an image URL is from a safe, trusted source.
// Rejects any javascript:, data:, vbscript:, or traversal attempts.
func ValidateImageURL(rawURL string) bool {
	return validateURLWithExtensions(rawURL, allowedImageExtensions)
}

// ValidateMediaOrDocURL checks if a media or document URL (including PDF) is safe and trusted.
func ValidateMediaOrDocURL(rawURL string) bool {
	return validateURLWithExtensions(rawURL, allowedDocExtensions)
}

func validateURLWithExtensions(rawURL string, allowedExts map[string]bool) bool {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" || len(trimmed) > 2048 {
		return false
	}

	// Reject any ASCII control characters or null bytes
	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] < 0x20 || trimmed[i] == 0x7f {
			return false
		}
	}

	// Immediate disallow of dangerous pseudo-protocols
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "data:") ||
		strings.HasPrefix(lower, "vbscript:") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(lower, "blob:") {
		return false
	}

	// Handle local uploads: e.g. /uploads/3c05fe9e3f8905b98a00c73e8f8feef2.jpg
	if strings.HasPrefix(trimmed, "/uploads/") {
		// Prevent path traversal
		if strings.Contains(trimmed, "..") || strings.Contains(trimmed, "\\") {
			return false
		}
		pathOnly := strings.SplitN(trimmed, "?", 2)[0]
		ext := strings.ToLower(filepath.Ext(pathOnly))
		return allowedExts[ext]
	}

	// Absolute HTTP(S) URL validation
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return false
	}

	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return false
	}

	host := strings.ToLower(parsed.Host)
	// Only allow http:// for loopback development
	if parsed.Scheme == "http" && !strings.HasPrefix(host, "localhost") && !strings.HasPrefix(host, "127.0.0.1") {
		return false
	}

	// Check if host is allowed or a valid subdomain of an allowed host
	if !isHostAllowed(host) {
		return false
	}

	// Check extension if applicable (e.g., picsum.photos dynamic URLs might not have extension)
	if strings.Contains(host, "picsum.photos") {
		return true
	}

	pathOnly := parsed.Path
	if pathOnly != "" {
		ext := strings.ToLower(filepath.Ext(pathOnly))
		if ext != "" && !allowedExts[ext] {
			return false
		}
	}

	return true
}

func isHostAllowed(host string) bool {
	if allowedRemoteHosts[host] {
		return true
	}
	// Strip port if present
	if idx := strings.IndexByte(host, ':'); idx != -1 {
		bareHost := host[:idx]
		if allowedRemoteHosts[bareHost] {
			return true
		}
	}
	// Check subdomain
	for allowed := range allowedRemoteHosts {
		if strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}
