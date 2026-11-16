package service

import (
	"context"
	"testing"
)

func TestNormalizeEthiopianPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{input: "0911234567", expected: "+251911234567", wantErr: false},
		{input: "+251911234567", expected: "+251911234567", wantErr: false},
		{input: "251911234567", expected: "+251911234567", wantErr: false},
		{input: "0711234567", expected: "+251711234567", wantErr: false},
		{input: "+251711234567", expected: "+251711234567", wantErr: false},
		{input: "0911-23-45-67", expected: "+251911234567", wantErr: false},
		{input: "(0911) 234567", expected: "+251911234567", wantErr: false},
		{input: "0811234567", expected: "", wantErr: true},
		{input: "12345", expected: "", wantErr: true},
		{input: "", expected: "", wantErr: true},
	}

	for _, tt := range tests {
		got, err := NormalizeEthiopianPhone(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("NormalizeEthiopianPhone(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("NormalizeEthiopianPhone(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestMaskPhoneNumber(t *testing.T) {
	phone := "+251911234567"
	masked := MaskPhoneNumber(phone)
	if masked != "+25191****567" {
		t.Errorf("MaskPhoneNumber(%q) = %q, want +25191****567", phone, masked)
	}

	short := "123"
	if MaskPhoneNumber(short) != "******" {
		t.Errorf("MaskPhoneNumber short phone failed")
	}
}

func TestGenerateSecureOTP(t *testing.T) {
	for i := 0; i < 20; i++ {
		otp, err := GenerateSecureOTP()
		if err != nil {
			t.Fatalf("GenerateSecureOTP() unexpected error: %v", err)
		}
		if len(otp) != 6 {
			t.Errorf("GenerateSecureOTP() length = %d, want 6", len(otp))
		}
		for _, c := range otp {
			if c < '0' || c > '9' {
				t.Errorf("GenerateSecureOTP() non-digit char: %c", c)
			}
		}
	}
}

func TestMockDevProvider(t *testing.T) {
	provider := NewMockDevProvider()
	if provider.ProviderCode() != "mock_dev" {
		t.Errorf("ProviderCode = %s, want mock_dev", provider.ProviderCode())
	}
	if err := provider.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}

	res, err := provider.SendSMS(context.Background(), "+251911234567", "Test code: 123456")
	if err != nil {
		t.Fatalf("SendSMS failed: %v", err)
	}
	if res.Status != "delivered" {
		t.Errorf("Status = %s, want delivered", res.Status)
	}
	if res.ProviderRef == nil || *res.ProviderRef == "" {
		t.Errorf("Expected valid provider reference")
	}
}
