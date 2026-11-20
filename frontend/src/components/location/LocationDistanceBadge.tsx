'use client';

import { useState, useEffect } from 'react';
import { MapPin, Navigation, Loader2 } from 'lucide-react';

interface Props {
  latitude?: number | null;
  longitude?: number | null;
  locationText: string;
  className?: string;
  interactive?: boolean;
}

// Haversine distance in kilometers
function calculateHaversine(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371; // Earth radius in km
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

export function LocationDistanceBadge({
  latitude,
  longitude,
  locationText,
  className = '',
  interactive = true,
}: Props) {
  const [distanceKm, setDistanceKm] = useState<number | null>(null);
  const [checking, setChecking] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const hasCoords =
    latitude != null &&
    longitude != null &&
    !isNaN(Number(latitude)) &&
    !isNaN(Number(longitude)) &&
    Number(latitude) >= -90 &&
    Number(latitude) <= 90 &&
    Number(longitude) >= -180 &&
    Number(longitude) <= 180;

  const checkDistance = () => {
    if (!hasCoords) return;
    if (typeof window === 'undefined' || !navigator.geolocation) {
      setErrorMsg('Geolocation not supported');
      return;
    }
    if (window.isSecureContext === false && window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1') {
      setErrorMsg('GPS requires HTTPS on mobile');
      return;
    }
    setChecking(true);
    setErrorMsg(null);

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const userLat = pos.coords.latitude;
        const userLng = pos.coords.longitude;
        const dist = calculateHaversine(userLat, userLng, Number(latitude), Number(longitude));
        setDistanceKm(Math.round(dist * 10) / 10);
        setChecking(false);
      },
      (err) => {
        setChecking(false);
        if (err.code === err.PERMISSION_DENIED) {
          setErrorMsg('Location permission required');
        } else {
          setErrorMsg('Unable to get distance');
        }
      },
      { timeout: 8000, enableHighAccuracy: true }
    );
  };

  // If coordinates exist, try gentle passive check if permission was already granted
  useEffect(() => {
    if (!hasCoords || typeof window === 'undefined') return;
    if ('permissions' in navigator) {
      navigator.permissions.query({ name: 'geolocation' as PermissionName }).then((res) => {
        if (res.state === 'granted') {
          navigator.geolocation.getCurrentPosition(
            (pos) => {
              const dist = calculateHaversine(pos.coords.latitude, pos.coords.longitude, Number(latitude), Number(longitude));
              setDistanceKm(Math.round(dist * 10) / 10);
            },
            () => {},
            { timeout: 5000 }
          );
        }
      }).catch(() => {});
    }
  }, [hasCoords, latitude, longitude]);

  const mapsUrl = hasCoords
    ? `https://www.google.com/maps/dir/?api=1&destination=${latitude},${longitude}`
    : `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(locationText + ', Ethiopia')}`;

  return (
    <div className={`inline-flex items-center gap-1.5 flex-wrap ${className}`}>
      <button
        type="button"
        onClick={interactive && hasCoords ? checkDistance : undefined}
        title={hasCoords ? 'Click to check distance from your device' : locationText}
        className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-xl text-xs font-semibold transition-all ${
          hasCoords
            ? 'bg-blue-50 hover:bg-blue-100 dark:bg-blue-950/40 dark:hover:bg-blue-900/40 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800/60 cursor-pointer shadow-xs active:scale-95'
            : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300 border border-zinc-200 dark:border-zinc-700'
        }`}
      >
        <MapPin size={12} className={hasCoords ? 'text-[#1769E0] shrink-0' : 'text-zinc-500 shrink-0'} />
        <span className="truncate max-w-[200px]">{locationText || 'Ethiopia'}</span>

        {checking && <Loader2 size={11} className="animate-spin text-[#1769E0] ml-1" />}

        {distanceKm != null && (
          <span className="ml-1 px-1.5 py-0.5 rounded-md bg-emerald-500 text-white text-[10px] font-black tracking-tight shrink-0 shadow-xs">
            {distanceKm < 1 ? `${Math.round(distanceKm * 1000)}m away` : `${distanceKm} km away`}
          </span>
        )}

        {hasCoords && distanceKm == null && !checking && interactive && (
          <span className="ml-1 text-[10px] font-bold text-blue-500 hover:underline">
            Check Distance
          </span>
        )}
      </button>

      {hasCoords && (
        <a
          href={mapsUrl}
          target="_blank"
          rel="noopener noreferrer"
          title="Open directions in Google Maps"
          className="p-1 text-zinc-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
        >
          <Navigation size={12} />
        </a>
      )}

      {errorMsg && (
        <span className="text-[10px] text-amber-500 font-medium">({errorMsg})</span>
      )}
    </div>
  );
}
