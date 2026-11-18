#!/usr/bin/env python3
"""
Category Taxonomy Rebuilder
Rebuilds all category JSON files with:
1. Single, clean canonical names (no "&" or mixed terms)
2. Alphabetical ordering
3. Synonym/alias system for search
4. Vertical hierarchy
"""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Set

# Category name cleanup and synonym extraction
def extract_synonyms(combined_name: str) -> tuple[str, List[str]]:
    """
    Extract canonical name and synonyms from combined names
    
    Examples:
    "Cameras & Photography" -> ("Cameras", ["Photography", "Photo", "Camera Equipment"])
    "Audio Equipment & Headphones" -> ("Audio Equipment", ["Headphones", "Audio", "Sound"])
    "TV, Audio & Video" -> ("Television", ["TV", "Audio", "Video", "Home Theater"])
    """
    
    # Synonym mapping for common terms
    synonym_map = {
        "cameras & photography": ("Cameras", ["Photography", "Photo Equipment", "Camera Gear", "Imaging"]),
        "cameras & photography equipment": ("Cameras", ["Photography", "Photo Equipment", "Camera Gear"]),
        "audio equipment & headphones": ("Audio Equipment", ["Headphones", "Sound Systems", "Audio Gear"]),
        "computer peripherals & accessories": ("Computer Peripherals", ["PC Accessories", "Computer Accessories"]),
        "mobile accessories & power banks": ("Mobile Accessories", ["Phone Accessories", "Power Banks", "Chargers"]),
        "smart home & automation devices": ("Smart Home", ["Home Automation", "IoT Devices", "Connected Home"]),
        "smartphones & mobile devices": ("Smartphones", ["Mobile Phones", "Cell Phones", "Handsets"]),
        "tablets & e-readers": ("Tablets", ["E-Readers", "Digital Readers", "Tablet Computers"]),
        "television & displays": ("Televisions", ["TV", "Displays", "Monitors", "Screens"]),
        "wearable technology & smartwatches": ("Wearables", ["Smartwatches", "Fitness Trackers", "Smart Jewelry"]),
        "gaming consoles & esports gear": ("Gaming Consoles", ["Esports", "Gaming Equipment", "Game Systems"]),
        "repair & maintenance services": ("Repair Services", ["Maintenance", "Fixing", "Technical Support"]),
        "trade-in & equipment exchange": ("Trade-In Services", ["Equipment Exchange", "Device Trade", "Buyback"]),
        "warranty, protection & overhaul": ("Warranty Services", ["Protection Plans", "Extended Warranty"]),
        "tv, audio & video": ("Home Entertainment", ["TV", "Audio", "Video", "Home Theater"]),
        "phones & tablets": ("Mobile Devices", ["Phones", "Tablets", "Smartphones"]),
        "computers & laptops": ("Computers", ["Laptops", "Desktops", "PCs", "Notebooks"]),
    }
    
    lower_name = combined_name.lower()
    
    # Check direct mapping
    if lower_name in synonym_map:
        return synonym_map[lower_name]
    
    # Auto-extract from "&" patterns
    if " & " in combined_name:
        parts = [p.strip() for p in combined_name.split("&")]
        canonical = parts[0]  # First part is primary
        synonyms = parts[1:] + [combined_name]  # Rest are synonyms
        return (canonical, synonyms)
    
    # Auto-extract from "," patterns
    if ", " in combined_name:
        parts = [p.strip() for p in combined_name.split(",")]
        canonical = parts[0]
        synonyms = parts[1:] + [combined_name]
        return (canonical, synonyms)
    
    # Single word - no extraction needed
    return (combined_name, [])


def clean_category_name(name: str) -> str:
    """Remove unnecessary words and clean up category names"""
    # Remove common suffixes
    name = re.sub(r'\s+(Equipment|Devices|Services|Gear|Products)$', '', name, flags=re.IGNORECASE)
    
    # But keep if it's the whole name
    if name.strip() in ["Equipment", "Devices", "Services", "Gear", "Products"]:
        return name
    
    return name.strip()


def generate_slug(name: str) -> str:
    """Generate URL-friendly slug from name"""
    slug = name.lower()
    slug = re.sub(r'[^\w\s-]', '', slug)  # Remove special chars
    slug = re.sub(r'[\s_]+', '-', slug)   # Replace spaces with hyphens
    slug = re.sub(r'-+', '-', slug)        # Remove duplicate hyphens
    return slug.strip('-')


