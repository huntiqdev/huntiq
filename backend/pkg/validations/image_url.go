package validations

import "strings"

// Allowed image URL patterns for security
var allowedImageDomains = []string{
	"/uploads/",                    // Local uploads
	"images.unsplash.com",          // Unsplash
	"lh3.googleusercontent.com",    // Google Images
	"res.cloudinary.com",           // Cloudinary
	"cdn.example.com",              // Your CDN (update as needed)
}

// ValidateImageURL checks if an image URL is from an allowed source
func ValidateImageURL(url string) bool {
	if url == "" {
		return false
	}

	for _, domain := range allowedImageDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}

	return false
}
