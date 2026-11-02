# HUNTIQ Marketplace Entity Model

---

## 1. Relational Entity Schema
- **User:** `user_id`, `phone_number_hash`, `account_tier` (free, gold, enterprise), `verified_badge`, `created_at`.
- **Listing:** `listing_id`, `user_id`, `category_id`, `title`, `description`, `price_etb`, `negotiable`, `listing_type` (rent/sale), `location_id`, `status` (active, expired, sold).
- **Category:** `category_id`, `slug`, `parent_id`, `name_en`, `name_am`, `name_or`, `name_ti`, `attribute_schema`.
- **Location:** `location_id`, `region`, `city`, `subcity`, `woreda`, `neighborhood_landmarks`, `coordinates`.
- **Promotion:** `promotion_id`, `listing_id`, `tier` (featured_7d, top_30d), `paid_via` (telebirr, cbe_birr), `receipt_ref`.