def rebuild_electronics_categories():
    """Rebuild electronics.json with clean names and synonyms"""
    
    electronics_structure = {
        "classification_system": "SYS_PROD",
        "family_slug": "electronics",
        "categories": []
    }
    
    # Root category
    root = {
        "id": 1,
        "parent_id": None,
        "parent_code": None,
        "classification_system_code": "SYS_PROD",
        "slug": "electronics",
        "names": {
            "en": "Electronics",
            "am": "ኤሌክትሮኒክስ",
            "om": "Elektirooniksii",
            "ti": "ኤሌክትሮኒክስ"
        },
        "description": "Electronic devices and equipment",
        "level": 0,
        "path": "/1/",
        "display_order": 1,
        "icon": "Smartphone",
        "image_url": None,
        "is_active": True,
        "is_leaf": False,
        "is_archived": False,
        "meta_title": None,
        "meta_description": None,
        "aliases": ["Consumer Electronics", "Electronic Goods", "Tech Products"],
        "redirects": []
    }
    
    electronics_structure["categories"].append(root)
    
    # Level 1 categories - Alphabetically ordered with clean names
    level1_categories = [
        {
            "name": "Accessories",
            "aliases": ["Electronic Accessories", "Tech Accessories", "Add-ons"],
            "icon": "Package"
        },
        {
            "name": "Audio Equipment",
            "aliases": ["Sound Systems", "Audio Gear", "Speakers", "Headphones"],
            "icon": "Volume2"
        },
        {
            "name": "Cameras",
            "aliases": ["Photography Equipment", "Photo Gear", "Imaging Devices", "Camera Equipment"],
            "icon": "Camera"
        },
        {
            "name": "Computer Peripherals",
            "aliases": ["PC Accessories", "Computer Accessories", "Input Devices"],
            "icon": "Monitor"
        },
        {
            "name": "Computers",
            "aliases": ["Laptops", "Desktops", "PCs", "Personal Computers", "Notebooks"],
            "icon": "Laptop"
        },
        {
            "name": "Gaming",
            "aliases": ["Gaming Consoles", "Video Games", "Esports Equipment", "Game Systems"],
            "icon": "Gamepad2"
        },
        {
            "name": "Mobile Accessories",
            "aliases": ["Phone Accessories", "Smartphone Accessories", "Power Banks", "Chargers"],
            "icon": "Zap"
        },
        {
            "name": "Smartphones",
            "aliases": ["Mobile Phones", "Cell Phones", "Handsets", "Mobile Devices"],
            "icon": "Smartphone"
        },
        {
            "name": "Smart Home",
            "aliases": ["Home Automation", "IoT Devices", "Connected Home", "Smart Devices"],
            "icon": "Home"
        },
        {
            "name": "Tablets",
            "aliases": ["E-Readers", "Tablet Computers", "Digital Readers", "Slates"],
            "icon": "Tablet"
        },
        {
            "name": "Televisions",
            "aliases": ["TV", "TVs", "Displays", "Home Entertainment", "Screens"],
            "icon": "Tv"
        },
        {
            "name": "Wearables",
            "aliases": ["Smartwatches", "Fitness Trackers", "Smart Jewelry", "Wearable Tech"],
            "icon": "Watch"
        }
    ]
    
    # Add level 1 categories
    for idx, cat_def in enumerate(level1_categories, start=2):
        cat = {
            "id": idx * 100,  # Use spaced IDs for easy insertion
            "parent_id": 1,
            "parent_code": "electronics",
            "classification_system_code": "SYS_PROD",
            "slug": generate_slug(cat_def["name"]),
            "names": {
                "en": cat_def["name"],
                "am": cat_def["name"],  # TODO: Add proper translations
                "om": cat_def["name"],
                "ti": cat_def["name"]
            },
            "description": f"{cat_def['name']} products and equipment",
            "level": 1,
            "path": f"/1/{idx * 100}/",
            "display_order": idx,
            "icon": cat_def.get("icon", "Box"),
            "image_url": None,
            "is_active": True,
            "is_leaf": False,
            "is_archived": False,
            "meta_title": None,
            "meta_description": None,
            "aliases": cat_def.get("aliases", []),
            "redirects": []
        }
        electronics_structure["categories"].append(cat)
    
    # Add Level 2 - Audio Equipment subcategories
    audio_subs = [
        {"name": "Bluetooth Speakers", "aliases": ["Wireless Speakers", "Portable Speakers"]},
        {"name": "Earbuds", "aliases": ["In-Ear Headphones", "Wireless Earbuds", "True Wireless"]},
        {"name": "Headphones", "aliases": ["Over-Ear Headphones", "Studio Headphones"]},
        {"name": "Home Theater", "aliases": ["Surround Sound", "Theater Systems"]},
        {"name": "Microphones", "aliases": ["Mics", "Recording Equipment"]},
        {"name": "Soundbars", "aliases": ["TV Speakers", "Sound Bars"]},
    ]
    
    for idx, sub in enumerate(audio_subs, start=1):
        cat = {
            "id": 200 + idx,
            "parent_id": 200,
            "parent_code": "audio-equipment",
            "classification_system_code": "SYS_PROD",
            "slug": generate_slug(sub["name"]),
            "names": {"en": sub["name"], "am": sub["name"], "om": sub["name"], "ti": sub["name"]},
            "description": f"{sub['name']} products",
            "level": 2,
            "path": f"/1/200/{200 + idx}/",
            "display_order": idx,
            "icon": "",
            "image_url": None,
            "is_active": True,
            "is_leaf": True,
            "is_archived": False,
            "meta_title": None,
            "meta_description": None,
            "aliases": sub.get("aliases", []),
            "redirects": []
        }
        electronics_structure["categories"].append(cat)
    
    # Add Level 2 - Cameras subcategories
    camera_subs = [
        {"name": "Action Cameras", "aliases": ["Sports Cameras", "Go Pros", "Adventure Cameras"]},
        {"name": "Camera Lenses", "aliases": ["Lenses", "Optical Lenses", "Camera Glass"]},
        {"name": "DSLR Cameras", "aliases": ["Digital SLR", "Reflex Cameras"]},
        {"name": "Instant Cameras", "aliases": ["Polaroid", "Instant Print"]},
        {"name": "Mirrorless Cameras", "aliases": ["Hybrid Cameras", "Compact System"]},
        {"name": "Point and Shoot", "aliases": ["Compact Cameras", "Digital Cameras"]},
        {"name": "Video Cameras", "aliases": ["Camcorders", "Recording Cameras"]},
    ]
    
    for idx, sub in enumerate(camera_subs, start=1):
        cat = {
            "id": 400 + idx,
            "parent_id": 400,
            "parent_code": "cameras",
            "classification_system_code": "SYS_PROD",
            "slug": generate_slug(sub["name"]),
            "names": {"en": sub["name"], "am": sub["name"], "om": sub["name"], "ti": sub["name"]},
            "description": f"{sub['name']} products",
            "level": 2,
            "path": f"/1/400/{400 + idx}/",
            "display_order": idx,
            "icon": "",
            "image_url": None,
            "is_active": True,
            "is_leaf": True,
            "is_archived": False,
            "meta_title": None,
            "meta_description": None,
            "aliases": sub.get("aliases", []),
            "redirects": []
        }
        electronics_structure["categories"].append(cat)
    
    # Add Level 2 - Computers subcategories
    computer_subs = [
        {"name": "All-in-One PCs", "aliases": ["AIO Computers", "All-in-One Desktops"]},
        {"name": "Desktop Computers", "aliases": ["Desktop PCs", "Tower PCs"]},
        {"name": "Gaming Laptops", "aliases": ["Gaming Notebooks", "High-Performance Laptops"]},
        {"name": "Laptops", "aliases": ["Notebooks", "Portable Computers"]},
        {"name": "Ultrabooks", "aliases": ["Thin Laptops", "Slim Notebooks"]},
        {"name": "Workstations", "aliases": ["Professional PCs", "High-End Computers"]},
    ]
    
    for idx, sub in enumerate(computer_subs, start=1):
        cat = {
            "id": 500 + idx,
            "parent_id": 500,
            "parent_code": "computers",
            "classification_system_code": "SYS_PROD",
            "slug": generate_slug(sub["name"]),
            "names": {"en": sub["name"], "am": sub["name"], "om": sub["name"], "ti": sub["name"]},
            "description": f"{sub['name']} products",
            "level": 2,
            "path": f"/1/500/{500 + idx}/",
            "display_order": idx,
            "icon": "",
            "image_url": None,
            "is_active": True,
            "is_leaf": True,
            "is_archived": False,
            "meta_title": None,
            "meta_description": None,
            "aliases": sub.get("aliases", []),
            "redirects": []
        }
        electronics_structure["categories"].append(cat)
    
    return electronics_structure


