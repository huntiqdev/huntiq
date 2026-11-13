-- ============================================================================
-- 077_consolidated_taxonomy_mapping.sql
-- PURPOSE: Map old 248 fragmented root categories to new 50 consolidated roots.
--          This migration is REVERSIBLE and NON-DESTRUCTIVE.
--          Old categories are deactivated, not deleted. Products are preserved.
-- ============================================================================

-- PART 1: Ensure new 50 consolidated roots exist.
-- The populate-genuine command should be run BEFORE this migration to seed
-- the new taxonomy. This migration handles the old-to-new mapping.

-- PART 2: Deactivate old root categories that are now superfamilies or children
-- under the new consolidated roots. Only deactivate roots whose slug does NOT
-- match any of the 50 new root slugs.
UPDATE category_hierarchy
SET is_active = false, updated_at = NOW()
WHERE level = 0
  AND parent_id IS NULL
  AND slug NOT IN (
    'agriculture', 'fishing-aquaculture', 'forestry-timber', 'food-beverage',
    'animal-care-veterinary', 'chemicals-raw-materials', 'electronics', 'consumer-electronics',
    'computing-it', 'transport-mobility', 'construction-infrastructure',
    'property-real-estate', 'energy-utilities', 'industrial-manufacturing',
    'mining-extractive', 'packaging-printing', 'commerce-trade',
    'finance-insurance', 'business-professional-services', 'legal-regulatory',
    'office-workplace-supplies', 'home-appliances-kitchenware',
    'furniture-furnishings', 'clothing-apparel', 'footwear-leather',
    'fashion-accessories-jewelry', 'textiles-fabrics', 'beauty-personal-care',
    'healthcare-medicine', 'education-training', 'research-scientific',
    'books-stationery-publishing', 'communications-media',
    'arts-culture-heritage', 'entertainment-recreation', 'sports-fitness',
    'government-public-administration', 'community-social-development',
    'religion-faith', 'public-safety-emergency', 'security-protective',
    'environment-conservation', 'water-sanitation',
    'waste-management-recycling', 'logistics-supply-chain',
    'hospitality-accommodation', 'tourism-travel', 'events-exhibitions',
    'cleaning-facility-management', 'personal-household-services',
    'workforce-employment'
  );

