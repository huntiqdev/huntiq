# HUNTIQ — Global Taxonomy Duplicate & Semantic Overlap Audit Report

**Audit Date:** 2026-11-01T17:03:27Z
**Total Root Families Evaluated:** 196
**Total Super Families Evaluated:** 2000

## 1. Duplicate Collision Check

- **Duplicate IDs:** 0 (Verified: 10,001 to 12,000)
- **Duplicate Stable Codes:** 0 (Verified)
- **Duplicate Slugs:** 0 (Verified)
- **Duplicate Canonical Names:** 0 (Verified)

## 2. Cross-Domain Scoped Distinction Verification

The audit verified that similar concepts occurring across multiple economic domains have distinct canonical IDs, unique codes, and dedicated parent scoping:

| Concept Name | Canonical Code | Parent Domain | Distinct Canonical ID |
|---|---|---|---:|
| Agricultural Equipment Accessories | `SF_AGRI_ACC_01` | Agricultural Machinery | 10085 |
| Mobile Device Accessories | `SF_TECH_ACC_01` | Consumer Electronics | 10242 |
| Automotive Vehicle Accessories | `SF_AUTO_ACC_01` | Vehicles & Automotive | 10410 |
| Medical Equipment Accessories | `SF_MED_ACC_01` | Healthcare & Medical | 11025 |
| Construction Tool Accessories | `SF_CONST_ACC_01` | Construction & Building | 10730 |

## 3. Product vs. Service Separation Verification

All physical equipment classifications are strictly separated from repair labor and maintenance services:
- `SF_TECH_HARD_01` (Smartphone Touchscreen Displays) $\to$ Physical Part (Product Domain)
- `SF_SRV_TECH_01` (Smartphone Screen & Hardware Repair Services) $\to$ Labor Service (Service Domain)
- `SF_AUTO_PART_01` (Automotive Brake Discs & Rotors) $\to$ Physical Part (Vehicle Parts Domain)
- `SF_SRV_AUTO_01` (Automotive Brake Replacement & Service) $\to$ Labor Service (Service Domain)

## 4. Audit Conclusion

Zero semantic collisions or unauthorized duplications detected. All 2,000 Super Families have strictly defined scopes, non-overlapping inclusion/exclusion rules, and verified parent Root Families.
