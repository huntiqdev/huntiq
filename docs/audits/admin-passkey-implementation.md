# HUNTIQ Enterprise Passkeys & WebAuthn Implementation Architecture

**Document:** `docs/audits/admin-passkey-implementation.md`  
**Date:** October 2026  
**Status:** Architecture Specification & Engineering Reference  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. WebAuthn Architecture Principles

HUNTIQ administrative authentication uses standard FIDO2 / WebAuthn specifications (W3C Web Authentication Level 2/3). 

### 1.1 Biometric Privacy & Device Verification
1. **Zero Biometric Collection:** HUNTIQ **never** collects, transmits, processes, or stores raw biometric images (fingerprints, face scans) or proprietary biometric templates.
2. **Local Hardware Matching:** Biometric matching occurs exclusively within the user's secure hardware element (Apple Secure Enclave, Windows Hello TPM, Android Titan M, YubiKey FIDO2 chip).
3. **Cryptographic Assertion Verification:** The device authenticator performs local user verification (biometric or PIN) and, upon success, signs a single-use server challenge with its hardware-backed private key. HUNTIQ verifies the cryptographic signature against the registered public key.
4. **User Verification Flag (`UV`):** For all privileged administrative portals (`platform_owner`, `platform_admin`, `security_admin`), `UserVerification: "required"` is strictly enforced. An assertion without user verification is rejected.

---

## 2. Server-Controlled Ceremonies

### 2.1 Credential Enrollment Ceremony (Registration)

```
Browser / Authenticator                             HUNTIQ Go Backend + PostgreSQL
      │                                                           │
      │ 1. POST /api/admin/passkeys/register/begin               │
      ├──────────────────────────────────────────────────────────>│
      │    (Authenticated session + fresh reauth)                 │ 2. Generate cryptographically
      │                                                           │    secure challenge (32 bytes).
      │                                                           │ 3. Store in `admin_auth_challenges`
      │                                                           │    with 60-second TTL.
      │ 4. Return CredentialCreationOptions                       │
      │<──────────────────────────────────────────────────────────┤
      │    (challenge, rpID, userID, pubKeyCredParams,            │
      │     authenticatorSelection: userVerification="required")  │
      │                                                           │
5. navigator.credentials.create()                                 │
   - Authenticator prompts for TouchID / Hello / Key             │
   - Generates key pair (ES256 / Ed25519)                        │
   - Signs challenge                                              │
      │                                                           │
      │ 6. POST /api/admin/passkeys/register/finish              │
      ├──────────────────────────────────────────────────────────>│
      │    (id, rawId, clientDataJSON, attestationObject)         │ 7. Validate:
      │                                                           │    - Challenge match & expiration
      │                                                           │    - Expected Origin & RP ID
      │                                                           │    - User verification flag
      │                                                           │ 8. Store in `admin_passkeys` table.
      │                                                           │ 9. Consume challenge.
      │ 10. 200 OK: Credential Enrolled                          │ 10. Record immutable audit log.
      │<──────────────────────────────────────────────────────────┤
```

### 2.2 Authentication Ceremony (Assertion)

```
Browser / Authenticator                             HUNTIQ Go Backend + PostgreSQL
      │                                                           │
      │ 1. POST /api/admin/passkeys/login/begin                  │
      │    { portal: "executive", email?: "..." }                 │
      ├──────────────────────────────────────────────────────────>│
      │                                                           │ 2. Generate challenge.
      │                                                           │ 3. Query allowed credential IDs
      │                                                           │    for portal/account.
      │ 4. Return CredentialRequestOptions                        │
      │<──────────────────────────────────────────────────────────┤
      │                                                           │
5. navigator.credentials.get()                                    │
   - Authenticator prompts for biometric / PIN                   │
   - Computes assertion signature                                │
      │                                                           │
      │ 6. POST /api/admin/passkeys/login/finish                 │
      ├──────────────────────────────────────────────────────────>│
      │    (id, clientDataJSON, authenticatorData, signature)     │ 7. Server verification:
      │                                                           │    - Verify challenge & single-use
      │                                                           │    - Verify Origin & RP ID hash
      │                                                           │    - Verify User Verification bit
      │                                                           │    - Verify signature using public key
      │                                                           │    - Verify & update signature counter
      │                                                           │ 8. Issue server session (`admin_sessions`)
      │ 9. Return Session Token & Set HttpOnly Cookie             │ 9. Record login in `admin_logs`.
      │<──────────────────────────────────────────────────────────┤
```

---

## 3. Supported Authenticator Ecosystem

1. **Platform Authenticators:**
   - Apple Touch ID & Face ID (macOS, iOS, iPadOS via Safari, Chrome, Edge)
   - Windows Hello (Fingerprint, Facial Recognition, Security PIN)
   - Android Biometric Prompt (Fingerprint, Face unlock)
   - Linux FIDO2 PAM / libfido2
2. **Roaming / External Authenticators:**
   - YubiKey 5 Series (USB-A, USB-C, NFC, Lightning)
   - Google Titan Security Keys
   - Feitian ePass FIDO2
3. **Cross-Device Passkeys (Hybrid / QR):**
   - FIDO Alliance Client-to-Authenticator Protocol (CTAP2) via Bluetooth BLE proximity handshake.

---

## 4. Cryptographic Storage & Security Measures

- **Public Key Format:** Standard COSE key structure or DER-encoded public key bytes (`BYTEA`).
- **Signature Counter Protection:** The `sign_count` field is verified on every login. If an incoming assertion has a counter less than or equal to the stored counter (for authenticators supporting counters), `clone_warning` is flagged, the session is denied, and security alerting is triggered.
- **Single-Use Challenge Binding:** Challenges in `admin_auth_challenges` have an immediate deletion / `used_at` flag applied within an atomic transaction. Replay attempts fail immediately.

---

## 5. Post-Implementation Security Audit Findings: Endpoints & Re-Auth

1. **Route Alignment Discrepancy:** The Go backend registered endpoints under `/admin/passkeys/*` (`/admin/passkeys/login/begin`, `/admin/passkeys/login/finish`), while the frontend client targeted `/admin/auth/passkey/*`. Routes and payload parameters (`identifier` vs `email`, `challenge` vs `challenge_id`) must be normalized.
2. **Re-Authentication Enforcement:** `mw.RequireReauth` must be attached to `POST /admin/passkeys/register/begin` and `DELETE /admin/passkeys/:id` to protect authenticator enrollment and deletion against session hijacking.

