package security

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

// WebAuthnManager handles WebAuthn ceremony initialization and verification
type WebAuthnManager struct {
	wa *webauthn.WebAuthn
}

// NewWebAuthnManager creates and configures a WebAuthn manager with explicit RP ID and origins
func NewWebAuthnManager(rpID string, rpOrigins []string) (*WebAuthnManager, error) {
	if rpID == "" {
		rpID = "localhost"
	}
	if len(rpOrigins) == 0 {
		rpOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	}

	// Normalize origins
	var cleanOrigins []string
	for _, o := range rpOrigins {
		o = strings.TrimSpace(o)
		if o != "" {
			cleanOrigins = append(cleanOrigins, o)
		}
	}

	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "HUNTIQ Platform Governance",
		RPID:          rpID,
		RPOrigins:     cleanOrigins,
		AttestationPreference: protocol.PreferNoAttestation,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			RequireResidentKey: protocol.ResidentKeyNotRequired(),
			UserVerification:   protocol.VerificationRequired, // Enterprise strict: UV required
		},
	})
	if err != nil {
		return nil, err
	}

	return &WebAuthnManager{wa: wa}, nil
}

// DeriveRPID extracts a hostname suitable for WebAuthn RP ID from a base URL
func DeriveRPID(frontendURL string) string {
	u, err := url.Parse(frontendURL)
	if err != nil || u.Hostname() == "" {
		return "localhost"
	}
	return u.Hostname()
}

// WebAuthnUserAdapter implements the webauthn.User interface for HUNTIQ models.User
type WebAuthnUserAdapter struct {
	User        *models.User
	Credentials []models.PasskeyCredential
}

// WebAuthnID returns a consistent byte slice representation of user ID
func (u *WebAuthnUserAdapter) WebAuthnID() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(u.User.ID))
	return buf
}

// WebAuthnName returns the administrative principal's username/email
func (u *WebAuthnUserAdapter) WebAuthnName() string {
	if u.User.Email != "" {
		return u.User.Email
	}
	return u.User.Phone
}

// WebAuthnDisplayName returns the human-readable name of the administrator
func (u *WebAuthnUserAdapter) WebAuthnDisplayName() string {
	if u.User.Name != "" {
		return u.User.Name
	}
	return u.WebAuthnName()
}

// WebAuthnCredentials converts stored passkey records to WebAuthn credentials
func (u *WebAuthnUserAdapter) WebAuthnCredentials() []webauthn.Credential {
	var list []webauthn.Credential
	for _, c := range u.Credentials {
		if c.RevokedAt != nil {
			continue // Skip revoked credentials
		}
		var transports []protocol.AuthenticatorTransport
		for _, t := range c.Transports {
			transports = append(transports, protocol.AuthenticatorTransport(t))
		}

		cred := webauthn.Credential{
			ID:              c.CredentialID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       transports,
			Flags: webauthn.CredentialFlags{
				UserPresent:    true,
				UserVerified:   true,
				BackupEligible: c.BackupEligible,
				BackupState:    c.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    c.AAGUID,
				SignCount: c.SignCount,
			},
		}
		list = append(list, cred)
	}
	return list
}

// BeginRegistration starts a WebAuthn enrollment ceremony
func (m *WebAuthnManager) BeginRegistration(user *WebAuthnUserAdapter) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	options, sessionData, err := m.wa.BeginRegistration(
		user,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationRequired,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
	)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

// FinishRegistration validates and parses the credential creation response
func (m *WebAuthnManager) FinishRegistration(user *WebAuthnUserAdapter, sessionData webauthn.SessionData, response *protocol.ParsedCredentialCreationData) (*webauthn.Credential, error) {
	credential, err := m.wa.CreateCredential(user, sessionData, response)
	if err != nil {
		return nil, err
	}
	// Verify user verification occurred
	if !credential.Flags.UserVerified {
		return nil, errors.New("user verification is required for administrative credentials")
	}
	return credential, nil
}

// BeginDiscoverableLogin starts a discoverable / resident key authentication ceremony
func (m *WebAuthnManager) BeginDiscoverableLogin() (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	options, sessionData, err := m.wa.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationRequired),
	)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

// BeginLogin starts a WebAuthn authentication ceremony
func (m *WebAuthnManager) BeginLogin(user *WebAuthnUserAdapter) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	if user == nil || len(user.Credentials) == 0 {
		return m.BeginDiscoverableLogin()
	}
	options, sessionData, err := m.wa.BeginLogin(
		user,
		webauthn.WithUserVerification(protocol.VerificationRequired),
	)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

// FinishLogin validates and verifies the cryptographic assertion
func (m *WebAuthnManager) FinishLogin(user *WebAuthnUserAdapter, sessionData webauthn.SessionData, response *protocol.ParsedCredentialAssertionData) (*webauthn.Credential, error) {
	credential, err := m.wa.ValidateLogin(user, sessionData, response)
	if err != nil {
		return nil, err
	}
	// Verify user verification was performed on hardware
	if !response.Response.AuthenticatorData.Flags.UserVerified() {
		return nil, errors.New("user verification is required for administrative authentication")
	}
	return credential, nil
}

// SerializeSessionData serializes WebAuthn session data to JSON for challenge storage
func SerializeSessionData(data *webauthn.SessionData) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DeserializeSessionData deserializes stored JSON back to webauthn.SessionData
func DeserializeSessionData(s string) (*webauthn.SessionData, error) {
	var data webauthn.SessionData
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// HashToken computes a SHA-256 hash for session tokens or reauth tokens
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	var hexChars = "0123456789abcdef"
	buf := make([]byte, 64)
	for i, b := range h {
		buf[i*2] = hexChars[b>>4]
		buf[i*2+1] = hexChars[b&0x0f]
	}
	return string(buf)
}
