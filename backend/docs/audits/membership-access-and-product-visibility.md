# SellBuy (HUNTIQ) — Membership Access & Product Visibility Audit Report

**Report File:** `docs/audits/membership-access-and-product-visibility.md`  
**Date:** 2026-10-23  
**Target Project:** SellBuy (HUNTIQ) Multi-Category Ethiopian Marketplace  
**Auditor:** Automated Inspection, Architectural Review & Live Browser Verification  
**Repository Root:** `/home/poven/GO PROJECT/sellbuy`

---

## 1. Executive Summary & Verification Decision

### Verdict: **PASS**
Tiered membership enforcement operates under a zero-trust model where authorization rules are enforced strictly on the Go API server. Guests and registered normal buyers can freely discover listings, view high-level summaries, and see prices. However, restricted specifications for Gold/Premium listings are stripped server-side from unauthorized API responses and masked in the frontend with interactive upgrade prompts.

---

## 2. Membership Tier Hierarchy & Access Matrix

| Buyer Tier | Discover Listings | View Basic Info & Images | Contact Seller | View Standard Specs | View Gold Protected Specs | View Premium Protected Specs |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| **Guest** (Unauthenticated) | **Yes** | **Yes** | Requires Auth | **Yes** | **Locked** (Masked) | **Locked** (Masked) |
| **Registered Buyer** (`normal`) | **Yes** | **Yes** | **Yes** | **Yes** | **Locked** (Masked) | **Locked** (Masked) |
| **Gold Member** (`gold`) | **Yes** | **Yes** | **Yes** | **Yes** | **Unlocked** (Full Specs) | **Locked** (Masked) |
| **Premium Member** (`premium`)| **Yes** | **Yes** | **Yes** | **Yes** | **Unlocked** (Full Specs) | **Unlocked** (Full Specs) |

---

## 3. Server-Enforced Field-Level Security

### A. Endpoint Protection (`GET /api/products/:id/attributes`)
In `backend/internal/handler/attribute_handler.go`:
```go
// Enforce buyer tier access for restricted specifications
buyerTier := "guest"
if bt, ok := c.Locals("buyerTier").(string); ok && bt != "" {
    buyerTier = bt
}

if product.RequiredBuyerTier == "gold" && buyerTier != "gold" && buyerTier != "premium" {
    return c.JSON(fiber.Map{
        "success":       true,
        "locked":        true,
        "required_tier": "gold",
        "data":          []interface{}{},
    })
}
```
**Zero-Trust Guarantee:**
- Unauthorized callers receive an empty array (`"data": []`).
- An unauthorized user cannot obtain restricted specifications by modifying client cookies, inspecting React props, or altering localStorage.

### B. Listing Metadata Sanitization (`GET /api/products` & `GET /api/products/:id`)
In `backend/internal/handler/product_handler.go`:
- Flag `specifications_locked: true` is computed dynamically based on the caller's JWT token claims.
- Seller private contact numbers are guarded behind active session verification.

---

## 4. Multi-Capability Account Model

The database migration `028_multi_role_capabilities.sql` decouples identities from rigid single-role enums:
- A single user account (`users`) can possess multiple concurrent capabilities: `buyer`, `seller`, `landlord`, `renter`.
- Active capabilities are tracked in `user_capabilities` with lifecycle states (`active`, `suspended`, `revoked`).
- JWT tokens embed active capabilities (`capabilities: ["buyer", "seller"]`).
- Middleware `RequireCapability(capability)` and `RequireRole(role)` validate against JWT capabilities or user state without requiring logout/login.

---

## 5. Live Test & Verification Evidence

1. **Go Handler Unit Tests:**
   ```bash
   go test ./internal/handler -run TestTierAuthorizationLogic -v
   ```
   - `GuestAccessToGoldSpecsBlocked` -> PASS (200 with `data: [], locked: true`)
   - `NormalBuyerAccessToGoldSpecsBlocked` -> PASS
   - `GoldMemberAccessToGoldSpecsGranted` -> PASS (200 with full attribute list)
   - `PremiumMemberAccessToGoldSpecsGranted` -> PASS

2. **Go Security Unit Tests:**
   ```bash
   go test ./internal/security -v
   ```
   - JWT generation, claims validation, capability encoding/decoding -> PASS.

3. **Live API Execution Evidence:**
   - `GET /api/products/2/attributes` unauthenticated:
     ```json
     {
       "data": [],
       "locked": true,
       "required_tier": "gold",
       "success": true
     }
     ```
   - Screenshot evidence: `product_2_hilux_1792733745798.png` confirms UI locked card: "GOLD MEMBERS ONLY - Restricted Product Specifications".

---

## 6. Business Decisions Requiring Product Review

1. **Tier Subscription Billing:**
   - Buyer tier upgrade endpoints currently transition status directly or via administrative review. Payment gateway integration (Chapa / Telebirr) for automatic tier provisioning can be scheduled for the monetization phase.
2. **Tier Expiration Daemon:**
   - Membership expiry cron worker runs daily; grace period configuration should be confirmed by business stakeholders.
