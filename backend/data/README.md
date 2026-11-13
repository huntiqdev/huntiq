# HUNTIQ Folder-Based Taxonomy Architecture

## Overview

The taxonomy is a folder-based hierarchical tree. Each folder is a category node:
the folder name is the slug, and a `node.json` file inside holds the multilingual
metadata. Sub-folders are children. A folder with no sub-folders is a leaf node.

There is **no hardcoded Go registry**. Adding a category means creating a folder.

## Directory Structure

```
backend/data/taxonomy/
├── category/                          ← SOURCE OF TRUTH (folder tree)
│   ├── agriculture/                   ← ROOT (folder name = slug)
│   │   ├── node.json                  ← metadata: name, icon, description
│   │   ├── cereal-grain-farming/       ← SUPERFAMILY
│   │   │   ├── node.json
│   │   │   ├── teff-farming/           ← LAYER 3
│   │   │   │   ├── node.json
│   │   │   │   └── highland-teff-varieties/  ← LAYER 4
│   │   │   │       ├── node.json
│   │   │   │       ├── white-magna-teff/     ← LAYER 5 (leaf)
│   │   │   │       │   └── node.json
│   │   │   │       └── ...
│   │   │   └── ...
│   │   └── ...
│   └── ... (50 roots total)
├── categories.json                    ← compiled (auto-generated)
├── taxonomy_tree.json                 ← compiled (auto-generated)
├── locales/                           ← compiled (auto-generated)
│   ├── en.json
│   ├── am.json
│   ├── om.json
│   └── ti.json
├── search_config.json                 ← search synonym config (editable)
└── schema/
    └── taxonomy_schema.json
```

## How to Add a Category

1. Create a folder inside the appropriate parent folder (folder name = slug, kebab-case).
2. Add a `node.json` file with the 4-language names:

```json
{
  "name": {
    "en": "Quinoa Farming",
    "am": "ኪኖዋ እርባታ",
    "om": "Qonna Kiinuwaa",
    "ti": "ኪኖዋ ሕርሻ"
  },
  "aliases": {
    "en": ["quinoa"],
    "am": ["ኪኖዋ"]
  }
}
```

3. Run the compiler:
```bash
cd backend
go run ./cmd/taxonomy-tool compile-roots
```

No parent files need editing. The compiler discovers children by scanning folders.

## How to Remove a Category

Delete the folder. Run compile-roots. Done.

## node.json Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Object with `en`, `am`, `om`, `ti` keys |
| `icon` | No | Lucide icon name |
| `description` | No | Category description |
| `aliases` | No | Object with per-language alias arrays for search |

## Commands

```bash
# Walk category/ folders → generate categories.json, taxonomy_tree.json, locales/
go run ./cmd/taxonomy-tool compile-roots

# Load compiled taxonomy into PostgreSQL
go run ./cmd/taxonomy-tool populate-genuine

# Run database category hierarchy integrity audit
go run ./cmd/taxonomy-tool audit
```

## 50 Consolidated Root Categories

| # | Slug | English Name |
|---|------|-------------|
| 1 | agriculture | Agriculture |
| 2 | fishing-aquaculture | Fishing & Aquaculture |
| 3 | forestry-timber | Forestry & Timber |
| 4 | food-beverage | Food & Beverage |
| 5 | animal-care-veterinary | Animal Care & Veterinary |
| 6 | chemicals-raw-materials | Chemicals & Raw Materials |
| 7 | electronics | Electronics & Devices |
| 8 | computing-it | Computing & Information Technology |
| 9 | transport-mobility | Transport & Mobility |
| 10 | construction-infrastructure | Construction & Infrastructure |
| 11 | property-real-estate | Property & Real Estate |
| 12 | energy-utilities | Energy & Utilities |
| 13 | industrial-manufacturing | Industrial Manufacturing |
| 14 | mining-extractive | Mining & Extractive Industries |
| 15 | packaging-printing | Packaging & Printing |
| 16 | commerce-trade | Commerce & Trade |
| 17 | finance-insurance | Finance & Insurance |
| 18 | business-professional-services | Business & Professional Services |
| 19 | legal-regulatory | Legal & Regulatory Services |
| 20 | office-workplace-supplies | Office & Workplace Supplies |
| 21 | home-appliances-kitchenware | Home Appliances & Kitchenware |
| 22 | furniture-furnishings | Furniture & Home Furnishings |
| 23 | clothing-apparel | Clothing & Apparel |
| 24 | footwear-leather | Footwear & Leather Goods |
| 25 | fashion-accessories-jewelry | Fashion Accessories & Jewelry |
| 26 | textiles-fabrics | Textiles & Fabrics |
| 27 | beauty-personal-care | Beauty & Personal Care |
| 28 | healthcare-medicine | Healthcare & Medicine |
| 29 | education-training | Education & Training |
| 30 | research-scientific | Research & Scientific Services |
| 31 | books-stationery-publishing | Books, Stationery & Publishing |
| 32 | communications-media | Communications & Media |
| 33 | arts-culture-heritage | Arts, Culture & Heritage |
| 34 | entertainment-recreation | Entertainment & Recreation |
| 35 | sports-fitness | Sports & Fitness |
| 36 | government-public-administration | Government & Public Administration |
| 37 | community-social-development | Community & Social Development |
| 38 | religion-faith | Religion & Faith |
| 39 | public-safety-emergency | Public Safety & Emergency Services |
| 40 | security-protective | Security & Protective Services |
| 41 | environment-conservation | Environment & Conservation |
| 42 | water-sanitation | Water & Sanitation |
| 43 | waste-management-recycling | Waste Management & Recycling |
| 44 | logistics-supply-chain | Logistics & Supply Chain |
| 45 | hospitality-accommodation | Hospitality & Accommodation |
| 46 | tourism-travel | Tourism & Travel |
| 47 | events-exhibitions | Events & Exhibitions |
| 48 | cleaning-facility-management | Cleaning & Facility Management |
| 49 | personal-household-services | Personal & Household Services |
| 50 | workforce-employment | Workforce & Employment |

## Statistics

- **Roots:** 50
- **Superfamilies:** 417
- **Total nodes:** 541
- **Max depth:** 5 layers
- **Languages:** English (en), Amharic (am), Afaan Oromoo (om), Tigrinya (ti)