def main():
    """Main execution"""
    print("=" * 60)
    print("Category Taxonomy Rebuilder")
    print("=" * 60)
    
    # Rebuild electronics
    print("\n[1/1] Rebuilding Electronics taxonomy...")
    electronics = rebuild_electronics_categories()
    
    # Save to file
    output_file = Path("/home/poven/GO PROJECT/sellbuy/data/master/categories/classification-systems/SYS_PROD/families/electronics_v2.json")
    output_file.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(electronics, f, indent=2, ensure_ascii=False)
    
    print(f"✓ Saved to: {output_file}")
    print(f"  Total categories: {len(electronics['categories'])}")
    print(f"  Root: 1")
    print(f"  Level 1: {len([c for c in electronics['categories'] if c['level'] == 1])}")
    print(f"  Level 2: {len([c for c in electronics['categories'] if c['level'] == 2])}")
    
    # Count aliases
    total_aliases = sum(len(c.get('aliases', [])) for c in electronics['categories'])
    print(f"  Total aliases: {total_aliases}")
    
    print("\n" + "=" * 60)
    print("✓ Rebuild complete!")
    print("=" * 60)
    print("\nNext steps:")
    print("1. Review electronics_v2.json")
    print("2. Backup original: mv electronics.json electronics_old.json")
    print("3. Replace: mv electronics_v2.json electronics.json")
    print("4. Import: go run cmd/tools/masterdata/main.go -action=import")


if __name__ == "__main__":
    main()
