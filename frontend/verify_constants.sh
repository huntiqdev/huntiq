#!/bin/bash

echo "🔍 Verifying Frontend Constants Match Backend..."
echo "=================================================="
echo ""

ISSUES=0

echo "1. Checking tier thresholds..."
if grep -q "gold:.*100_000" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ❌ Old incorrect threshold found (100_000)"
    ISSUES=$((ISSUES + 1))
elif grep -q "gold:.*5900" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ Correct Gold threshold (5900)"
else
    echo "   ⚠️  Could not verify Gold threshold"
fi

if grep -q "premium:.*5_000_000" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ❌ Old incorrect threshold found (5_000_000)"
    ISSUES=$((ISSUES + 1))
elif grep -q "premium:.*15900" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ Correct Premium threshold (15900)"
else
    echo "   ⚠️  Could not verify Premium threshold"
fi
echo ""

echo "2. Checking fee percentages..."
if grep -q "startupPct: 2\.0.*weeklyIncrementPct: 0\.25" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ❌ Old incorrect Normal fees found (2.0/0.25)"
    ISSUES=$((ISSUES + 1))
elif grep -q "startupPct: 1\.0.*weeklyIncrementPct: 0\.5" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ Correct Normal fees (1.0/0.5)"
else
    echo "   ⚠️  Could not verify Normal fees"
fi

if grep -q "startupPct: 1\.5.*weeklyIncrementPct: 0\.15" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ❌ Old incorrect Gold fees found (1.5/0.15)"
    ISSUES=$((ISSUES + 1))
elif grep -q "startupPct: 0\.5.*weeklyIncrementPct: 0\.25" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ Correct Gold fees (0.5/0.25)"
else
    echo "   ⚠️  Could not verify Gold fees"
fi

if grep -q "startupPct: 1\.2.*weeklyIncrementPct: 0\.10" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ❌ Old incorrect Premium fees found (1.2/0.10)"
    ISSUES=$((ISSUES + 1))
elif grep -q "startupPct: 0\.25.*weeklyIncrementPct: 0\.1" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ Correct Premium fees (0.25/0.1)"
else
    echo "   ⚠️  Could not verify Premium fees"
fi
echo ""

echo "3. Checking validation constants..."
if grep -q "MAX_PRODUCT_DESCRIPTION_LENGTH = 2000" frontend/src/lib/constants/validation.ts 2>/dev/null; then
    echo "   ❌ Old incorrect description length (2000)"
    ISSUES=$((ISSUES + 1))
elif grep -q "MAX_PRODUCT_DESCRIPTION_LENGTH = 5000" frontend/src/lib/constants/validation.ts 2>/dev/null; then
    echo "   ✅ Correct description length (5000)"
else
    echo "   ⚠️  Could not verify description length"
fi

if grep -q "MIN_IMAGES_PER_PRODUCT = 1" frontend/src/lib/constants/validation.ts 2>/dev/null; then
    echo "   ❌ Old incorrect min images (1)"
    ISSUES=$((ISSUES + 1))
elif grep -q "MIN_IMAGES_PER_PRODUCT = 2" frontend/src/lib/constants/validation.ts 2>/dev/null; then
    echo "   ✅ Correct min images (2)"
else
    echo "   ⚠️  Could not verify min images"
fi

if grep -q "MAX_IMAGES_PER_PRODUCT = 5" frontend/src/lib/constants/validation.ts 2>/dev/null; then
    echo "   ⚠️  Conflicting MAX_IMAGES_PER_PRODUCT found (should be tier-based)"
fi
echo ""

echo "4. Checking API methods..."
if grep -q "getPublicSettings" frontend/src/lib/api/public.ts 2>/dev/null; then
    echo "   ✅ getPublicSettings() API method exists"
else
    echo "   ❌ Missing getPublicSettings() API method"
    ISSUES=$((ISSUES + 1))
fi
echo ""

echo "5. Checking for hardcoded limits..."
HARDCODED=$(grep -r "500\.0\|1_000_000" frontend/src/components/referral/ 2>/dev/null | wc -l)
if [ "$HARDCODED" -gt 0 ]; then
    echo "   ⚠️  Found $HARDCODED hardcoded withdrawal limits in referral components"
    echo "      (These will be fixed when SettingsProvider is integrated)"
else
    echo "   ✅ No hardcoded withdrawal limits"
fi
echo ""

echo "6. Checking helper functions support dynamic settings..."
if grep -q "function getTierFeeParams(tier: SellerTier, settings?: any)" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ getTierFeeParams() supports settings parameter"
else
    echo "   ❌ getTierFeeParams() missing settings support"
    ISSUES=$((ISSUES + 1))
fi

if grep -q "function getMaxImages(tier: SellerTier, settings?: any)" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ getMaxImages() supports settings parameter"
else
    echo "   ❌ getMaxImages() missing settings support"
    ISSUES=$((ISSUES + 1))
fi

if grep -q "function calcListingFee" frontend/src/lib/constants/fees.ts 2>/dev/null; then
    echo "   ✅ calcListingFee() function exists"
else
    echo "   ❌ calcListingFee() function missing"
    ISSUES=$((ISSUES + 1))
fi
echo ""

echo "=================================================="
if [ $ISSUES -eq 0 ]; then
    echo "✅ ALL CRITICAL CHECKS PASSED!"
    echo ""
    echo "Next steps:"
    echo "1. Create SettingsProvider in hooks/useSettings.tsx"
    echo "2. Wrap app with SettingsProvider"
    echo "3. Update components to use useSettings() hook"
    exit 0
else
    echo "❌ Found $ISSUES issue(s) - review above"
    echo ""
    echo "Critical issues must be fixed before deployment."
    exit 1
fi
