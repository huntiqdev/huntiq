'use client';

import { useState, useEffect, useMemo } from 'react';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import {
  MapPin,
  Navigation,
  ChevronDown,
  X,
  Compass,
  Check,
  AlertCircle,
  WifiOff,
  Search,
} from 'lucide-react';
import { LocationSelector, LocationValue } from './LocationSelector';

export interface LocationNode {
  id: number;
  name: string;
  slug: string;
  location_type?: string;
  level?: number;
  latitude?: number | null;
  longitude?: number | null;
  children?: LocationNode[];
}

export interface MasterLocationBarProps {
  activeRegion?: string;
  activeZone?: string;
  activeSubcity?: string;
  activeWoreda?: string;
  activeCoords?: { lat: number; lng: number } | null;
  activeMaxDistance?: number;
  onLocationChange: (loc: {
    region?: string;
    zone?: string;
    subcity?: string;
    woreda?: string;
    coords?: { lat: number; lng: number } | null;
    maxDistance?: number;
    sortBy?: 'nearest' | 'trending';
  }) => void;
  className?: string;
}

export function MasterLocationBar({
  activeRegion = '',
  activeZone = '',
  activeSubcity = '',
  activeWoreda = '',
  activeCoords = null,
  activeMaxDistance,
  onLocationChange,
  className = '',
}: MasterLocationBarProps) {
  const { t, locale } = useLocale();
  const [dbLocations, setDbLocations] = useState<LocationNode[]>([]);
  const [selectedParentNode, setSelectedParentNode] = useState<LocationNode | null>(null);
  const [childLocations, setChildLocations] = useState<LocationNode[]>([]);
  const [isLoadingLocations, setIsLoadingLocations] = useState(true);

  const [isLocating, setIsLocating] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [radiusPickerOpen, setRadiusPickerOpen] = useState(false);
  const [insecureAlert, setInsecureAlert] = useState<string | null>(null);
  const [isOnline, setIsOnline] = useState(true);

  // Monitor mobile device network online status
  useEffect(() => {
    if (typeof window === 'undefined') return;
    setIsOnline(navigator.onLine);
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  // Fetch locations purely from the database backend
  useEffect(() => {
    let active = true;
    setIsLoadingLocations(true);

    api.getLocationTree(locale || 'en')
      .then((tree) => {
        if (!active) return;
        if (Array.isArray(tree) && tree.length > 0) {
          // If root is 'Ethiopia' (level 0), unwrap to primary administrative regions/cities
          if (tree.length === 1 && tree[0].children && tree[0].children.length > 0) {
            setDbLocations(tree[0].children);
          } else {
            setDbLocations(tree);
          }
        } else {
          // Fallback to legacy database regions endpoint if tree empty
          api.getRegions().then((regs) => {
            if (!active) return;
            const mapped: LocationNode[] = (regs || []).map((r, idx) => ({
              id: idx + 1,
              name: (locale === 'am' && r.am) ? r.am : (locale === 'or' && r.or) ? r.or : r.en,
              slug: r.en.toLowerCase().replace(/\s+/g, '-'),
              location_type: 'region',
            }));
            setDbLocations(mapped);
          }).catch(() => {});
        }
      })
      .catch(() => {
        // Fallback to database regions endpoint
        api.getRegions().then((regs) => {
          if (!active) return;
          const mapped: LocationNode[] = (regs || []).map((r, idx) => ({
            id: idx + 1,
            name: (locale === 'am' && r.am) ? r.am : (locale === 'or' && r.or) ? r.or : r.en,
            slug: r.en.toLowerCase().replace(/\s+/g, '-'),
            location_type: 'region',
          }));
          setDbLocations(mapped);
        }).catch(() => {});
      })
      .finally(() => {
        if (active) setIsLoadingLocations(false);
      });

    return () => {
      active = false;
    };
  }, [locale]);

  // When activeRegion changes, synchronize selected parent node and its child locations from DB
  useEffect(() => {
    if (!activeRegion) {
      setSelectedParentNode(null);
      setChildLocations([]);
      return;
    }

    const matched = dbLocations.find(
      (n) => n.name.toLowerCase() === activeRegion.toLowerCase() || n.slug === activeRegion.toLowerCase()
    );

    if (matched) {
      setSelectedParentNode(matched);
      if (matched.children && matched.children.length > 0) {
        setChildLocations(matched.children);
      } else {
        // Fetch children dynamically from backend
        api.getLocationChildren(matched.id, locale || 'en').then((children) => {
          if (Array.isArray(children)) setChildLocations(children);
        }).catch(() => {});
      }
    }
  }, [activeRegion, dbLocations, locale]);

  // Compute active location display label
  const locationLabel = useMemo(() => {
    const parts: string[] = [];
    if (activeSubcity) parts.push(activeSubcity);
    else if (activeZone) parts.push(activeZone);
    if (activeRegion) parts.push(activeRegion);
    if (activeWoreda) parts.push(activeWoreda);

    if (parts.length > 0) {
      return parts.join(', ');
    }
    if (activeCoords) {
      return `GPS: ${activeCoords.lat.toFixed(4)}, ${activeCoords.lng.toFixed(4)}`;
    }
    return locale === 'am' ? 'መላው ኢትዮጵያ' : 'All Ethiopia';
  }, [activeRegion, activeZone, activeSubcity, activeWoreda, activeCoords, locale]);

  const hasLocation = Boolean(activeRegion || activeZone || activeSubcity || activeWoreda || activeCoords);

  // 1-Click GPS Location Detection (Accepts any actual XY coordinates freely)
  const handleNearMeClick = () => {
    setInsecureAlert(null);

    if (typeof window === 'undefined') return;

    // Mobile HTTP Check: Chrome/Safari on mobile blocks geolocation over unencrypted HTTP
    if (window.isSecureContext === false && window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1') {
      setInsecureAlert(
        locale === 'am'
          ? 'በሞባይል አውታረ መረብ የጂፒኤስ ፍተሻ HTTPS ይፈልጋል። እባክዎ ከታች ካሉት አካባቢዎች ይምረጡ።'
          : 'Mobile browser GPS requires HTTPS or localhost. Please tap your location directly from the database list below.'
      );
      return;
    }

    if (!navigator.geolocation) {
      setInsecureAlert(
        locale === 'am'
          ? 'የእርስዎ ስልክ/ብሮውዘር የጂፒኤስ ፍተሻ አይደግፍም። አካባቢዎን ከታች ይምረጡ።'
          : 'Your browser does not support geolocation. Please select your location from the list below.'
      );
      return;
    }

    setIsLocating(true);

    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;

        try {
          // Attempt reverse geocoding from database to resolve administrative division if known
          const res = await api.reverseGeocode(lat, lng, locale || 'en');
          const data = res?.data || res;
          const levels = data?.resolved_levels || {};

          onLocationChange({
            region: levels.region || undefined,
            zone: levels.zone || undefined,
            subcity: levels.subcity || undefined,
            woreda: levels.woreda || undefined,
            coords: { lat, lng },
            maxDistance: activeMaxDistance || 25,
            sortBy: 'nearest',
          });
        } catch {
          // If not registered in hierarchy or reverse geocoding is unavailable,
          // ALWAYS ACCEPT coordinates directly without rejecting!
          onLocationChange({
            coords: { lat, lng },
            maxDistance: activeMaxDistance || 25,
            sortBy: 'nearest',
          });
        } finally {
          setIsLocating(false);
        }
      },
      (err) => {
        setIsLocating(false);
        let msg =
          locale === 'am'
            ? 'አካባቢዎን ማግኘት አልተቻለም። እባክዎ አካባቢዎን ከታች ይምረጡ።'
            : 'Could not detect your exact position. Please select your location below.';

        if (err.code === err.PERMISSION_DENIED) {
          msg =
            locale === 'am'
              ? 'የአካባቢ ፈቃድ አልተሰጠም። እባክዎ አካባቢዎን ከታች ይምረጡ።'
              : 'Location permission was denied. Tap your location from the list below.';
        } else if (err.code === err.TIMEOUT) {
          msg =
            locale === 'am'
              ? 'የጂፒኤስ ጥሪ ጊዜ አልፏል። እባክዎ አካባቢዎን ከታች ይምረጡ።'
              : 'GPS request timed out. Please tap your location below.';
        }
        setInsecureAlert(msg);
      },
      { enableHighAccuracy: true, timeout: 8000, maximumAge: 60000 }
    );
  };

  // Node Selection: Uses database coordinates when available to sync XY GPS
  const handleSelectNode = (node: LocationNode) => {
    setInsecureAlert(null);
    const coords = (node.latitude && node.longitude) ? { lat: node.latitude, lng: node.longitude } : null;

    onLocationChange({
      region: node.name,
      subcity: undefined,
      zone: undefined,
      woreda: undefined,
      coords,
      maxDistance: node.location_type === 'city' ? 25 : 75,
      sortBy: coords ? 'nearest' : 'trending',
    });
  };

  const handleSelectChildNode = (child: LocationNode) => {
    setInsecureAlert(null);
    const coords = (child.latitude && child.longitude) ? { lat: child.latitude, lng: child.longitude } : activeCoords;

    onLocationChange({
      region: activeRegion,
      subcity: child.location_type === 'subcity' ? child.name : undefined,
      zone: child.location_type === 'zone' ? child.name : undefined,
      coords,
      maxDistance: 15,
      sortBy: coords ? 'nearest' : 'trending',
    });
  };

  const handleClearLocation = () => {
    setInsecureAlert(null);
    setSelectedParentNode(null);
    setChildLocations([]);
    onLocationChange({
      region: undefined,
      zone: undefined,
      subcity: undefined,
      woreda: undefined,
      coords: null,
      maxDistance: undefined,
      sortBy: 'trending',
    });
  };

  const handleDistanceSelect = (km: number | undefined) => {
    setRadiusPickerOpen(false);
    onLocationChange({
      region: activeRegion || undefined,
      zone: activeZone || undefined,
      subcity: activeSubcity || undefined,
      woreda: activeWoreda || undefined,
      coords: activeCoords,
      maxDistance: km,
      sortBy: activeCoords ? 'nearest' : 'trending',
    });
  };

  return (
    <div className={`w-full ${className}`}>
      {/* ── Offline Network Banner (for mobile phones with unstable connections) ── */}
      {!isOnline && (
        <div className="mb-2 px-3 py-1.5 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-700 dark:text-amber-400 text-xs font-bold flex items-center gap-2">
          <WifiOff size={14} className="shrink-0 animate-pulse text-amber-600" />
          <span>
            {locale === 'am'
              ? 'የኔትወርክ ግንኙነት የለም። የተቀመጡ መረጃዎች እየታዩ ነው።'
              : 'Network offline. Showing cached listings.'}
          </span>
        </div>
      )}

      {/* ── Insecure Context / GPS Permission Alert ── */}
      {insecureAlert && (
        <div className="mb-2.5 px-3 py-2 rounded-xl bg-blue-500/10 border border-blue-500/30 text-blue-800 dark:text-blue-300 text-xs font-semibold flex items-center justify-between gap-2 animate-fade-in">
          <div className="flex items-center gap-2">
            <AlertCircle size={15} className="text-[#1769E0] shrink-0" />
            <span>{insecureAlert}</span>
          </div>
          <button
            onClick={() => setInsecureAlert(null)}
            className="text-blue-600 hover:text-blue-900 dark:hover:text-white p-1"
          >
            <X size={14} />
          </button>
        </div>
      )}

      {/* ── MASTER LOCATION BAR CONTAINER ── */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-2.5 sm:p-3 shadow-xs">
        {/* Row 1: Active Location Header + GPS Action + Radius */}
        <div className="flex items-center justify-between gap-2 flex-wrap sm:flex-nowrap pb-2 border-b border-zinc-150 dark:border-zinc-800/80">
          {/* Left: Active Location Indicator */}
          <div className="flex items-center gap-2 min-w-0">
            <div className="w-8 h-8 rounded-xl bg-[#1769E0]/10 dark:bg-[#1769E0]/20 text-[#1769E0] dark:text-[#38A5FF] flex items-center justify-center shrink-0">
              <MapPin size={17} className="animate-bounce-subtle" />
            </div>
            <div className="min-w-0">
              <div className="flex items-center gap-1.5">
                <span className="text-[10px] font-black uppercase tracking-wider text-zinc-400 dark:text-zinc-500">
                  {locale === 'am' ? 'ዋና የአካባቢ ማጣሪያ' : 'Master Location Filter'}
                </span>
                {activeCoords && (
                  <span className="px-1.5 py-0.2 rounded text-[9px] font-black bg-teal-500/10 text-teal-600 dark:text-teal-400 border border-teal-500/20">
                    GPS SYNC
                  </span>
                )}
              </div>
              <p className="text-xs sm:text-sm font-black text-zinc-900 dark:text-white truncate">
                {locationLabel}
              </p>
            </div>

            {hasLocation && (
              <button
                type="button"
                onClick={handleClearLocation}
                className="ml-1 p-1 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg text-zinc-400 hover:text-red-500 transition-colors"
                title="Clear location filter"
              >
                <X size={14} />
              </button>
            )}
          </div>

          {/* Right: GPS Near Me Button & Distance Selector */}
          <div className="flex items-center gap-2 shrink-0">
            {/* GPS Near Me Button */}
            <button
              type="button"
              onClick={handleNearMeClick}
              disabled={isLocating}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-black transition-all cursor-pointer shadow-xs ${
                activeCoords
                  ? 'bg-teal-500 text-white shadow-teal-500/25 hover:bg-teal-600'
                  : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-[#1769E0] hover:text-white text-zinc-800 dark:text-zinc-200 border border-zinc-200/80 dark:border-zinc-700/80'
              }`}
              title="Detect listings near my current GPS location"
            >
              <Navigation
                size={13}
                className={isLocating ? 'animate-spin' : activeCoords ? 'fill-white' : ''}
              />
              <span>
                {isLocating
                  ? locale === 'am'
                    ? 'በመፈለግ ላይ…'
                    : 'Locating…'
                  : locale === 'am'
                  ? 'አካባቢዬ (GPS)'
                  : 'Near Me'}
              </span>
            </button>

            {/* Radius Picker (Visible when coordinates or specific area active) */}
            {(activeCoords || activeRegion) && (
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setRadiusPickerOpen((v) => !v)}
                  className="flex items-center gap-1 px-2.5 py-1.5 rounded-xl text-xs font-bold bg-zinc-100 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 hover:border-zinc-400 border border-zinc-200 dark:border-zinc-700"
                >
                  <Compass size={12} className="text-zinc-500" />
                  <span>
                    {activeMaxDistance ? `${activeMaxDistance} km` : locale === 'am' ? 'ሙሉ ራዲየስ' : 'Any km'}
                  </span>
                  <ChevronDown size={11} className="text-zinc-400" />
                </button>

                {radiusPickerOpen && (
                  <div className="absolute right-0 top-full mt-1 w-36 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-xl z-50 p-1 text-xs">
                    <p className="px-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400">
                      {locale === 'am' ? 'የፍተሻ ርቀት' : 'Max Radius'}
                    </p>
                    {[5, 10, 25, 50, 100].map((km) => (
                      <button
                        key={km}
                        onClick={() => handleDistanceSelect(km)}
                        className={`w-full text-left px-2 py-1.5 rounded-lg font-bold flex items-center justify-between ${
                          activeMaxDistance === km
                            ? 'bg-[#1769E0] text-white'
                            : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-700 dark:text-zinc-300'
                        }`}
                      >
                        <span>{km} km</span>
                        {activeMaxDistance === km && <Check size={12} />}
                      </button>
                    ))}
                    <button
                      onClick={() => handleDistanceSelect(undefined)}
                      className="w-full text-left px-2 py-1.5 rounded-lg font-bold hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-500"
                    >
                      {locale === 'am' ? 'ገደብ የለሽ' : 'Entire Area'}
                    </button>
                  </div>
                )}
              </div>
            )}

            {/* Advanced Hierarchy Modal Trigger */}
            <button
              type="button"
              onClick={() => setModalOpen(true)}
              className="flex items-center gap-1 px-2.5 py-1.5 rounded-xl text-xs font-bold text-[#1769E0] dark:text-[#38A5FF] hover:bg-[#1769E0]/10 transition-colors"
              title="Search by specific Woreda or custom location"
            >
              <Search size={12} />
              <span className="hidden sm:inline">
                {locale === 'am' ? 'ዝርዝር ወረዳ' : 'Woredas'}
              </span>
            </button>
          </div>
        </div>

        {/* Row 2: Database-Driven Location Pills (Purely fetched from PostgreSQL loccatdata) */}
        <div className="pt-2">
          <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-0.5">
            {/* All Ethiopia Pill */}
            <button
              type="button"
              onClick={handleClearLocation}
              className={`flex items-center gap-1 px-3 py-1 rounded-xl text-xs font-black whitespace-nowrap shrink-0 transition-all ${
                !hasLocation
                  ? 'bg-[#1769E0] text-white shadow-xs'
                  : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300'
              }`}
            >
              <span>🇪🇹</span>
              <span>{locale === 'am' ? 'መላው ኢትዮጵያ' : 'All Ethiopia'}</span>
            </button>

            {/* Loading Skeleton if locations still fetching */}
            {isLoadingLocations && dbLocations.length === 0 && (
              <div className="flex items-center gap-2">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="h-7 w-20 bg-zinc-200 dark:bg-zinc-800 rounded-xl animate-pulse shrink-0"
                  />
                ))}
              </div>
            )}

            {/* Render Every Region/City Purely from Database */}
            {dbLocations.map((locNode) => {
              const isActive =
                activeRegion.toLowerCase() === locNode.name.toLowerCase() ||
                activeRegion.toLowerCase() === locNode.slug.toLowerCase();

              return (
                <button
                  key={locNode.id || locNode.slug}
                  type="button"
                  onClick={() => handleSelectNode(locNode)}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded-xl text-xs font-bold whitespace-nowrap shrink-0 transition-all ${
                    isActive
                      ? 'bg-[#1769E0] text-white shadow-xs'
                      : 'bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300'
                  }`}
                >
                  <span>{locNode.location_type === 'city' ? '🏙️' : '🌾'}</span>
                  <span>{locNode.name}</span>
                </button>
              );
            })}
          </div>

          {/* Row 3: Child Subcities / Zones (Dynamically populated from backend database) */}
          {childLocations.length > 0 && (
            <div className="flex items-center gap-1.5 overflow-x-auto no-scrollbar pt-2 mt-1.5 border-t border-zinc-100 dark:border-zinc-800/60 animate-fade-in">
              <span className="text-[10px] font-black uppercase tracking-wider text-indigo-500 dark:text-indigo-400 shrink-0">
                {selectedParentNode?.name}:
              </span>
              {childLocations.map((child) => {
                const isActive =
                  activeSubcity.toLowerCase() === child.name.toLowerCase() ||
                  activeZone.toLowerCase() === child.name.toLowerCase();

                return (
                  <button
                    key={child.id || child.slug}
                    type="button"
                    onClick={() => handleSelectChildNode(child)}
                    className={`px-2.5 py-0.8 rounded-lg text-xs font-semibold whitespace-nowrap shrink-0 transition-all ${
                      isActive
                        ? 'bg-indigo-600 text-white font-black shadow-xs'
                        : 'bg-indigo-50 dark:bg-indigo-950/30 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-100 dark:hover:bg-indigo-900/50'
                    }`}
                  >
                    {child.name}
                  </button>
                );
              })}
            </div>
          )}
        </div>
      </div>

      {/* ── ADVANCED HIERARCHY MODAL (Woredas & Kebeles — Pure DB Driven) ── */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-fade-in">
          <div className="relative w-full max-w-lg bg-white dark:bg-zinc-900 rounded-3xl p-5 shadow-2xl border border-zinc-200 dark:border-zinc-800">
            <div className="flex items-center justify-between pb-3 border-b border-zinc-200 dark:border-zinc-800">
              <div className="flex items-center gap-2">
                <MapPin size={18} className="text-[#1769E0]" />
                <h3 className="text-base font-black text-zinc-900 dark:text-white">
                  {locale === 'am' ? 'አካባቢ ይምረጡ (ክልል፣ ዞን፣ ወረዳ)' : 'Select Location'}
                </h3>
              </div>
              <button
                type="button"
                onClick={() => setModalOpen(false)}
                className="p-1 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl text-zinc-400 hover:text-zinc-900 dark:hover:text-white"
              >
                <X size={18} />
              </button>
            </div>

            <div className="py-4">
              <LocationSelector
                showPrecise={true}
                defaultValue={{
                  region: activeRegion,
                  zone: activeZone,
                  subcity: activeSubcity,
                  woreda: activeWoreda,
                  latitude: activeCoords?.lat,
                  longitude: activeCoords?.lng,
                }}
                onChange={(val: LocationValue) => {
                  onLocationChange({
                    region: val.region || undefined,
                    zone: val.zone || undefined,
                    subcity: val.subcity || undefined,
                    woreda: val.woreda || undefined,
                    coords:
                      val.latitude && val.longitude
                        ? { lat: val.latitude, lng: val.longitude }
                        : activeCoords,
                    maxDistance: activeMaxDistance || 25,
                    sortBy: val.latitude && val.longitude ? 'nearest' : 'trending',
                  });
                  setModalOpen(false);
                }}
              />
            </div>

            <div className="pt-3 border-t border-zinc-150 dark:border-zinc-800 flex justify-end">
              <button
                type="button"
                onClick={() => setModalOpen(false)}
                className="px-4 py-2 rounded-xl bg-[#1769E0] text-white text-xs font-bold hover:bg-[#1055BD]"
              >
                {locale === 'am' ? 'ዝጋ' : 'Close'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
