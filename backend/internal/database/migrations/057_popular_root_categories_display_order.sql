-- 057_popular_root_categories_display_order.sql
-- Prioritize high-demand consumer and commercial marketplace categories first

UPDATE category_hierarchy SET display_order = 50 WHERE level = 0 AND slug IN (
    'aerospace-defence', 'banking-financial-services', 'freelance-gig-services', 
    'nonprofit-development-organizations', 'pharmaceuticals-drug-industry', 
    'telecommunications', 'tenders-business-opportunities', 'tobacco-regulated-substances'
);

UPDATE category_hierarchy SET display_order = 1 WHERE level = 0 AND slug IN ('electronics');
UPDATE category_hierarchy SET display_order = 2 WHERE level = 0 AND slug IN ('vehicles', 'transport-mobility');
UPDATE category_hierarchy SET display_order = 3 WHERE level = 0 AND slug IN ('property', 'property-real-estate');
UPDATE category_hierarchy SET display_order = 4 WHERE level = 0 AND slug IN ('business', 'construction-infrastructure');
UPDATE category_hierarchy SET display_order = 5 WHERE level = 0 AND slug IN ('agriculture');
UPDATE category_hierarchy SET display_order = 6 WHERE level = 0 AND slug IN ('clothing-apparel', 'fashion', 'footwear-leather');
UPDATE category_hierarchy SET display_order = 7 WHERE level = 0 AND slug IN ('home-appliances-kitchenware', 'furniture-furnishings', 'home');
UPDATE category_hierarchy SET display_order = 8 WHERE level = 0 AND slug IN ('computing-it');
UPDATE category_hierarchy SET display_order = 9 WHERE level = 0 AND slug IN ('food-beverage');
UPDATE category_hierarchy SET display_order = 10 WHERE level = 0 AND slug IN ('industrial-manufacturing');
