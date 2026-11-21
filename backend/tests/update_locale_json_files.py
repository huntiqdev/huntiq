#!/usr/bin/env python3
import json

LOCALES = ['en', 'am', 'or', 'ti']
FILE_PATHS = {loc: f"frontend/src/lib/language/{loc}.json" for loc in LOCALES}

# New dictionaries per locale
ADDITIONS = {
    "footer": {
        "verified_platform": {
            "en": "Verified Platform",
            "am": "የተረጋገጠ መድረክ",
            "or": "Waltajjii Mirkanaa'e",
            "ti": "ዝተረጋገጸ መድረኽ"
        },
        "customer_support": {
            "en": "Customer Support",
            "am": "የደንበኞች አገልግሎት",
            "or": "Tajaajila Maamiltootaa",
            "ti": "ኣገልግሎት ዓማዊል"
        },
        "faq_title": {
            "en": "Frequently Asked Questions",
            "am": "ተደጋግመው የሚጠየቁ ጥያቄዎች",
            "or": "Gaaffilee Yeroo Baay'ee Gaafataman",
            "ti": "ተደጋጋሚ ሕቶታት"
        },
        "support_ticket": {
            "en": "Support Ticket",
            "am": "የድጋፍ ጥያቄ",
            "or": "Waraqaa Deeggarsaa",
            "ti": "ናይ ደገፍ ቲኬት"
        },
        "buying": {
            "en": "Buying",
            "am": "መግዛት",
            "or": "Bittaa",
            "ti": "ምግዛእ"
        },
        "all_products_sale": {
            "en": "All Products for Sale",
            "am": "ለሽያጭ የቀረቡ ሁሉም እቃዎች",
            "or": "Meeshaalee Gurgurtaaf Dhiyaatan Hunda",
            "ti": "ንመሸጣ ዝቐረቡ ኩሎም ኣቑሑ"
        },
        "brand_new_items": {
            "en": "Brand New Items",
            "am": "አዳዲስ እቃዎች",
            "or": "Meeshaalee Haaraa",
            "ti": "ሓደስቲ ኣቑሑ"
        },
        "pre_owned_items": {
            "en": "Quality Pre-Owned",
            "am": "ጥራት ያላቸው ያገለገሉ",
            "or": "Meeshaalee Qulqullina Qaban Itti Fayyadaman",
            "ti": "ጽሬት ዘለዎም ዝተጠቕሙ"
        },
        "land_property": {
            "en": "Land & Property",
            "am": "መሬት እና ንብረት",
            "or": "Lafa fi Qabeenya",
            "ti": "መሬትን ንብረትን"
        },
        "buyer_protection": {
            "en": "Buyer Protection",
            "am": "የገዢ ጥበቃ",
            "or": "Eegumsa Bittootaa",
            "ti": "ዕቝባ ዓዳጊ"
        },
        "selling_renting": {
            "en": "Selling & Renting",
            "am": "መሸጥ እና ማከራየት",
            "or": "Gurgurtaa fi Kiraa",
            "ti": "ምሻጥን ምኽራይን"
        },
        "start_selling": {
            "en": "Start Selling",
            "am": "መሸጥ ይጀምሩ",
            "or": "Gurguruu Jalqabaa",
            "ti": "ምሻጥ ጀምሩ"
        },
        "list_rental_property": {
            "en": "List Rental / Property",
            "am": "ኪራይ / ንብረት ይመዝግቡ",
            "or": "Kiraa / Qabeenya Galmeessaa",
            "ti": "ክራይ / ንብረት መዝግቡ"
        },
        "seller_guidelines": {
            "en": "Seller Requirements",
            "am": "የሻጮች መመሪያ",
            "or": "Qajeelfama Gurgurtootaa",
            "ti": "መምርሒ ነጋዶ"
        },
        "seller_dashboard": {
            "en": "Seller Dashboard",
            "am": "የሻጭ ዳሽቦርድ",
            "or": "Daashboordii Gurgurtaa",
            "ti": "ዳሽቦርድ ሸቃጢ"
        },
        "gold_premium_tiers": {
            "en": "Gold & Premium Tiers",
            "am": "የወርቅ እና ፕሪሚየም ደረጃዎች",
            "or": "Sadarkaa Warqii fi Piriimiyeemii",
            "ti": "ደረጃታት ወርቅን ፕሪምየምን"
        },
        "policies_language": {
            "en": "Policies & Language",
            "am": "መመሪያዎች እና ቋንቋ",
            "or": "Imaammataa fi Afaan",
            "ti": "ፖሊሲታትን ቋንቋን"
        },
        "safety_guidelines": {
            "en": "Safety Guidelines",
            "am": "የደህንነት መመሪያዎች",
            "or": "Qajeelfama Nageenyaa",
            "ti": "መምርሒታት ድሕንነት"
        },
        "protected_transactions": {
            "en": "Protected Transactions",
            "am": "ደህንነታቸው የተጠበቁ ግብይቶች",
            "or": "Daldala Eegumsa Qabu",
            "ti": "ድሑናት ንግዳዊ ልውውጣት"
        },
        "ethiopian_languages": {
            "en": "4 Ethiopian Languages",
            "am": "4 የኢትዮጵያ ቋንቋዎች",
            "or": "Afaanota Itoophiyaa 4",
            "ti": "4 ናይ ኢትዮጵያ ቋንቋታት"
        }
    },
    "seller_nav": {
        "landlord_console": {
            "en": "Landlord Console",
            "am": "የአከራይ ኮንሶል",
            "or": "Koonsolii Abbaa Qabeenyaa",
            "ti": "ኮንሶል ኣካራዪ"
        },
        "seller_studio": {
            "en": "Seller Studio",
            "am": "የሻጭ ስቱዲዮ",
            "or": "Istuudiyoo Gurgurtaa",
            "ti": "ስቱዲዮ ሸቃጢ"
        },
        "provider_studio": {
            "en": "Provider Studio",
            "am": "የአገልግሎት ሰጪ ስቱዲዮ",
            "or": "Istuudiyoo Dhiyeessaa",
            "ti": "ስቱዲዮ ኣቕራቢ"
        },
        "properties_rentals": {
            "en": "Properties & Rentals",
            "am": "ንብረቶች እና ኪራዮች",
            "or": "Qabeenyaa fi Kiraa",
            "ti": "ንብረታትን ክራያትን"
        },
        "lease_requests": {
            "en": "Lease Requests",
            "am": "የኪራይ ጥያቄዎች",
            "or": "Gaaffilee Kiraa",
            "ti": "ናይ ክራይ ሕቶታት"
        },
        "rental_income": {
            "en": "Rental Income",
            "am": "የኪራይ ገቢ",
            "or": "Galii Kiraa",
            "ti": "ናይ ክራይ ኣታዊ"
        },
        "back_to_marketplace": {
            "en": "Back to Marketplace",
            "am": "ወደ ገበያ ይመለሱ",
            "or": "Gara Gabaatti Deebi'aa",
            "ti": "ናብ ዕዳጋ ተመለሱ"
        },
        "sign_out": {
            "en": "Sign Out",
            "am": "ውጣ",
            "or": "Ba'i",
            "ti": "ውጻእ"
        },
        "choose_console": {
            "en": "Choose Your Provider Console",
            "am": "የአገልግሎት ሰጪ ኮንሶልዎን ይምረጡ",
            "or": "Koonsolii Dhiyeessaa Keessan Filadhaa",
            "ti": "ናይ ኣቕራቢ ኮንሶልኩም ምረጹ"
        },
        "active": {
            "en": "Active",
            "am": "ንቁ",
            "or": "Hojii irra",
            "ti": "ንጡፍ"
        },
        "activate": {
            "en": "Activate",
            "am": "አግብር",
            "or": "Hojii irra oolchi",
            "ti": "ኣንቅሕ"
        }
    },
    "buyer_nav": {
        "buyer_panel": {
            "en": "Buyer Panel",
            "am": "የገዢ ፓነል",
            "or": "Paanaalii Bittaa",
            "ti": "ፓነል ዓዳጊ"
        },
        "dashboard": {
            "en": "Dashboard",
            "am": "ዳሽቦርድ",
            "or": "Daashboordii",
            "ti": "ዳሽቦርድ"
        },
        "my_inquiries": {
            "en": "My Inquiries",
            "am": "የእኔ ጥያቄዎች",
            "or": "Gaaffilee Koo",
            "ti": "ሕቶታተይ"
        },
        "referral": {
            "en": "Referral",
            "am": "የሪፈራል ፕሮግራም",
            "or": "Wamicha",
            "ti": "ሪፈራል"
        },
        "profile": {
            "en": "Profile",
            "am": "መገለጫ",
            "or": "Piroofayilii",
            "ti": "መግለጺ"
        },
        "marketplace_capabilities": {
            "en": "Marketplace Capabilities",
            "am": "የገበያ አገልግሎት አማራጮች",
            "or": "Dandeettii Gabaa",
            "ti": "ናይ ዕዳጋ ዓቕምታት"
        },
        "start_selling": {
            "en": "Start Selling",
            "am": "መሸጥ ይጀምሩ",
            "or": "Gurguruu Jalqabaa",
            "ti": "ምሻጥ ጀምሩ"
        },
        "list_for_rent": {
            "en": "List for Rent",
            "am": "ለኪራይ ያቅርቡ",
            "or": "Kiraaf Dhiyeessaa",
            "ti": "ንክራይ ኣቕርቡ"
        },
        "upgrade_to_unlock": {
            "en": "Upgrade to unlock",
            "am": "ለመክፈት ደረጃ ያሳድጉ",
            "or": "Banuuf sadarkaa ol guddisaa",
            "ti": "ንምኽፋት ደረጃ ኣዕብዩ"
        },
        "upgrade_to_premium": {
            "en": "Upgrade to Premium",
            "am": "ወደ ፕሪሚየም ደረጃ ያሳድጉ",
            "or": "Gara Piriimiyeemiitti Ol Guddisaa",
            "ti": "ናብ ፕሪምየም ደረጃ ኣዕብዩ"
        }
    },
    "dashboard": {
        "fee_summary": {
            "en": "Fee Summary",
            "am": "የክፍያ ማጠቃለያ",
            "or": "Cuunfaa Kaffaltii",
            "ti": "መጠቃለሊ ክፍሊት"
        },
        "total_fees_paid": {
            "en": "Total Fees Paid",
            "am": "ጠቅላላ የተከፈለ ክፍያ",
            "or": "Kaffaltii Waliigalaa Kaffalame",
            "ti": "ጠቕላላ ዝተኸፍለ ክፍሊት"
        },
        "net_profit": {
            "en": "Net Profit",
            "am": "የተጣራ ትርፍ",
            "or": "Bu'aa Qulqulluu",
            "ti": "ጽሩይ መኽሰብ"
        },
        "unpaid_listing_fees": {
            "en": "Unpaid Listing Fees",
            "am": "ያልተከፈሉ የማስታወቂያ ክፍያዎች",
            "or": "Kaffaltii Beeksisaa Hin Kaffalmin",
            "ti": "ዘይተኸፍሉ ናይ ምዝገባ ክፍሊታት"
        },
        "no_listings_yet": {
            "en": "No listings yet",
            "am": "እስካሁን ምንም ማስታወቂያ የለም",
            "or": "Ammas beeksisni hin jiru",
            "ti": "ክሳብ ሕጂ ዝኾነ ምልክታ የለን"
        },
        "active_listings_remaining": {
            "en": "Active Listings — Time Remaining",
            "am": "ንቁ ማስታወቂያዎች — የቀረው ጊዜ",
            "or": "Beeksisa Hojii Irra Jiru — Yeroo Hafe",
            "ti": "ንጡፋት ምልክታታት — ዝተረፈ ግዜ"
        },
        "no_expiry": {
            "en": "No expiry",
            "am": "ጊዜው አያበቃም",
            "or": "Yeroon hin dhumu",
            "ti": "ግዜኡ ዘየብቅዕ"
        },
        "recent_orders": {
            "en": "Recent Order Requests",
            "am": "የቅርብ ጊዜ የትዕዛዝ ጥያቄዎች",
            "or": "Gaaffilee Ajajaa Dhihoo",
            "ti": "ናይ ቀረባ ግዜ ናይ ትእዛዝ ሕቶታት"
        },
        "seller_tier": {
            "en": "Seller Tier",
            "am": "የሻጭ ደረጃ",
            "or": "Sadarkaa Gurgurtaa",
            "ti": "ደረጃ ሸቃጢ"
        },
        "value_sold": {
            "en": "Value Sold",
            "am": "የተሸጠ ዋጋ",
            "or": "Gatii Gurgurame",
            "ti": "ዝተሸጠ ዋጋ"
        },
        "locked_features": {
            "en": "Locked features",
            "am": "የተዘጉ አገልግሎቶች",
            "or": "Tajaajila Cufame",
            "ti": "ዝተዓጸዉ ኣገልግሎታት"
        },
        "max_tier_unlocked": {
            "en": "Maximum tier — all features unlocked",
            "am": "ከፍተኛ ደረጃ — ሁሉም አገልግሎቶች ተከፍተዋል",
            "or": "Sadarkaa olaanaa — tajaajilli hundi banameera",
            "ti": "ዝለዓለ ደረጃ — ኩሎም ኣገልግሎታት ተኸፊቶም"
        },
        "customer_access_tier": {
            "en": "Customer Access Tier",
            "am": "የደንበኛ መዳረሻ ደረጃ",
            "or": "Sadarkaa Seensa Maamilaa",
            "ti": "ደረጃ መእተዊ ዓማዊል"
        },
        "my_purchases_orders": {
            "en": "My Purchases & Orders",
            "am": "የእኔ ግዢዎች እና ትዕዛዞች",
            "or": "Bittaa fi Ajajoota Koo",
            "ti": "ግዝአተይን ትእዛዛተይን"
        },
        "browse_rentals": {
            "en": "Browse Rental Properties",
            "am": "የኪራይ ንብረቶችን ይመልከቱ",
            "or": "Qabeenya Kiraa Daawwadhaa",
            "ti": "ናይ ክራይ ንብረታት ርኣዩ"
        },
        "leases_landlords": {
            "en": "Leases & Landlord listings",
            "am": "ውሎች እና የአከራይ ማስታወቂያዎች",
            "or": "Waliigaltee fi Beeksisa Abbaa Qabeenyaa",
            "ti": "ውዕላትን ናይ ኣካራዪ ምልክታታትን"
        },
        "recent_requests": {
            "en": "Recent Requests",
            "am": "የቅርብ ጊዜ ጥያቄዎች",
            "or": "Gaaffilee Dhihoo",
            "ti": "ናይ ቀረባ ግዜ ሕቶታት"
        },
        "no_requests_yet": {
            "en": "No order requests yet",
            "am": "እስካሁን ምንም የትዕዛዝ ጥያቄ የለም",
            "or": "Ammas gaaffiin ajajaa hin jiru",
            "ti": "ክሳብ ሕጂ ዝኾነ ናይ ትእዛዝ ሕቶ የለን"
        }
    },
    "buyer_tier_features": {
        "gold_buyer": {
            "en": "Gold Buyer",
            "am": "የወርቅ ገዢ",
            "or": "Bittaa Warqii",
            "ti": "ወርቂ ዓዳጊ"
        },
        "verified_trader": {
            "en": "Verified Trader",
            "am": "የተረጋገጠ ነጋዴ",
            "or": "Daldalaa Mirkanaa'e",
            "ti": "ዝተረጋገጸ ነጋዳይ"
        },
        "gold_benefit_contact": {
            "en": "Direct verified seller & landlord phone contacts",
            "am": "ቀጥተኛ የተረጋገጡ የሻጭ እና የአከራይ ስልክ አድራሻዎች",
            "or": "Bilbila gurgurtootaa fi abbootii qabeenyaa mirkanaa'an kallattiin",
            "ti": "ቀጥታዊ ናይ ዝተረጋገጹ ነጋዶን ኣካረይትን ስልኪ ርክባት"
        },
        "gold_benefit_sms": {
            "en": "Priority SMS alerts for discounted listings",
            "am": "ለተቀነሱ ማስታወቂያዎች የቅድሚያ የኤስኤምኤስ ማንቂያዎች",
            "or": "Beeksisa gatiin gad bu'eef akeekkachiisa SMS dursaa",
            "ti": "ዋጋ ንዝተነከየሎም ምልክታታት ናይ ቅድሚት SMS ምልክታ"
        },
        "gold_benefit_badge": {
            "en": "Official Gold Buyer Verified Badge",
            "am": "ይፋዊ የወርቅ ገዢ የተረጋገጠ ባጅ",
            "or": "Asxaa Mirkanaa'aa Bittaa Warqii Seera Qabeessa",
            "ti": "ወግዓዊ ናይ ወርቂ ዓዳጊ ዝተረጋገጸ ባጅ"
        },
        "premium_vip": {
            "en": "Premium VIP",
            "am": "ፕሪሚየም ቪአይፒ",
            "or": "Piriimiyeemii VIP",
            "ti": "ፕሪምየም ቪኣይፒ"
        },
        "maximum_privilege": {
            "en": "Maximum Privilege",
            "am": "ከፍተኛ መብት",
            "or": "Mirga Olaanaa",
            "ti": "ዝለዓለ መሰል"
        },
        "premium_benefit_all_gold": {
            "en": "All Gold Buyer privileges included",
            "am": "ሁሉም የወርቅ ገዢ መብቶች ተካተዋል",
            "or": "Mirgi Bittaa Warqii hundi dabalameera",
            "ti": "ኩሎም ናይ ወርቂ ዓዳጊ መሰላት ተኻቲቶም"
        },
        "premium_benefit_legal_check": {
            "en": "Free Title Deed & Vehicle Inspection Legal Checks",
            "am": "ነፃ የባለቤትነት ማረጋገጫ (ካርታ) እና የተሽከርካሪ ሕጋዊ ምርመራ",
            "or": "Waraqaa Ragaa fi Qorannoo Seeraa Konkolaataa Bilaashaa",
            "ti": "ናጻ ናይ ዋንነት ምስክር ወረቐትን ሕጋዊ ፍተሻ ተሽከርከርትን"
        }
    },
    "listing_fee": {
        "startup_rate": {
            "en": "Startup Rate",
            "am": "የመነሻ ተመን",
            "or": "Reetii Jalqabaa",
            "ti": "ናይ መበገሲ ተመን"
        },
        "duration_hours": {
            "en": "Duration Hours",
            "am": "የቆይታ ሰዓታት",
            "or": "Sa'aatii Turtii",
            "ti": "ናይ ጻንሒት ሰዓታት"
        },
        "rate_calculation": {
            "en": "Rate Calculation",
            "am": "የተመን ስሌት",
            "or": "Herrega Reetii",
            "ti": "ስሌት ዋጋ"
        },
        "calculated_fee_rate": {
            "en": "Calculated Fee Rate",
            "am": "የተሰላ የክፍያ ተመን",
            "or": "Reetii Kaffaltii Herregame",
            "ti": "ዝተሓሰበ መጠን ተመን"
        },
        "total_fee_amount": {
            "en": "Total Fee Amount",
            "am": "ጠቅላላ የክፍያ መጠን",
            "or": "Hanga Kaffaltii Waliigalaa",
            "ti": "ጠቕላላ መጠን ክፍሊት"
        },
        "you_receive": {
            "en": "You Receive",
            "am": "የሚደርስዎት",
            "or": "Kan Isin Qaqqabu",
            "ti": "ዝበጽሓኩም"
        },
        "pay_listing_fee": {
            "en": "Pay Listing Fee",
            "am": "የማስታወቂያ ክፍያ ይክፈሉ",
            "or": "Kaffaltii Beeksisaa Kaffalaa",
            "ti": "ናይ ምልክታ ክፍሊት ክፈሉ"
        },
        "amount_due": {
            "en": "Amount Due",
            "am": "የሚፈለግ ክፍያ",
            "or": "Hanga Kaffalamuu Qabu",
            "ti": "ዝኽፈል መጠን"
        },
        "telebirr_direct": {
            "en": "Telebirr SuperApp Direct",
            "am": "ቴሌብር ቀጥታ ክፍያ",
            "or": "Telebirr Kallattiin",
            "ti": "ቴሌብር ቀጥታዊ ክፍሊት"
        },
        "official_payment_gateway": {
            "en": "Official Mobile Payment Gateway",
            "am": "ይፋዊ የሞባይል ክፍያ በር",
            "or": "Kaffaltii Moobaayilaa Seera Qabeessa",
            "ti": "ወግዓዊ ናይ ሞባይል ክፍሊት መገዲ"
        }
    },
    "not_found": {
        "return_home": {
            "en": "Return to Home",
            "am": "ወደ መነሻ ይመለሱ",
            "or": "Gara Manatti Deebi'aa",
            "ti": "ናብ መበገሲ ተመለሱ"
        },
        "browse_marketplace": {
            "en": "Browse Marketplace",
            "am": "ገበያውን ያስሱ",
            "or": "Gabaa Daawwadhaa",
            "ti": "ዕዳጋ ርኣዩ"
        },
        "page_not_found": {
            "en": "Page Not Found",
            "am": "ገጹ አልተገኘም",
            "or": "Fuulli Hin Argamne",
            "ti": "ገጽ ኣይተረኽበን"
        },
        "product_not_found": {
            "en": "Product Not Found",
            "am": "ምርቱ አልተገኘም",
            "or": "Meeshaan Hin Argamne",
            "ti": "ፍርያት ኣይተረኽበን"
        }
    }
}

