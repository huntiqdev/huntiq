# HUNTIQ — Legacy Tree Integrity & Listing Mapping Forensic Report

**Project:** HUNTIQ Universal Marketplace  
**Document:** Forensic Investigation into Products 31, 32, and 34 Category Paths  
**Status:** Complete & Verified  

---

## 1. Executive Summary

During the post-migration forensic audit, an inquiry was conducted into the category paths assigned to:
- **Product 31:** Sony Alpha A7 IV Full-Frame Mirrorless Camera Body (Category ID `6446`, Path `/41/6446/`)
- **Product 32:** Tafe 75HP 4WD Agricultural Farm Tractor with Disc Plough (Category ID `29`, Path `/6/29/`)
- **Product 34:** Toyota Hilux 2.8 D-4D Double Cab 4x4 Automatic 2021 (Category ID `230`, Path `/77/227/230/`)

All other products were mapped to Migration 035 canonical categories (IDs `7629`–`7710`). This report investigates whether categories `6446`, `29`, and `230` are **legitimately retained** or **incorrectly attached to legacy trees**.

---

## 2. Product 31: Sony Alpha A7 IV Camera Body

### Forensic Findings
- **Assigned Category:** ID `6446` — `Mirrorless & DSLR Digital Cameras`
- **Current Path:** `/41/6446/` (Root ID `41`: `Cameras, Photography and Optics`)
- **Migration 035 Canonical Tree Inspection:**
  A forensic scan of all categories introduced in Migration 035 (IDs `7598`–`7792`) revealed only ONE camera category:
  ```sql
  SELECT id, name_en, path FROM category_hierarchy WHERE id >= 7598 AND name_en ILIKE '%camera%';
  -- Result: ID 7766 | CCTV IP Cameras & Network Video Recorders (ሲሲቲቪ) | Path: /7765/7766/
  ```
  Migration 035 **did not include consumer digital photography cameras**; it only modeled industrial CCTV security surveillance equipment (ID `7766`).
- **Legacy Hierarchy Depth:**
  Root 41 contains a deeply articulated, professional photography taxonomy:
  - ID `6446`: Mirrorless & DSLR Digital Cameras
  - ID `6447`: 45MP+ 8K Cinema Capable Full-Frame Mirrorless Camera Bodies (exact leaf match for Sony A7 IV)
  - ID `6455`: Camera Lenses and Optical Equipment
  - ID `6467`: Handheld Pro Camera Video Stabilizers

### Verdict
**LEGITIMATELY RETAINED.** Root 41 and ID `6446` are the authoritative photography taxonomy in HUNTIQ. Mapping Product 31 to ID 7766 (CCTV) would have been a catastrophic classification error. Product 31 is correctly classified under ID `6446` (or its child `6447`).

---

## 3. Product 32: Tafe 75HP 4WD Agricultural Farm Tractor

### Forensic Findings
- **Assigned Category:** ID `29` — `Tractors & Farm Machinery`
- **Current Path:** `/6/29/` (Root ID `6`: `Agriculture & Farming`)
- **Database Hierarchy Comparison:**
  The database contains two distinct tractor category paths:
  1. *Legacy Core Path:* Root `6` -> ID `29` (`Tractors & Farm Machinery`)
  2. *Expanded Agricultural Path:* Root `80` (`Agriculture, Farming and Agricultural Land`) -> ID `7274` (`Agricultural Tractors, Harvesters and Implements`) -> **ID `7275` (`75HP to 120HP Turbocharged 4WD Heavy Soil Cultivation Tractors`)**

### Verdict
**LEGITIMATELY RETAINED, BUT SUB-OPTIMAL DEPTH.** 
Category ID `29` is an active, valid Level 1 category under Root 6. However, ID `7275` under Root 80 is an exact, horsepower-specific match for "Tafe 75HP Tractor". 
- *Current State:* Fully functional, valid foreign key, indexed path.
- *Recommended Enhancement:* When legacy tree 6 is formally consolidated into Root 80, migrate Product 32 to ID `7275`.

---

## 4. Product 34: Toyota Hilux 2.8 D-4D Double Cab 4x4

### Forensic Findings
- **Assigned Category:** ID `230` — `Double Cabin Pickups (Hilux, D-Max)`
- **Current Path:** `/77/227/230/` (Root ID `77`: `Trucks, Buses and Commercial Transport` -> ID `227`: `Light and Medium Cargo Trucks`)
- **Migration 035 Canonical Tree Inspection:**
  Migration 035 created a parallel canonical pickup node:
  - Root `7655` (`Commercial Vehicles, Trucks & Fleet Transport`) -> ID `7656` (`Light & Medium Cargo Trucks & Double-Cab Pickups`) -> **ID `7657` (`Double & Single Cabin Utility Pickups (ፒክአፕ)`)**

### Verdict
**DUAL-PATH ARTIFACT.**
Both ID `230` and ID `7657` are valid pickup categories:
- ID `230` was selected during initial migration because its legacy name explicitly contained `(Hilux, D-Max)`.
- ID `7657` is the canonical Migration 035 node.
- *Current State:* Product 34 is safely accessible via ID `230`.
- *Recommended Action:* To achieve 100% canonical alignment across all automotive assets, re-target Product 34 to ID `7657` in a scheduled maintenance window.
