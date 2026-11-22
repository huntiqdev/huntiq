package validations

import "testing"

func TestValidateImageAndDocURLs(t *testing.T) {
	tests := []struct {
		url      string
		isImage  bool
		isDoc    bool
		expected bool
	}{
		// Valid local uploads
		{"/uploads/3c05fe9e3f8905b98a00c73e8f8feef2.jpg", true, true, true},
		{"/uploads/pic.png", true, true, true},
		{"/uploads/avatar.webp", true, true, true},
		{"/uploads/document.pdf", false, true, true}, // PDF valid for doc, not image

		// Valid remote hosts
		{"https://images.unsplash.com/photo-1234.jpg", true, true, true},
		{"https://picsum.photos/400/300", true, true, true},
		{"https://res.cloudinary.com/huntiq/image/upload/v1/sample.png", true, true, true},
		{"https://api.huntiq.et/uploads/trade_license.pdf", false, true, true},

		// Malicious & Exploit Payloads
		{"javascript:alert(1)", true, true, false},
		{"javascript:alert(1)//images.unsplash.com", true, true, false},
		{"data:text/html;base64,PHNjcmlwdD4=", true, true, false},
		{"vbscript:msgbox(1)", true, true, false},
		{"file:///etc/passwd", true, true, false},
		{"/uploads/../etc/passwd", true, true, false},
		{"/uploads/bad.svg", true, true, false},
		{"https://attacker.com/fake.jpg?domain=images.unsplash.com", true, true, false},
		{"http://evil.com/exploit.jpg", true, true, false},
		{"", true, true, false},
		{"   ", true, true, false},
		{"/uploads/foo\x00bar.jpg", true, true, false},
	}

	for _, tt := range tests {
		imgRes := ValidateImageURL(tt.url)
		docRes := ValidateMediaOrDocURL(tt.url)

		if tt.isImage && imgRes != tt.expected {
			t.Errorf("ValidateImageURL(%q) = %v; want %v", tt.url, imgRes, tt.expected)
		}
		if tt.isDoc && docRes != tt.expected {
			t.Errorf("ValidateMediaOrDocURL(%q) = %v; want %v", tt.url, docRes, tt.expected)
		}
	}
}