# Load current files
data = {}
for loc in LOCALES:
    with open(FILE_PATHS[loc], 'r', encoding='utf-8') as f:
        data[loc] = json.load(f)

# Apply additions
for section, items in ADDITIONS.items():
    for key, loc_vals in items.items():
        for loc in LOCALES:
            if section not in data[loc]:
                data[loc][section] = {}
            data[loc][section][key] = loc_vals[loc]

# Fix any legacy "ሓበሻ ማርኬት" brand reference in ti.json
if "footer" in data["ti"] and "copyright" in data["ti"]["footer"]:
    data["ti"]["footer"]["copyright"] = "© 2026 HUNTIQ ። መሰሉ ብሕጊ ዝተሓለወ እዩ።"

# Verify 100% key parity
def get_all_paths(d, prefix=""):
    paths = set()
    for k, v in d.items():
        cur = f"{prefix}.{k}" if prefix else k
        if isinstance(v, dict):
            paths.update(get_all_paths(v, cur))
        else:
            paths.add(cur)
    return paths

paths_by_loc = {loc: get_all_paths(data[loc]) for loc in LOCALES}
base_paths = paths_by_loc['en']

parity_ok = True
for loc in ['am', 'or', 'ti']:
    missing = base_paths - paths_by_loc[loc]
    extra = paths_by_loc[loc] - base_paths
    if missing:
        print(f"Error: {loc} is missing {len(missing)} keys!")
        parity_ok = False
    if extra:
        print(f"Error: {loc} has {len(extra)} extra keys!")
        parity_ok = False

if parity_ok:
    print(f"Parity verified! Exactly {len(base_paths)} keys in each locale.")
    for loc in LOCALES:
        with open(FILE_PATHS[loc], 'w', encoding='utf-8') as f:
            json.dump(data[loc], f, ensure_ascii=False, indent=2)
            f.write('\n')
    print("All 4 locale files updated successfully.")
else:
    print("Aborting save due to parity errors.")
