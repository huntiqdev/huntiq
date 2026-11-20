'use client';

import { useState, useEffect, useRef } from 'react';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { Navigation, MapPin, X, ChevronDown, Search, Compass, AlertCircle, Check } from 'lucide-react';

interface Props {
  activeCoords: { lat: number; lng: number } | null;
  activeLocationLabel?: string;
  activeMaxDistance?: number;
  onSearch: (coords: { lat: number; lng: number } | null, label?: string, maxDistance?: number) => void;
  className?: string;
}

export function NearMeFindButton({
  activeCoords,
  activeLocationLabel = '',
  activeMaxDistance = 25,
  onSearch,
  className = '',
}: Props) {
  const { t, locale } = useLocale();
  const [isLocating, setIsLocating] = useState(false);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [manualQuery, setManualQuery] = useState('');
  const [suggestions, setSuggestions] = useState<any[]>([]);
  const [selectedRadius, setSelectedRadius] = useState<number>(activeMaxDistance || 25);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Debounced search for manual location input from PostgreSQL database
  useEffect(() => {
    if (!manualQuery.trim()) {
      setSuggestions([]);
      return;
    }

    // Check if input is direct raw XY coordinates: "9.03, 38.74" or "9.03 38.74"
    const coordMatch = manualQuery.trim().match(/^(-?\d+(\.\d+)?)[,\s]+(-?\d+(\.\d+)?)$/);
    if (coordMatch) {
      setSuggestions([]);
      return;
    }

    const timer = setTimeout(() => {
      api.searchLocations(manualQuery.trim())
        .then((res) => {
          setSuggestions(Array.isArray(res) ? res.slice(0, 5) : []);
        })
        .catch(() => setSuggestions([]));
    }, 250);

    return () => clearTimeout(timer);
  }, [manualQuery]);

  // Automatic GPS Location Detection
  const handleAutomaticLocate = () => {
    setStatusMessage(null);

    if (typeof window === 'undefined') return;

    // Check mobile insecure context
    const isInsecureContext =
      window.isSecureContext === false &&
      window.location.hostname !== 'localhost' &&
      window.location.hostname !== '127.0.0.1';

    if (isInsecureContext || !navigator.geolocation) {
      setDropdownOpen(true);
      setStatusMessage(
        locale === 'am'
          ? 'በሞባይል አውታረ መረብ የጂፒኤስ ፍተሻ HTTPS ይፈልጋል። እባክዎ አካባቢዎን እዚህ ይጻፉ።'
          : 'Mobile browser GPS requires HTTPS or localhost. Please type your location manually below.'
      );
      return;
    }

    setIsLocating(true);

    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;

        try {
          const res = await api.reverseGeocode(lat, lng, locale || 'en');
          const data = res?.data || res;
          const label = data?.name || `${lat.toFixed(2)}, ${lng.toFixed(2)}`;
          onSearch({ lat, lng }, label, selectedRadius);
        } catch {
          // If reverse geocode unavailable, still apply actual GPS coordinates directly
          onSearch({ lat, lng }, `${lat.toFixed(2)}, ${lng.toFixed(2)}`, selectedRadius);
        } finally {
          setIsLocating(false);
          setDropdownOpen(false);
        }
      },
      (err) => {
        setIsLocating(false);
        setDropdownOpen(true);
        let msg =
          locale === 'am'
            ? 'አካባቢዎን ማግኘት አልተቻለም። እባክዎ እዚህ ይጻፉ።'
            : 'Could not acquire GPS position. Please enter your location manually.';
        if (err.code === err.PERMISSION_DENIED) {
          msg =
            locale === 'am'
              ? 'የአካባቢ ፈቃድ አልተሰጠም። እባክዎ አካባቢዎን እዚህ ይጻፉ።'
              : 'Location permission was denied. Please enter your location manually.';
        }
        setStatusMessage(msg);
      },
      { enableHighAccuracy: true, timeout: 8000, maximumAge: 60000 }
    );
  };

  // Manual submission: handles both raw coordinates and text locations
  const handleManualSubmit = (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    const query = manualQuery.trim();
    if (!query) return;

    // Check if input is direct coordinates: e.g. "9.03, 38.74"
    const coordMatch = query.match(/^(-?\d+(\.\d+)?)[,\s]+(-?\d+(\.\d+)?)$/);
    if (coordMatch) {
      const lat = parseFloat(coordMatch[1]);
      const lng = parseFloat(coordMatch[3]);
      if (!isNaN(lat) && !isNaN(lng)) {
        onSearch({ lat, lng }, `${lat.toFixed(2)}, ${lng.toFixed(2)}`, selectedRadius);
        setDropdownOpen(false);
        setManualQuery('');
        return;
      }
    }

    // If query matches a suggestion or place name
    if (suggestions.length > 0) {
      handleSelectSuggestion(suggestions[0]);
    } else {
      // Pass as query label
      onSearch(null, query, selectedRadius);
      setDropdownOpen(false);
    }
  };

  const handleSelectSuggestion = (s: any) => {
    const lat = s.latitude ? Number(s.latitude) : null;
    const lng = s.longitude ? Number(s.longitude) : null;
    const label = s.subcity || s.woreda || s.zone || s.region || s.name || manualQuery;

    if (lat && lng) {
      onSearch({ lat, lng }, label, selectedRadius);
    } else {
      onSearch(null, label, selectedRadius);
    }
    setDropdownOpen(false);
    setManualQuery('');
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onSearch(null, '', undefined);
    setDropdownOpen(false);
  };

  const isActive = Boolean(activeCoords || activeLocationLabel);

  return (
    <div className={`relative inline-block ${className}`} ref={dropdownRef}>
      {/* ── SINGLE "NEAR ME FIND" BUTTON ── */}
      <div className="inline-flex items-center rounded-xl shadow-xs transition-all">
        <button
          type="button"
          onClick={() => {
            if (isActive) {
              setDropdownOpen((v) => !v);
            } else {
              handleAutomaticLocate();
            }
          }}
          disabled={isLocating}
          className={`flex items-center gap-2 px-3.5 py-2 rounded-xl text-xs font-black transition-all cursor-pointer border ${
            isActive
              ? 'bg-teal-600 hover:bg-teal-700 text-white border-teal-600 shadow-md shadow-teal-600/20'
              : 'bg-white dark:bg-zinc-800 hover:bg-teal-50 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 border-zinc-200 dark:border-zinc-700 hover:border-teal-500'
          }`}
          title="Find products nearest to your location"
        >
          <Navigation
            size={14}
            className={`${
              isLocating ? 'animate-spin text-teal-400' : isActive ? 'fill-white text-white' : 'text-teal-600'
            }`}
          />
          <span>
            {isLocating
              ? locale === 'am'
                ? 'በመፈለግ ላይ…'
                : 'Locating…'
              : isActive
              ? `${locale === 'am' ? 'አካባቢዬ:' : 'Near Me:'} ${activeLocationLabel || 'Active'} (${selectedRadius}km)`
              : locale === 'am'
              ? 'አካባቢዬ ፈልግ (Near Me Find)'
              : 'Near Me Find'}
          </span>

          {/* If not active, show small chevron to open manual entry */}
          {!isActive && (
            <span
              onClick={(e) => {
                e.stopPropagation();
                setDropdownOpen((v) => !v);
              }}
              className="pl-1 text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 cursor-pointer"
            >
              <ChevronDown size={12} />
            </span>
          )}

          {/* If active, 1-click remove button */}
          {isActive && (
            <span
              onClick={handleClear}
              className="ml-1 hover:bg-white/20 rounded-full p-0.5 text-white/90 hover:text-white cursor-pointer"
              title="Clear near me filter"
            >
              <X size={13} />
            </span>
          )}
        </button>
      </div>

      {/* ── SLEEK MANUAL & CUSTOM LOCATION DROPDOWN (ONLY SHOWN WHEN REQUESTED) ── */}
      {dropdownOpen && (
        <div className="absolute left-0 sm:right-auto top-full mt-2 w-80 sm:w-96 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-2xl z-50 p-3 text-xs animate-fade-in">
          {/* Header */}
          <div className="flex items-center justify-between pb-2 border-b border-zinc-150 dark:border-zinc-800">
            <div className="flex items-center gap-1.5 font-black text-zinc-900 dark:text-white">
              <MapPin size={14} className="text-teal-600" />
              <span>{locale === 'am' ? 'በአካባቢ ፈልግ' : 'Find Listings Near Me'}</span>
            </div>
            <button
              type="button"
              onClick={() => setDropdownOpen(false)}
              className="p-1 text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200"
            >
              <X size={14} />
            </button>
          </div>

          {/* Status Alert if GPS was blocked or unavailable */}
          {statusMessage && (
            <div className="my-2 p-2 rounded-xl bg-blue-50 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-800 text-blue-700 dark:text-blue-300 flex items-start gap-1.5 text-[11px]">
              <AlertCircle size={13} className="shrink-0 mt-0.5" />
              <span>{statusMessage}</span>
            </div>
          )}

          {/* 1-Tap Auto GPS Button inside dropdown */}
          <div className="py-2">
            <button
              type="button"
              onClick={handleAutomaticLocate}
              disabled={isLocating}
              className="w-full py-2 px-3 rounded-xl bg-teal-50 dark:bg-teal-950/30 hover:bg-teal-100 dark:hover:bg-teal-900/40 text-teal-700 dark:text-teal-300 border border-teal-200 dark:border-teal-800 font-extrabold flex items-center justify-center gap-1.5 cursor-pointer transition-colors"
            >
              <Navigation size={13} className={isLocating ? 'animate-spin' : ''} />
              <span>
                {isLocating
                  ? locale === 'am'
                    ? 'በመፈለግ ላይ…'
                    : 'Detecting Location…'
                  : locale === 'am'
                  ? 'የአሁን ጂፒኤስ ቦታዬን ተጠቀም (Auto GPS)'
                  : 'Use Current Device Location (Auto GPS)'}
              </span>
            </button>
          </div>

          <div className="flex items-center gap-2 my-1 text-[10px] text-zinc-400 font-black uppercase tracking-wider">
            <div className="flex-1 h-px bg-zinc-200 dark:bg-zinc-800" />
            <span>{locale === 'am' ? 'ወይም በእጅ ጻፍ' : 'OR TYPE MANUALLY'}</span>
            <div className="flex-1 h-px bg-zinc-200 dark:bg-zinc-800" />
          </div>

          {/* Manual Input Form: Address, Place Name, or Coordinates */}
          <form onSubmit={handleManualSubmit} className="space-y-2 pt-1">
            <div className="relative">
              <input
                type="text"
                value={manualQuery}
                onChange={(e) => setManualQuery(e.target.value)}
                placeholder={
                  locale === 'am'
                    ? 'ከተማ፣ ክፍለ-ከተማ ወይም ጂፒኤስ (ምሳሌ፡ ቦሌ፣ ሀዋሳ፣ 9.03, 38.74)'
                    : 'City, area, or coordinates (e.g. Bole, Hawassa, or 9.03, 38.74)'
                }
                className="w-full pl-8 pr-3 py-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl outline-none focus:ring-2 focus:ring-teal-500 text-xs"
              />
              <Search size={13} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-zinc-400" />
            </div>

            {/* Suggestions dropdown from database */}
            {suggestions.length > 0 && (
              <div className="border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden divide-y divide-zinc-100 dark:divide-zinc-800">
                {suggestions.map((s, idx) => {
                  const title = s.subcity || s.woreda || s.zone || s.region || s.name;
                  const subtitle = [s.subcity, s.zone, s.region].filter(Boolean).join(' › ');
                  return (
                    <button
                      key={idx}
                      type="button"
                      onClick={() => handleSelectSuggestion(s)}
                      className="w-full text-left px-2.5 py-1.5 hover:bg-zinc-50 dark:hover:bg-zinc-800/80 flex items-center justify-between text-xs"
                    >
                      <div>
                        <p className="font-bold text-zinc-800 dark:text-zinc-200">{title}</p>
                        {subtitle && <p className="text-[10px] text-zinc-400">{subtitle}</p>}
                      </div>
                      <MapPin size={11} className="text-zinc-400" />
                    </button>
                  );
                })}
              </div>
            )}

            {/* Distance Radius Picker */}
            <div className="pt-1">
              <p className="text-[10px] font-black uppercase tracking-wider text-zinc-400 mb-1 flex items-center gap-1">
                <Compass size={11} />
                <span>{locale === 'am' ? 'የፍተሻ ርቀት (ኪ.ሜ)' : 'Search Radius'}</span>
              </p>
              <div className="grid grid-cols-5 gap-1">
                {[5, 10, 25, 50, 100].map((km) => (
                  <button
                    key={km}
                    type="button"
                    onClick={() => setSelectedRadius(km)}
                    className={`py-1 rounded-lg font-extrabold text-[11px] transition-all ${
                      selectedRadius === km
                        ? 'bg-teal-600 text-white shadow-xs'
                        : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-200 dark:hover:bg-zinc-700'
                    }`}
                  >
                    {km} km
                  </button>
                ))}
              </div>
            </div>

            {/* Submit Action */}
            <div className="pt-2 flex items-center justify-between gap-2">
              {isActive && (
                <button
                  type="button"
                  onClick={handleClear}
                  className="px-3 py-1.5 rounded-xl border border-zinc-200 dark:border-zinc-700 text-zinc-500 font-bold hover:bg-zinc-100 dark:hover:bg-zinc-800"
                >
                  {locale === 'am' ? 'አጥፋ' : 'Clear'}
                </button>
              )}
              <button
                type="submit"
                className="flex-1 py-2 px-3 rounded-xl bg-teal-600 hover:bg-teal-700 text-white font-black text-xs transition-colors flex items-center justify-center gap-1 shadow-md shadow-teal-600/20"
              >
                <span>{locale === 'am' ? 'ፈልግ (Find)' : 'Find Nearest Products'}</span>
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
