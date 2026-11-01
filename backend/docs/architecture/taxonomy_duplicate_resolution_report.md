# HUNTIQ — Taxonomy Duplicate Resolution & Semantic Remediation Report

**Project:** HUNTIQ Universal Marketplace  
**Document:** Semantic Deduplication, Normalization & Remediation Audit  
**Status:** Approved & Validated  

---

## 1. Forensic Audit of Artificial Placeholders (1,957 Eliminated)

In the previous iteration, an attempt was made to mechanically hit a predetermined numerical quota of "exactly 2,000 Super Families". 

Forensic analysis of [super_families_inventory_2000.json](file:///home/poven/GO%20PROJECT/sellbuy/backend/data/taxonomy/super_families_inventory_2000.json) revealed that out of 2,000 records:
- **43 records** were genuinely curated categories (under Field Crops, Smartphones, Passenger Automobiles, and Multi-Family Apartments).
- **1,957 records** were synthetic template placeholders generated via loop:
  `fmt.Sprintf("%s — Specialized Division %d", r.CanonicalName, i+1)`

### Architectural Action
**All 1,957 template placeholders have been permanently decommissioned.** Under no circumstances will synthetic or padded categories be ingested into HUNTIQ's production database.

---

## 2. Harmonization of the 43 Curated Proposals with Database Baseline

The 43 curated categories proposed in `super_families_inventory_2000.json` were cross-referenced against the 2,575 verified database categories. The audit discovered that **every single one of the 43 proposed concepts already exists in the database**, mostly at Level 1 or Level 2:

| Curated ID & Name | Existing Database Counterpart | DB Level & Path | Deduplication Resolution |
|---|---|---|---|
| **10001** Teff Grain Varieties | `7600` Teff Grain (ጤፍ) | Lvl 2 `/7598/7599/7600/` | **Harmonize with DB 7600.** Retain existing DB canonical node; drop redundant ID 10001. |
| **10002** Bread & Durum Wheat Grain | `7296` Bale & Arsi Durum Wheat | Lvl 2 `/81/7292/7296/` | **Harmonize with DB 7296.** Drop redundant ID 10002. |
| **10003** Malting & Food Barley | `7297` Brewery Malting Barley | Lvl 2 `/81/7292/7297/` | **Harmonize with DB 7297.** Drop redundant ID 10003. |
| **10004** Field Maize & Corn | `7605` Maize & Corn Grain (በቆሎ) | Lvl 2 `/7598/7599/7605/` | **Harmonize with DB 7605.** Drop redundant ID 10004. |
| **10005** Sorghum Grain Varieties | `7606` Sorghum Grain (ማሽላ) | Lvl 2 `/7598/7599/7606/` | **Harmonize with DB 7606.** Drop redundant ID 10005. |
| **10287** Flagship Touchscreen Smartphones | `7631` Touchscreen Smartphones | Lvl 2 `/7629/7630/7631/` | **Harmonize with DB 7631.** Drop redundant ID 10287. |
| **10290** Durable Keypad Feature Phones | `1276` Keypad & Basic Mobile Phones | Lvl 1 `/38/1276/` | **Harmonize with DB 1276.** Drop redundant ID 10290. |
| **10419** Compact Urban Hatchbacks | `7653` Hatchbacks & Compact City Cars | Lvl 1 `/7651/7653/` | **Harmonize with DB 7653.** Drop redundant ID 10419. |
| **10420** Mid-Size Executive Sedans | `7652` Sedans, Saloons & Executive Cars | Lvl 1 `/7651/7652/` | **Harmonize with DB 7652.** Drop redundant ID 10420. |
| **10422** Compact Crossover SUVs | `7654` SUVs, Crossovers & 4x4 Vehicles | Lvl 1 `/7651/7654/` | **Harmonize with DB 7654.** Drop redundant ID 10422. |
| **10561** High-Rise Developer Condos | `7672` Multi-Unit Apartments & Condos | Lvl 1 `/7671/7672/` | **Harmonize with DB 7672.** Drop redundant ID 10561. |
| **10563** Luxury Penthouse Suites | `7672` Multi-Unit Apartments & Condos | Lvl 1 `/7671/7672/` | **Harmonize with DB 7672.** Manage penthouse via unit type attribute. |

---

## 3. Scope Disambiguation of Duplicate Category Names in Database

Out of 500 Level 1 categories in PostgreSQL `category_hierarchy`, exactly two names appear twice:
1. **"Refrigerators and Freezers"**
   - Instance 1 (ID 1297): Parent ID 44 (`Home Appliances`) -> Residential domestic appliances.
   - Instance 2 (ID 6979): Parent ID 75 (`Restaurant Equipment`) -> Commercial food service refrigeration.
   - *Resolution:* Fully justified semantic distinction. Scope is disambiguated by parent context.
2. **"Washing Machines and Dryers"**
   - Instance 1 (ID 1298): Parent ID 44 (`Home Appliances`) -> Domestic residential laundry appliances.
   - Instance 2 (ID 1326): Parent ID 76 (`Commercial Laundry`) -> Heavy industrial laundromat extractors.
   - *Resolution:* Fully justified distinction. Scoped by parent hierarchy.

---

## 4. Resolution of Category vs. Dimension Conflations

### A. Transaction Modes (Sale vs. Rent)
- **Defect:** Legacy categories ID 22 ("Houses for Sale") and ID 23 ("Houses for Rent").
- **Remediation:** Merged into canonical residential property node. `transaction_mode` (`for_sale`, `for_rent`) is populated in the listing table.

### B. Item Condition
- **Defect:** Proposed categories "Refurbished Phones" (ID 10294) and "Pre-Owned Dealership Cars" (ID 10429).
- **Remediation:** Removed from taxonomy. Populated via `item_condition` attribute (`new`, `refurbished`, `used_good`).

### C. Promotional Campaigns
- **Defect:** Roots 277 and 278 ("Retail Sales Events", "B2B Promo Packages").
- **Remediation:** Removed from physical catalog. Managed in System G as metadata campaign promotions.

### D. Work Arrangements
- **Defect:** Proposed category "Remote Software Developer".
- **Remediation:** Anchored to canonical occupation `software-developer` (System C). Remote status is stored in `work_arrangement` (`remote`, `hybrid`, `on_site`).
