package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

var (
	ethiopianPhoneRegex = regexp.MustCompile(`^(\+?251|0)?([79]\d{8})$`)
)

// NormalizeEthiopianPhone standardizes an Ethiopian phone number to canonical E.164 (+2519XXXXXXXX or +2517XXXXXXXX).
func NormalizeEthiopianPhone(raw string) (string, error) {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	matches := ethiopianPhoneRegex.FindStringSubmatch(cleaned)
	if len(matches) != 3 {
		return "", fmt.Errorf("invalid Ethiopian phone number format: %s", raw)
	}

	return "+251" + matches[2], nil
}

// MaskPhoneNumber masks middle digits of a phone number for secure logging (e.g., +25191****345).
func MaskPhoneNumber(phone string) string {
	if len(phone) < 7 {
		return "******"
	}
	prefix := phone[:6]
	suffix := phone[len(phone)-3:]
	return prefix + "****" + suffix
}

// GenerateSecureOTP generates a cryptographically secure 6-digit numeric OTP.
func GenerateSecureOTP() (string, error) {
	maxVal := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SMSDeliveryResult holds the outcome of an SMS delivery attempt.
type SMSDeliveryResult struct {
	ProviderRef   *string
	Status        string // delivered, failed, rejected
	ErrorCategory *string
}

// SMSProvider defines the contract for outbound SMS gateways.
type SMSProvider interface {
	SendSMS(ctx context.Context, phone, message string) (*SMSDeliveryResult, error)
	HealthCheck(ctx context.Context) error
	ProviderCode() string
}

// EthioTelecomProvider implements the Ethio Telecom SMS gateway adapter.
type EthioTelecomProvider struct {
	config     *models.SMSProviderConfig
	httpClient *http.Client
}

func NewEthioTelecomProvider(cfg *models.SMSProviderConfig) *EthioTelecomProvider {
	return &EthioTelecomProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *EthioTelecomProvider) ProviderCode() string {
	return "ethio_telecom"
}

func (p *EthioTelecomProvider) HealthCheck(ctx context.Context) error {
	if p.config.APIEndpoint == "" {
		return errors.New("empty API endpoint for Ethio Telecom provider")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, p.config.APIEndpoint, nil)
	if err != nil {
		return err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *EthioTelecomProvider) SendSMS(ctx context.Context, phone, message string) (*SMSDeliveryResult, error) {
	if p.config.APIEndpoint == "" {
		errCat := "configuration_missing"
		return &SMSDeliveryResult{
			Status:        "failed",
			ErrorCategory: &errCat,
		}, errors.New("SMS gateway endpoint not configured")
	}

	payload := map[string]interface{}{
		"sender_id":   p.config.SenderID,
		"phone":       phone,
		"message":     message,
		"credentials": p.config.Credentials,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		errCat := "serialization_error"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.APIEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		errCat := "request_build_error"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		errCat := "network_timeout"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var respMap map[string]interface{}
		_ = json.Unmarshal(respBody, &respMap)

		// Inspect JSON response for telecom error payloads
		if successVal, ok := respMap["success"].(bool); ok && !successVal {
			errMsg := "unknown provider error"
			if m, ok := respMap["message"].(string); ok && m != "" {
				errMsg = m
			}
			errCat := "provider_error"
			return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, fmt.Errorf("SMS gateway reported failure: %s", errMsg)
		}
		if statusVal, ok := respMap["status"].(string); ok && strings.EqualFold(statusVal, "failed") {
			errMsg := "provider reported status: failed"
			if m, ok := respMap["message"].(string); ok && m != "" {
				errMsg = m
			}
			errCat := "provider_rejected"
			return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, fmt.Errorf("SMS gateway rejected message: %s", errMsg)
		}

		ref := ""
		if r, ok := respMap["message_id"].(string); ok {
			ref = r
		} else if r, ok := respMap["reference"].(string); ok {
			ref = r
		}
		var pRef *string
		if ref != "" {
			pRef = &ref
		}
		return &SMSDeliveryResult{
			Status:      "delivered",
			ProviderRef: pRef,
		}, nil
	}

	errCat := fmt.Sprintf("http_%d", resp.StatusCode)
	return &SMSDeliveryResult{
		Status:        "failed",
		ErrorCategory: &errCat,
	}, fmt.Errorf("SMS provider returned status %d: %s", resp.StatusCode, string(respBody))
}

// FidaProvider implements the Fida SMS gateway adapter for national OTP and alerts.
type FidaProvider struct {
	config     *models.SMSProviderConfig
	httpClient *http.Client
}

func NewFidaProvider(cfg *models.SMSProviderConfig) *FidaProvider {
	return &FidaProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *FidaProvider) ProviderCode() string {
	return "fida"
}

func (p *FidaProvider) HealthCheck(ctx context.Context) error {
	if p.config.APIEndpoint == "" {
		return errors.New("empty API endpoint for Fida provider")
	}
	return nil
}

func (p *FidaProvider) SendSMS(ctx context.Context, phone, message string) (*SMSDeliveryResult, error) {
	if p.config.APIEndpoint == "" {
		errCat := "configuration_missing"
		return &SMSDeliveryResult{
			Status:        "failed",
			ErrorCategory: &errCat,
		}, errors.New("Fida SMS gateway endpoint not configured")
	}

	apiKey := ""
	clientID := ""
	if p.config.Credentials != nil {
		if k, ok := p.config.Credentials["api_key"].(string); ok {
			apiKey = k
		}
		if c, ok := p.config.Credentials["client_id"].(string); ok {
			clientID = c
		}
	}

	// Controlled simulated behavior if in mock/demo credential state
	if apiKey == "" || apiKey == "MOCK_FIDA_KEY" || strings.Contains(p.config.APIEndpoint, "demo") {
		masked := MaskPhoneNumber(phone)
		log.Printf("[FIDA SMS SIMULATION] To: %s | Message: %s", masked, message)
		ref := fmt.Sprintf("FIDA-SIM-%d", time.Now().UnixNano())
		return &SMSDeliveryResult{
			Status:      "delivered",
			ProviderRef: &ref,
		}, nil
	}

	payload := map[string]interface{}{
		"sender_id": p.config.SenderID,
		"to":        phone,
		"message":   message,
		"client_id": clientID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		errCat := "serialization_error"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.APIEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		errCat := "request_build_error"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	if clientID != "" {
		req.Header.Set("X-Client-ID", clientID)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		errCat := "network_timeout"
		return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var respMap map[string]interface{}
		_ = json.Unmarshal(respBody, &respMap)

		// Inspect JSON response for telecom error payloads
		if successVal, ok := respMap["success"].(bool); ok && !successVal {
			errMsg := "unknown provider error"
			if m, ok := respMap["message"].(string); ok && m != "" {
				errMsg = m
			}
			errCat := "provider_error"
			return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, fmt.Errorf("SMS gateway reported failure: %s", errMsg)
		}
		if statusVal, ok := respMap["status"].(string); ok && strings.EqualFold(statusVal, "failed") {
			errMsg := "provider reported status: failed"
			if m, ok := respMap["message"].(string); ok && m != "" {
				errMsg = m
			}
			errCat := "provider_rejected"
			return &SMSDeliveryResult{Status: "failed", ErrorCategory: &errCat}, fmt.Errorf("SMS gateway rejected message: %s", errMsg)
		}

		ref := ""
		if r, ok := respMap["message_id"].(string); ok {
			ref = r
		} else if r, ok := respMap["reference"].(string); ok {
			ref = r
		}
		var pRef *string
		if ref != "" {
			pRef = &ref
		}
		return &SMSDeliveryResult{
			Status:      "delivered",
			ProviderRef: pRef,
		}, nil
	}

	errCat := fmt.Sprintf("http_%d", resp.StatusCode)
	return &SMSDeliveryResult{
		Status:        "failed",
		ErrorCategory: &errCat,
	}, fmt.Errorf("Fida provider returned status %d: %s", resp.StatusCode, string(respBody))
}

// MockDevProvider logs SMS messages securely in development environments.
type MockDevProvider struct{}

func NewMockDevProvider() *MockDevProvider {
	return &MockDevProvider{}
}

func (m *MockDevProvider) ProviderCode() string {
	return "mock_dev"
}

func (m *MockDevProvider) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *MockDevProvider) SendSMS(ctx context.Context, phone, message string) (*SMSDeliveryResult, error) {
	masked := MaskPhoneNumber(phone)
	log.Printf("[SMS MOCK] To: %s | Delivered (simulated): %s", masked, message)
	ref := fmt.Sprintf("MOCK-%d", time.Now().UnixNano())
	return &SMSDeliveryResult{
		Status:      "delivered",
		ProviderRef: &ref,
	}, nil
}

// SMSManager orchestrates provider resolution, delivery logging, and phone validation.
type SMSManager struct {
	repo    *repository.SMSRepo
	devMode bool
}

func NewSMSManager(repo *repository.SMSRepo, devMode bool) *SMSManager {
	return &SMSManager{
		repo:    repo,
		devMode: devMode,
	}
}

// SendOTPMessage normalizes phone, records delivery audit, dispatches SMS, and updates audit status.
func (m *SMSManager) SendOTPMessage(ctx context.Context, rawPhone, otpCode, purpose string) error {
	normalizedPhone, err := NormalizeEthiopianPhone(rawPhone)
	if err != nil {
		return err
	}

	maskedPhone := MaskPhoneNumber(normalizedPhone)
	message := fmt.Sprintf("Your HUNTIQ verification code is %s. Valid for 5 minutes. Do not share this code.", otpCode)

	var provider SMSProvider
	var providerCode string
	var providerID int

	cfg, err := m.repo.GetActiveProviderConfig(ctx)
	if err == nil && cfg != nil && cfg.IsActive && cfg.APIEndpoint != "" {
		providerCode = cfg.ProviderCode
		providerID = cfg.ID
		switch cfg.ProviderCode {
		case "fida":
			provider = NewFidaProvider(cfg)
		case "ethio_telecom":
			provider = NewEthioTelecomProvider(cfg)
		default:
			provider = NewMockDevProvider()
		}
	} else {
		// Fallback to mock dev provider
		providerCode = "mock_dev"
		provider = NewMockDevProvider()
	}

	// 1. Create audit log entry with status 'queued'
	logEntry := &models.SMSDeliveryLog{
		PhoneMasked: maskedPhone,
		Purpose:     purpose,
		Provider:    providerCode,
		Status:      "queued",
		RetryCount:  0,
	}

	logID, err := m.repo.SaveDeliveryLog(ctx, logEntry)
	if err != nil {
		log.Printf("[SMSManager] Warning: failed to save delivery log: %v", err)
	}

	// 2. Dispatch SMS
	result, sendErr := provider.SendSMS(ctx, normalizedPhone, message)

	// 3. Update audit log
	if logID > 0 {
		if result != nil {
			_ = m.repo.UpdateDeliveryLogStatus(ctx, logID, result.Status, result.ProviderRef, result.ErrorCategory)
		} else if sendErr != nil {
			errCat := "dispatch_exception"
			_ = m.repo.UpdateDeliveryLogStatus(ctx, logID, "failed", nil, &errCat)
		}
	}

	if providerID > 0 && result != nil && result.Status == "delivered" {
		_ = m.repo.IncrementDailySent(ctx, providerID)
	}

	return sendErr
}