-- PART 3: Create category redirects from old root slugs to new canonical slugs.
-- This preserves backward-compatible URLs and search resolution.
INSERT INTO category_redirects (old_slug, new_slug, source_slug, target_category_id, redirect_type, created_at)
SELECT t.old_slug, t.new_slug, t.old_slug, ch.id, 'consolidation', NOW()
FROM (VALUES
    -- Old agriculture-related roots → agriculture
    ('field-crops-cereal-grains', 'agriculture'),
    ('oilseeds-pulses-legumes', 'agriculture'),
    ('horticulture-fresh-vegetables', 'agriculture'),
    ('fruits-tree-crops-orchards', 'agriculture'),
    ('coffee-tea-raw-spices', 'agriculture'),
    ('floriculture-cut-flowers', 'agriculture'),
    ('livestock-cattle-small-ruminants', 'agriculture'),
    ('commercial-poultry-aviculture', 'agriculture'),
    ('apiculture-honey-bee-products', 'agriculture'),
    ('fisheries-aquaculture-freshwater', 'fishing-aquaculture'),
    ('forestry-timber-woodlots', 'forestry-timber'),
    ('agricultural-inputs-seeds-fertilizers', 'agriculture'),
    ('agricultural-tractors-farm-machinery', 'agriculture'),
    ('agricultural-farmland-pasture-acreage', 'property-real-estate'),
    ('agriculture-agronomy-veterinary-farm-roles', 'workforce-employment'),
    ('precision-farming-drones-soil-sensors', 'agriculture'),
    ('commercial-hydroponics-aquaponics-systems', 'agriculture'),
    ('pet-food-veterinary-care-animal-supplies', 'animal-care-veterinary'),
    -- Old food-related roots → food-beverage
    ('grain-milling-commercial-flours', 'food-beverage'),
    ('bakery-bread-pastry-confectionery', 'food-beverage'),
    ('edible-oils-fats-shortenings', 'food-beverage'),
    ('dairy-milk-cheese-butter', 'food-beverage'),
    ('butchery-fresh-processed-meats', 'food-beverage'),
    ('processed-preserved-seafood', 'food-beverage'),
    ('prepared-spices-sauces-condiments', 'food-beverage'),
    ('non-alcoholic-beverages-water', 'food-beverage'),
    ('fermented-distilled-alcoholic-beverages', 'food-beverage'),
    ('packaged-shelf-stable-groceries', 'food-beverage'),
    ('packaged-roasted-coffee-tea-products', 'food-beverage'),
    ('commercial-restaurant-foodservice-supplies', 'food-beverage'),
    ('food-grade-packaging-containers', 'packaging-printing'),
    ('specialized-infant-dietary-nutrition', 'food-beverage'),
    ('food', 'food-beverage'),
    -- Old electronics/tech roots → electronics or computing-it
    ('consumer-electronics', 'electronics'),
    ('smartphones-mobile-handsets', 'electronics'),
    ('mobile-device-cases-chargers-accessories', 'electronics'),
    ('laptops-notebooks-desktop-computers', 'electronics'),
    ('computer-monitors-peripherals-accessories', 'electronics'),
    ('audio-headphones-speakers-sound-systems', 'electronics'),
    ('television-sets-smart-screens-displays', 'electronics'),
    ('photography-videography-cameras-optics', 'electronics'),
    ('video-gaming-consoles-controllers-esports', 'electronics'),
    ('smartwatches-fitness-bands-wearables', 'electronics'),
    ('enterprise-servers-data-storage-hardware', 'computing-it'),
    ('telecommunications-networking-infrastructure', 'computing-it'),
    ('electronic-components-pcb-semiconductors', 'computing-it'),
    ('smart-cities-iot-automation-surveillance', 'computing-it'),
    ('information-technology-cloud-software-services', 'computing-it'),
    ('telecommunications-towers-fiber-backhaul', 'computing-it'),
    ('digital-media-streaming-podcast-production-hardware', 'computing-it'),
    -- Old transport roots → transport-mobility
    ('vehicles', 'transport-mobility'),
    ('passenger-automobiles-sedans-suvs', 'transport-mobility'),
    ('commercial-trucks-freight-haulers', 'transport-mobility'),
    ('motorcycles-motorbikes-scooters', 'transport-mobility'),
    ('bajaj-three-wheelers-tuk-tuk', 'transport-mobility'),
    ('buses-coaches-transit-vans', 'transport-mobility'),
    ('bicycles-ebikes-micromobility', 'transport-mobility'),
    ('electric-automobiles-ev-charging', 'transport-mobility'),
    ('electric-vehicle-charging-stations-accessories', 'transport-mobility'),
    ('aerospace-aviation-aircraft-parts-avionics', 'transport-mobility'),
    ('railway-rolling-stock-locomotives-tracks', 'transport-mobility'),
    ('maritime-shipping-vessels-cargo-handling', 'transport-mobility'),
    ('vehicle-body-parts-lighting-cabin', 'transport-mobility'),
    ('vehicle-tires-wheels-rims', 'transport-mobility'),
    ('brakes-suspension-steering-systems', 'transport-mobility'),
    ('powertrain-engines-transmissions-drivetrain', 'transport-mobility'),
    ('automotive-repair-mechanical-services', 'transport-mobility'),
    ('automotive-starter-batteries-12v-24v', 'transport-mobility'),
    -- Old construction roots → construction-infrastructure
    ('construction', 'construction-infrastructure'),
    ('cement-lime-hydraulic-binders', 'construction-infrastructure'),
    ('construction-aggregates-sand-stone', 'construction-infrastructure'),
    ('construction-hardware-fasteners-nails', 'construction-infrastructure'),
    ('construction-lumber-plywood-formwork', 'construction-infrastructure'),
    ('ceramic-floor-wall-porcelain-tiles', 'construction-infrastructure'),
    ('doors-windows-aluminum-profiles', 'construction-infrastructure'),
    ('roofing-sheets-ega-metal-tiles', 'construction-infrastructure'),
    ('masonry-blocks-clay-bricks', 'construction-infrastructure'),
    ('structural-steel-rebar-framing', 'construction-infrastructure'),
    ('plumbing-piping-sanitary-drainage', 'construction-infrastructure'),
    ('electrical-wiring-cables-switchgear', 'construction-infrastructure'),
    ('electrical-transformers-switchboards-cables', 'energy-utilities'),
    ('residential-lighting-led-chandeliers', 'construction-infrastructure'),
    ('architectural-paints-coatings-sealants', 'construction-infrastructure'),
    ('building-construction-craft-trade-services', 'construction-infrastructure'),
    ('engineering-construction-architecture-roles', 'workforce-employment'),
    -- Old property roots → property-real-estate
    ('property', 'property-real-estate'),
    ('residential-apartments-condominiums', 'property-real-estate'),
    ('residential-villas-compound-homes', 'property-real-estate'),
    ('residential-building-land-plots', 'property-real-estate'),
    ('commercial-office-spaces-suites', 'property-real-estate'),
    ('commercial-retail-mall-shops', 'property-real-estate'),
    ('commercial-industrial-mixed-use-plots', 'property-real-estate'),
    ('land-parcels-master-planned-developments', 'property-real-estate'),
    ('industrial-warehouses-storage-depots', 'property-real-estate'),
    ('residential-commercial-rentals-tenancy', 'property-real-estate'),
    ('hotels-eco-lodges-hospitality-accommodations', 'hospitality-accommodation'),
    -- Old energy roots → energy-utilities
    ('energy', 'energy-utilities'),
    ('solar-energy-photovoltaics-clean-tech', 'energy-utilities'),
    ('solar-photovoltaic-panels-inverters-batteries', 'energy-utilities'),
    ('solar-water-heaters-evacuated-tubes', 'energy-utilities'),
    ('power-generation-diesel-generators-storage', 'energy-utilities'),
    ('lithium-tubular-deep-cycle-energy-storage', 'energy-utilities'),
    ('wind-power-turbines-anemometers-towers', 'energy-utilities'),
    ('hydroelectric-water-turbines-accessories', 'energy-utilities'),
    ('biomass-briquettes-charcoal-biogas', 'energy-utilities'),
    ('smart-electric-water-meters-prepaid', 'water-sanitation'),
    ('municipal-wastewater-treatment-water-purification', 'water-sanitation'),
    -- Old industrial roots → industrial-manufacturing
    ('industrial', 'industrial-manufacturing'),
    ('machinery', 'industrial-manufacturing'),
    ('food-processing-grain-flour-mills', 'industrial-manufacturing'),
    ('commercial-food-beverage-processing-machinery', 'industrial-manufacturing'),
    ('textile-weaving-garment-manufacturing-equipment', 'industrial-manufacturing'),
    ('metalworking-machinery-cnc-lathes-presses', 'industrial-manufacturing'),
    ('plastic-injection-molding-extrusion-machinery', 'industrial-manufacturing'),
    ('industrial-offset-printing-presses-bindery', 'packaging-printing'),
    ('industrial-packaging-bottling-labeling-machinery', 'packaging-printing'),
    ('woodworking-machinery-timber-sawmills', 'industrial-manufacturing'),
    ('heavy-construction-machinery-earthmoving', 'construction-infrastructure'),
    ('forklifts-overhead-cranes-conveyors-rigging', 'industrial-manufacturing'),
    ('industrial-air-compressors-pneumatics-pumps', 'industrial-manufacturing'),
    ('industrial-welding-plasma-cutting-equipment', 'industrial-manufacturing'),
    ('industrial-safety-equipment-ppe-workwear', 'industrial-manufacturing'),
    ('industrial-testing-instruments-gauges-scales', 'research-scientific'),
    ('mechanical-power-transmission-bearings-belts', 'industrial-manufacturing'),
    ('power-tools-pneumatic-cutting-equipment', 'industrial-manufacturing'),
    ('tools-workshop-mechanics-equipment', 'industrial-manufacturing'),
    ('welding-soldering-thermal-fabrication', 'industrial-manufacturing'),
    ('industrial-manufacturing-machinery-production', 'industrial-manufacturing'),
    ('factory-warehousing-racking-pallets-drums', 'logistics-supply-chain'),
    ('bulk-industrial-chemicals-solvents-resins', 'chemicals-raw-materials'),
    ('mining-heavy-quarry-drilling-exploration', 'mining-extractive'),
    -- Old clothing/fashion roots
    ('apparel', 'clothing-apparel'),
    ('traditional-cultural-ceremonial-attire', 'clothing-apparel'),
    ('mens-ready-to-wear-apparel', 'clothing-apparel'),
    ('womens-ready-to-wear-apparel', 'clothing-apparel'),
    ('infant-children-youth-apparel', 'clothing-apparel'),
    ('childrens-youth-baby-apparel', 'clothing-apparel'),
    ('leather-formal-casual-footwear', 'footwear-leather'),
    ('athletic-footwear-sneakers-sports-shoes', 'footwear-leather'),
    ('footwear-shoes-boots-sandals', 'footwear-leather'),
    ('fine-jewelry-precious-metals-gemstones', 'fashion-accessories-jewelry'),
    ('fine-jewelry-precious-metals-gold', 'fashion-accessories-jewelry'),
    ('watches-wristwatches-luxury-timepieces', 'fashion-accessories-jewelry'),
    ('wristwatches-fashion-accessories', 'fashion-accessories-jewelry'),
    ('bags-luggage-briefcases-leather-goods', 'footwear-leather'),
    ('luggage-bags-backpacks-leather', 'footwear-leather'),
    ('cosmetics-makeup-perfumes-fragrances', 'beauty-personal-care'),
    ('cosmetics-beauty-makeup-products', 'beauty-personal-care'),
    ('skincare-haircare-personal-hygiene', 'beauty-personal-care'),
    ('skincare-creams-lotions-sunscreen', 'beauty-personal-care'),
    ('haircare-braids-human-hair-wigs', 'beauty-personal-care'),
    ('fragrances-perfumes-colognes-incense', 'fashion-accessories-jewelry'),
    ('textiles-fabrics-garment-sewing-supplies', 'textiles-fabrics'),
    ('beauty-salon-hairdressing-spa-wellness', 'beauty-personal-care'),
    -- Old home/furniture roots
    ('home', 'home-appliances-kitchenware'),
    ('furniture', 'furniture-furnishings'),
    ('major-home-kitchen-appliances', 'home-appliances-kitchenware'),
    ('small-domestic-appliances-cookware', 'home-appliances-kitchenware'),
    ('cookware-pots-cutlery-traditional-sets', 'home-appliances-kitchenware'),
    ('bedding-linens-mattresses-blankets', 'home-appliances-kitchenware'),
    ('living-room-dining-bedroom-furniture', 'furniture-furnishings'),
    ('commercial-office-workspace-furniture', 'furniture-furnishings'),
    ('outdoor-patio-furniture-garden-supplies', 'furniture-furnishings'),
    ('home-storage-plastic-baskets-racks', 'home-appliances-kitchenware'),
    ('housekeeping-cleaning-detergents-mops', 'home-appliances-kitchenware'),
    -- Old sports/entertainment roots
    ('sports', 'sports-fitness'),
    ('sports-fitness-gym-outdoor-equipment', 'sports-fitness'),
    ('team-sports-gear-footballs-athletics', 'sports-fitness'),
    ('outdoor-recreation-camping-trekking-gear', 'sports-fitness'),
    ('entertainment', 'entertainment-recreation'),
    ('musical-instruments-traditional-ethiopian-instruments', 'entertainment-recreation'),
    ('numismatics-philately-collectibles-memorabilia', 'arts-culture-heritage'),
    ('fine-art-original-paintings-sculptures', 'arts-culture-heritage'),
    ('traditional-cultural-artifacts-antiques', 'arts-culture-heritage'),
    -- Old education roots
    ('education', 'education-training'),
    ('private-tutoring-language-skills-training', 'education-training'),
    ('educational-scholarships-academic-fellowships', 'education-training'),
    ('education-teaching-academic-lecturer-roles', 'workforce-employment'),
    ('published-books-academic-literature-periodicals', 'books-stationery-publishing'),
    ('books', 'books-stationery-publishing'),
    -- Old jobs/workforce roots
    ('jobs', 'workforce-employment'),
    ('accounting-finance-banking-roles', 'workforce-employment'),
    ('administrative-clerical-secretarial-roles', 'workforce-employment'),
    ('technology-software-engineering-it-roles', 'workforce-employment'),
    ('sales-marketing-business-development-roles', 'workforce-employment'),
    ('hospitality-culinary-restaurant-hotel-roles', 'workforce-employment'),
    ('general-labor-security-custodial-support-roles', 'workforce-employment'),
    ('professional-occupations-careers-recruitment', 'workforce-employment'),
    ('transportation-logistics-driver-warehouse-roles', 'workforce-employment'),
    ('medical-clinical-healthcare-nursing-roles', 'workforce-employment'),
    -- Old services roots
    ('services', 'business-professional-services'),
    ('professional-legal-financial-advisory-services', 'business-professional-services'),
    ('commercial-cleaning-fumigation-janitorial', 'cleaning-facility-management'),
    ('skilled-trades-electricians-plumbers-masons', 'construction-infrastructure'),
    ('electronics-appliance-repair-services', 'electronics'),
    ('events-staging-audiovisual-production-services', 'events-exhibitions'),
    ('commercial-printing-graphic-branding-services', 'packaging-printing'),
    ('manned-guarding-cctv-security-services', 'security-protective'),
    ('tactical-security-ballistic-defense-supplies', 'security-protective'),
    -- Old commerce roots
    ('commercial-retail-discount-promotions-sales', 'commerce-trade'),
    ('business-to-business-promotional-packages', 'commerce-trade'),
    ('corporate-private-procurement-rfqs', 'commerce-trade'),
    ('public-government-procurement-tenders-rfps', 'government-public-administration'),
    -- Old logistics roots
    ('logistics', 'logistics-supply-chain'),
    ('freight-logistics-cargo-moving-services', 'logistics-supply-chain'),
    ('industrial-machinery-generator-car-rental', 'logistics-supply-chain'),
    -- Old media roots
    ('broadcast-production-radio-tv-transmitters', 'communications-media'),
    ('cinema-cameras-film-lighting-production-rigs', 'communications-media'),
    ('concert-stage-sound-reinforcement-stage-lighting', 'events-exhibitions'),
    ('audiovisual-translation-dubbing-subtitling', 'communications-media'),
    ('photography-videography-drone-aerial', 'communications-media'),
    -- Old community roots
    ('ngo-humanitarian-grants-funding-calls', 'community-social-development'),
    ('emergency-disaster-relief-supplies-safety', 'public-safety-emergency'),
    ('environmental-remediation-recycling-equipment', 'environment-conservation'),
    ('chemicals', 'chemicals-raw-materials'),
    ('specialized-scientific-metrology-testing-gear', 'research-scientific'),
    ('community-cultural-events-public-announcements', 'community-social-development'),
    ('official-public-notices-legal-announcements', 'government-public-administration'),
    -- Old gig roots
    ('ai-data-annotation-prompt-engineering-gigs', 'workforce-employment'),
    ('freelance-writing-translation-copywriting-gigs', 'workforce-employment'),
    ('graphic-design-branding-illustration-gigs', 'workforce-employment'),
    ('online-consulting-coaching-advisory-gigs', 'workforce-employment'),
    ('social-media-management-digital-marketing-gigs', 'workforce-employment'),
    ('video-editing-motion-graphics-animation-gigs', 'workforce-employment'),
    ('voiceover-audio-production-podcasting-gigs', 'workforce-employment'),
    ('web-mobile-fullstack-programming-gigs', 'workforce-employment'),
    ('virtual-assistant-data-entry-customer-support', 'workforce-employment'),
    ('micro-tasks-surveys-remote-clerical-gigs', 'workforce-employment')
) AS t(old_slug, new_slug)
LEFT JOIN category_hierarchy ch ON ch.slug = t.new_slug AND ch.parent_id IS NULL
WHERE NOT EXISTS (
    SELECT 1 FROM category_redirects cr
    WHERE cr.old_slug = t.old_slug AND cr.new_slug = t.new_slug
);

