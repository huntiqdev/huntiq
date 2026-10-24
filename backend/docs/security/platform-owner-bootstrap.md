# HUNTIQ Platform Owner Bootstrap & Provisioning Protocol

**Document:** `docs/security/platform-owner-bootstrap.md`  
**Date:** October 2026  
**Status:** Official Security Procedure  
**Classification:** RESTRICTED — Administrative Custodians Only  

---

## 1. Zero-Trust Bootstrap Principles

1. **No Public Registration:** There is no registration URL, signup form, or public endpoint that can create a Platform Owner.
2. **No Hard-Coded Production Credentials:** The repository codebase contains zero default credentials, hardcoded passwords, or backdoor keys.
3. **Identity Verification:** The initial Platform Owner must correspond to an explicitly verified authorized identity.
4. **Hardware Passkey Mandatory:** The Platform Owner account cannot be used with password-only authentication; an authenticating passkey with user verification is strictly required.

---

## 2. Initial Platform Owner Provisioning Procedure

When deploying HUNTIQ to a fresh production environment where no Platform Owner exists:

### Step 1: Run the Secure Bootstrap Command on the Secure Host
```bash
cd /home/poven/GO\ PROJECT/sellbuy/backend
go run ./cmd/admin-recovery --action=bootstrap-owner \
  --email="governance@huntiq.com" \
  --phone="0911000000" \
  --name="HUNTIQ Platform Owner"
```

### Step 2: Receive Single-Use Enrollment Ticket
The CLI generates:
- An initial account record with `role = 'admin'` and `admin_role = 'platform_owner'`.
- A high-entropy, cryptographically random single-use enrollment token.
- A restricted link: `https://huntiq.com/executive-console/login?enrollment=<TOKEN>` (valid for 15 minutes).

### Step 3: Complete Biometric / Hardware Passkey Enrollment
1. The authorized executive opens the enrollment link in a secure browser.
2. The browser initiates the WebAuthn registration ceremony.
3. The executive verifies themselves using Touch ID, Windows Hello, or an external FIDO2 Security Key (YubiKey).
4. The Go backend verifies the attestation, stores the credential public key in `admin_passkeys`, consumes the enrollment ticket, and creates the initial session.
5. An immutable audit entry is written to `admin_logs`.
