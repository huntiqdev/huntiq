'use client';
import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { useLocale } from '@/hooks/useLocale';
import { I18nLocation, locName, locPath } from '@/types/location';
import { Search, X, MapPin, Lock } from 'lucide-react';

export interface LocationValue {
  region?: string;
  zone?: string;
  subcity?: string;
  woreda?: string;
  kebele?: string;
  specific_area?: string;
  latitude?: number;
  longitude?: number;
}

interface Props {
  onChange: (v: LocationValue) => void;
  searchable?: boolean;
  /** show kebele + specific_area fields */
  showPrecise?: boolean;
  /** show GPS coordinate fields (gold/premium only) */
  showCoords?: boolean;
  defaultValue?: LocationValue;
}

const SEL = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 outline-none focus:ring-2 focus:ring-orange-600 transition-all';
const INP = 'w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 outline-none focus:ring-2 focus:ring-orange-600 transition-all';

export function LocationSelector({
  onChange, searchable, showPrecise = true, showCoords = false,
  defaultValue = {},
}: Props) {
  const { t, locale } = useLocale();

  const [regions,   setRegions]   = useState<I18nLocation[]>([]);
  const [zones,     setZones]     = useState<I18nLocation[]>([]);
  const [subcities, setSubcities] = useState<I18nLocation[]>([]);
  const [woredas,   setWoredas]   = useState<I18nLocation[]>([]);

  const [region,       setRegion]       = useState(defaultValue.region ?? '');
  const [zone,         setZone]         = useState(defaultValue.zone ?? '');
  const [subcity,      setSubcity]      = useState(defaultValue.subcity ?? '');
  const [woreda,       setWoreda]       = useState(defaultValue.woreda ?? '');
  const [kebele,       setKebele]       = useState(defaultValue.kebele ?? '');
  const [specificArea, setSpecificArea] = useState(defaultValue.specific_area ?? '');
  const [lat,          setLat]          = useState(defaultValue.latitude?.toString() ?? '');
  const [lng,          setLng]          = useState(defaultValue.longitude?.toString() ?? '');

  const [search,        setSearch]        = useState('');
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchOpen,    setSearchOpen]    = useState(false);

  useEffect(() => {
    api.getRegions().then(r => setRegions((r ?? []) as I18nLocation[])).catch(() => {});
  }, []);

  useEffect(() => {
    setZone(''); setSubcity(''); setWoreda('');
    setZones([]); setSubcities([]); setWoredas([]);
    if (!region) return;
    api.getZones(region).then(z => setZones((z ?? []) as I18nLocation[])).catch(() => {});
    api.getSubcities(region).then(s => setSubcities((s ?? []) as I18nLocation[])).catch(() => {});
  }, [region]);

  useEffect(() => {
    setSubcity(''); setWoreda('');
    setWoredas([]);
    if (!region || !zone) {
      // When zone is cleared, reload all subcities for the region
      if (region) {
        api.getSubcities(region).then(s => setSubcities((s ?? []) as I18nLocation[])).catch(() => {});
      }
      return;
    }
    api.getSubcities(region, zone)
      .then(s => setSubcities((s ?? []) as I18nLocation[]))
      .catch(() => {});
  }, [zone]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    setWoreda(''); setWoredas([]);
    if (!region || !subcity) return;
    api.getWoredas(region, subcity, zone || undefined)
      .then(w => setWoredas((w ?? []) as I18nLocation[]))
      .catch(() => {});
  }, [subcity, region, zone]);

  // Load woredas directly from zone when no subcity level exists in this region
  useEffect(() => {
    if (subcity || !region || !zone) return;
    api.getWoredas(region, '', zone)
      .then(w => { if ((w ?? []).length > 0) setWoredas(w as I18nLocation[]); })
      .catch(() => {});
  }, [zone, region]); // eslint-disable-line react-hooks/exhaustive-deps

  const emit = useCallback((
    r: string, z: string, s: string, w: string,
    k: string, area: string, la: string, lo: string,
  ) => {
    const v: LocationValue = {};
    if (r) v.region = r;
    if (z) v.zone = z;
    if (s) v.subcity = s;
    if (w) v.woreda = w;
    if (k) v.kebele = k;
    if (area) v.specific_area = area;
    const latN = parseFloat(la), lngN = parseFloat(lo);
    if (!isNaN(latN)) v.latitude = latN;
    if (!isNaN(lngN)) v.longitude = lngN;
    onChange(v);
  }, [onChange]);

  useEffect(() => {
    emit(region, zone, subcity, woreda, kebele, specificArea, lat, lng);
  // Debounce free-text fields (kebele, specificArea, lat, lng) to avoid
  // calling onChange on every keystroke — hierarchy selects fire immediately.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [region, zone, subcity, woreda]);

  // Debounced emit for free-text fields — fires 400ms after user stops typing
  useEffect(() => {
    const timer = setTimeout(() => {
      emit(region, zone, subcity, woreda, kebele, specificArea, lat, lng);
    }, 400);
    return () => clearTimeout(timer);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [kebele, specificArea, lat, lng]);

  // Debounced search
  useEffect(() => {
    if (!search.trim()) { setSearchResults([]); setSearchOpen(false); return; }
    const timer = setTimeout(() => {
      api.searchLocations(search)
        .then(r => { setSearchResults(r ?? []); setSearchOpen(true); })
        .catch(() => {});
    }, 280);
    return () => clearTimeout(timer);
  }, [search]);

  const applyResult = (r: any) => {
    setRegion(r.region ?? '');
    setZone(r.zone ?? '');
    setSubcity(r.subcity ?? '');
    setWoreda(r.woreda ?? '');
    setSearch(''); setSearchResults([]); setSearchOpen(false);
  };

  const clearAll = () => {
    setRegion(''); setZone(''); setSubcity(''); setWoreda('');
    setKebele(''); setSpecificArea(''); setLat(''); setLng('');
    setSearch(''); setSearchResults([]); setSearchOpen(false);
  };

  const hasSelection = region || zone || subcity || woreda;

  return (
    <div className="flex flex-col gap-2">

      {/* Search box */}
      {searchable && (
        <div className="relative">
          <Search size={13} className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none" />
          <input
            type="text"
            value={search}
            onChange={e => setSearch(e.target.value)}
            onFocus={() => searchResults.length > 0 && setSearchOpen(true)}
            placeholder={t.ui?.search_locations ?? 'Search location…'}
            className={`${SEL} pl-9 pr-8`}
          />
          {search && (
            <button type="button" onClick={() => { setSearch(''); setSearchResults([]); setSearchOpen(false); }}
              className="absolute right-2.5 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-600">
              <X size={13} />
            </button>
          )}
          {searchOpen && searchResults.length > 0 && (
            <div className="absolute z-30 w-full mt-1 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded-xl shadow-2xl overflow-hidden max-h-56 overflow-y-auto">
              {searchResults.map((r, i) => (
                <button type="button" key={i} onClick={() => applyResult(r)}
                  className="w-full text-left px-3 py-2.5 text-sm hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors border-b border-zinc-100 dark:border-zinc-800 last:border-0">
                  {locPath(r, locale)}
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Region */}
      <select className={SEL} value={region} onChange={e => setRegion(e.target.value)}>
        <option value="">{t.ui?.all_regions ?? 'All Regions'}</option>
        {regions.map(r => (
          <option key={r.en} value={r.en}>{locName(r, locale)}</option>
        ))}
      </select>

      {/* Zone */}
      {zones.length > 0 && (
        <select className={SEL} value={zone} onChange={e => setZone(e.target.value)}>
          <option value="">{t.ui?.all_zones ?? 'All Zones'}</option>
          {zones.map(z => (
            <option key={z.en} value={z.en}>{locName(z, locale)}</option>
          ))}
        </select>
      )}

      {/* Subcity */}
      {subcities.length > 0 && (
        <select className={SEL} value={subcity} onChange={e => setSubcity(e.target.value)}>
          <option value="">{t.ui?.all_subcities ?? 'All Subcities'}</option>
          {subcities.map(s => (
            <option key={s.en} value={s.en}>{locName(s, locale)}</option>
          ))}
        </select>
      )}

      {/* Woreda */}
      {woredas.length > 0 && (
        <select className={SEL} value={woreda} onChange={e => setWoreda(e.target.value)}>
          <option value="">{t.ui?.all_woredas ?? 'All Woredas'}</option>
          {woredas.map(w => (
            <option key={w.en} value={w.en}>{locName(w, locale)}</option>
          ))}
        </select>
      )}

      {/* Kebele — free text, shown after woreda is selected */}
      {showPrecise && woreda && (
        <input
          type="text"
          value={kebele}
          onChange={e => setKebele(e.target.value)}
          placeholder={t.ui?.kebele_placeholder ?? 'Kebele (optional, e.g. Kebele 01)'}
          className={INP}
        />
      )}

      {/* Specific area — free text description */}
      {showPrecise && region && (
        <input
          type="text"
          value={specificArea}
          onChange={e => setSpecificArea(e.target.value)}
          placeholder={t.ui?.specific_area_placeholder ?? 'Specific area / landmark (optional)'}
          className={INP}
        />
      )}

      {/* GPS coordinates — gold/premium only */}
      {showCoords ? (
        <div className="flex flex-col gap-1.5 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-xl p-3">
          <p className="text-[10px] font-black uppercase tracking-widest text-yellow-700 dark:text-yellow-400 flex items-center gap-1">
            <MapPin size={10} /> {t.ui?.gps_label ?? 'GPS Coordinates (Gold / Premium)'}
          </p>
          <div className="flex gap-2">
            <input
              type="number"
              step="0.0000001"
              value={lat}
              onChange={e => setLat(e.target.value)}
              placeholder={t.ui?.latitude_placeholder ?? 'Latitude (e.g. 9.0320)'}
              className={`${INP} flex-1`}
            />
            <input
              type="number"
              step="0.0000001"
              value={lng}
              onChange={e => setLng(e.target.value)}
              placeholder={t.ui?.longitude_placeholder ?? 'Longitude (e.g. 38.7469)'}
              className={`${INP} flex-1`}
            />
          </div>
          <p className="text-[10px] text-yellow-600 dark:text-yellow-500">
            {t.ui?.gps_note ?? 'Buyers with Gold/Premium access will see a map link to your exact location.'}
          </p>
        </div>
      ) : (
        <div className="flex items-center gap-2 bg-zinc-50 dark:bg-zinc-800/60 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2">
          <Lock size={11} className="text-zinc-400 shrink-0" />
          <p className="text-[10px] text-zinc-400">
            {t.ui?.gps_locked ?? 'GPS coordinates available for Gold & Premium sellers'}
          </p>
        </div>
      )}

      {/* Clear */}
      {hasSelection && (
        <button type="button" onClick={clearAll}
          className="flex items-center gap-1.5 text-xs text-zinc-400 hover:text-red-500 transition-colors self-start">
          <X size={11} /> {t.ui?.clear_location ?? 'Clear location'}
        </button>
      )}
    </div>
  );
}