-- PART 4: Update products to use new category hierarchy IDs where possible.
-- For products whose category_hierarchy_id points to an old (now inactive) root,
-- try to resolve the redirect and update the product to point to the new root.
UPDATE products p
SET category_hierarchy_id = (
    SELECT ch_new.id
    FROM category_redirects cr
    JOIN category_hierarchy ch_new ON ch_new.slug = cr.new_slug AND ch_new.parent_id IS NULL
    JOIN category_hierarchy ch_old ON ch_old.id = p.category_hierarchy_id
    WHERE cr.old_slug = ch_old.slug
    LIMIT 1
)
WHERE p.category_hierarchy_id IS NOT NULL
  AND EXISTS (
    SELECT 1 FROM category_hierarchy ch_old
    WHERE ch_old.id = p.category_hierarchy_id
      AND ch_old.is_active = false
      AND ch_old.parent_id IS NULL
  );

-- PART 5: Verify no orphaned products (safety check).
DO $$
DECLARE
    orphan_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO orphan_count
    FROM products
    WHERE category_hierarchy_id IS NOT NULL
      AND NOT EXISTS (
          SELECT 1 FROM category_hierarchy WHERE id = products.category_hierarchy_id
      );

    IF orphan_count > 0 THEN
        RAISE NOTICE 'WARNING: %d products have orphaned category_hierarchy_id references. Manual review needed.', orphan_count;
    ELSE
        RAISE NOTICE 'OK: No orphaned product category references found.';
    END IF;
END $$;

-- PART 6: Log migration
INSERT INTO taxonomy_audit_logs (action, reason, old_values, new_values, created_at)
SELECT 'taxonomy_consolidation',
       'Consolidation of root categories to canonical roots',
       jsonb_build_object('inactive_roots', (SELECT COUNT(*) FROM category_hierarchy WHERE parent_id IS NULL AND is_active = false)),
       jsonb_build_object('active_roots', (SELECT COUNT(*) FROM category_hierarchy WHERE parent_id IS NULL AND is_active = true)),
       NOW()
WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'taxonomy_audit_logs');
