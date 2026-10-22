'use client';

import React, { useState } from 'react';
import { Calendar, Filter, ChevronLeft, ChevronRight, Lock, User, Sparkles, AlertCircle, Wrench, Check } from 'lucide-react';

export interface PropertyUnit {
  id: string;
  name: string;
  type: string;
  ratePerNight: number;
  location: string;
  schedule: Record<string, 'available' | 'occupied' | 'reserved' | 'maintenance'>;
  tenantDetails?: Record<string, { name: string; leaseId: string; status: string }>;
}

const DEFAULT_UNITS: PropertyUnit[] = [
  {
    id: 'U-101',
    name: 'Bole Luxury Penthouse 401',
    type: 'Penthouse',
    ratePerNight: 220,
    location: 'Bole Hub',
    schedule: {
      'Day 1': 'occupied', 'Day 2': 'occupied', 'Day 3': 'occupied', 'Day 4': 'occupied',
      'Day 5': 'available', 'Day 6': 'available', 'Day 7': 'reserved', 'Day 8': 'reserved',
      'Day 9': 'available', 'Day 10': 'available', 'Day 11': 'maintenance', 'Day 12': 'available'
    },
    tenantDetails: {
      'Day 1': { name: 'Meron Tadesse', leaseId: 'LA-9912', status: 'Checked In' },
      'Day 7': { name: 'David Smith', leaseId: 'LA-1044', status: 'Deposit Paid' }
    }
  },
  {
    id: 'U-102',
    name: 'Kazanchis Executive Suite 202',
    type: 'Suite',
    ratePerNight: 150,
    location: 'Central Addis',
    schedule: {
      'Day 1': 'available', 'Day 2': 'available', 'Day 3': 'reserved', 'Day 4': 'reserved',
      'Day 5': 'occupied', 'Day 6': 'occupied', 'Day 7': 'occupied', 'Day 8': 'available',
      'Day 9': 'available', 'Day 10': 'maintenance', 'Day 11': 'maintenance', 'Day 12': 'available'
    },
    tenantDetails: {
      'Day 3': { name: 'Abebe Bikila', leaseId: 'LA-4022', status: 'Confirmed' },
      'Day 5': { name: 'Sarah Jenkins', leaseId: 'LA-8810', status: 'Checked In' }
    }
  },
  {
    id: 'U-103',
    name: 'Hawassa Lakeview Villa 1',
    type: 'Villa',
    ratePerNight: 350,
    location: 'Hawassa Resort',
    schedule: {
      'Day 1': 'occupied', 'Day 2': 'occupied', 'Day 3': 'occupied', 'Day 4': 'occupied',
      'Day 5': 'occupied', 'Day 6': 'occupied', 'Day 7': 'occupied', 'Day 8': 'occupied',
      'Day 9': 'reserved', 'Day 10': 'reserved', 'Day 11': 'available', 'Day 12': 'available'
    },
    tenantDetails: {
      'Day 1': { name: 'Global Tech Corp', leaseId: 'LA-7700', status: 'Corporate Lease' },
      'Day 9': { name: 'Helen Haile', leaseId: 'LA-3021', status: 'Deposit Paid' }
    }
  }
];

const DAYS = ['Day 1', 'Day 2', 'Day 3', 'Day 4', 'Day 5', 'Day 6', 'Day 7', 'Day 8', 'Day 9', 'Day 10', 'Day 11', 'Day 12'];

interface AvailabilityGridProps {
  units?: PropertyUnit[];
}

export function AvailabilityGrid({ units }: AvailabilityGridProps) {
  const [selectedLocation, setSelectedLocation] = useState<string>('all');
  const [activeCell, setActiveCell] = useState<{ unitId: string; day: string } | null>(null);
  const [actionSuccess, setActionSuccess] = useState<string | null>(null);

  const displayUnits = units && units.length > 0 ? units : DEFAULT_UNITS;

  const filteredUnits = displayUnits.filter(
    (u) => selectedLocation === 'all' || u.location.toLowerCase().includes(selectedLocation.toLowerCase())
  );

  const handleCellClick = (unit: PropertyUnit, day: string) => {
    setActiveCell({ unitId: unit.id, day });
  };

  const handleSetStatus = (newStatus: 'available' | 'occupied' | 'reserved' | 'maintenance') => {
    if (!activeCell) return;
    setActionSuccess(`Unit ${activeCell.unitId} status on ${activeCell.day} updated to ${newStatus.toUpperCase()}`);
    setTimeout(() => setActionSuccess(null), 3000);
    setActiveCell(null);
  };

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4">
      {/* Grid Top Controls */}
      <div className="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-zinc-800">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
            <Calendar size={20} />
          </div>
          <div>
            <h2 className="text-base font-extrabold text-white flex items-center gap-2">
              Property Availability & Room Allocation Matrix
              <span className="px-2 py-0.5 text-[10px] font-mono bg-blue-500/10 text-blue-400 border border-blue-500/20 rounded-md">
                12-Day View
              </span>
            </h2>
            <p className="text-xs text-zinc-400">Real-time room occupancy, check-ins, and maintenance holds</p>
          </div>
        </div>

        {/* Filters & Legends */}
        <div className="flex flex-wrap items-center gap-3">
          <div className="flex items-center gap-1.5 text-xs text-zinc-400 bg-zinc-950 p-1.5 rounded-xl border border-zinc-800">
            <span className="px-2 py-1 bg-emerald-500/20 text-emerald-400 rounded-md font-bold text-[10px]">🟢 Available</span>
            <span className="px-2 py-1 bg-blue-500/20 text-blue-400 rounded-md font-bold text-[10px]">🔵 Rented/Occupied</span>
            <span className="px-2 py-1 bg-amber-500/20 text-amber-400 rounded-md font-bold text-[10px]">🟡 Reserved</span>
            <span className="px-2 py-1 bg-purple-500/20 text-purple-400 rounded-md font-bold text-[10px]">🟣 Maintenance</span>
          </div>

          <select
            value={selectedLocation}
            onChange={(e) => setSelectedLocation(e.target.value)}
            className="bg-zinc-950 border border-zinc-800 text-zinc-200 text-xs font-bold rounded-xl px-3 py-2 focus:outline-none focus:border-orange-500"
          >
            <option value="all">All Locations</option>
            <option value="Bole">Bole Hub</option>
            <option value="Hawassa">Hawassa Resort</option>
            <option value="Dire Dawa">Dire Dawa</option>
          </select>
        </div>
      </div>

      {actionSuccess && (
        <div className="p-3 bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 rounded-xl text-xs font-bold flex items-center gap-2 animate-in fade-in">
          <Check size={16} />
          {actionSuccess}
        </div>
      )}

      {/* Grid Timeline Matrix */}
      <div className="overflow-x-auto">
        <table className="w-full border-collapse text-xs text-left">
          <thead>
            <tr className="bg-zinc-950 border-b border-zinc-800 text-zinc-400 font-mono text-[11px]">
              <th className="p-3 sticky left-0 z-10 bg-zinc-950 min-w-[220px] font-bold text-zinc-200 border-r border-zinc-800">
                Property / Unit Name
              </th>
              {DAYS.map((day) => (
                <th key={day} className="p-2 text-center min-w-[64px] border-r border-zinc-800/60 font-semibold">
                  {day}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-800/60">
            {filteredUnits.map((unit) => (
              <tr key={unit.id} className="hover:bg-zinc-950/40 transition-colors">
                {/* Fixed Left Unit Info */}
                <td className="p-3 sticky left-0 z-10 bg-zinc-900 border-r border-zinc-800 min-w-[220px]">
                  <div className="font-bold text-white text-xs truncate">{unit.name}</div>
                  <div className="flex items-center justify-between text-[10px] text-zinc-400 mt-0.5">
                    <span className="px-1.5 py-0.5 bg-zinc-800 rounded font-mono">{unit.type}</span>
                    <span className="font-mono text-orange-400 font-bold">${unit.ratePerNight}/night</span>
                  </div>
                </td>

                {/* Day Matrix Cells */}
                {DAYS.map((day) => {
                  const status = unit.schedule[day] || 'available';
                  const tenantInfo = unit.tenantDetails?.[day];

                  let colorClass = 'bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 border-emerald-500/20';
                  let label = 'Free';

                  if (status === 'occupied') {
                    colorClass = 'bg-blue-600/20 text-blue-300 hover:bg-blue-600/30 border-blue-500/30 font-bold';
                    label = 'Rented';
                  } else if (status === 'reserved') {
                    colorClass = 'bg-amber-500/20 text-amber-300 hover:bg-amber-500/30 border-amber-500/30';
                    label = 'Reserved';
                  } else if (status === 'maintenance') {
                    colorClass = 'bg-purple-500/20 text-purple-300 hover:bg-purple-500/30 border-purple-500/30';
                    label = 'Hold';
                  }

                  return (
                    <td key={day} className="p-1 border-r border-zinc-800/40 text-center">
                      <button
                        onClick={() => handleCellClick(unit, day)}
                        className={`w-full py-2.5 rounded-lg border text-[10px] transition-all flex flex-col items-center justify-center gap-0.5 group relative ${colorClass}`}
                        title={`${unit.name} - ${day}: ${status.toUpperCase()} ${tenantInfo ? `(${tenantInfo.name})` : ''}`}
                      >
                        <span className="leading-none">{label}</span>

                        {/* Status Icon */}
                        {status === 'occupied' && <User size={10} className="text-blue-400" />}
                        {status === 'maintenance' && <Wrench size={10} className="text-purple-400" />}
                        {status === 'reserved' && <Lock size={10} className="text-amber-400" />}
                      </button>
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Cell Quick Action Dialog */}
      {activeCell && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
          <div className="w-full max-w-sm bg-zinc-900 border border-zinc-800 rounded-2xl p-5 shadow-2xl space-y-4 animate-in zoom-in-95">
            <div className="flex items-center justify-between border-b border-zinc-800 pb-3">
              <h3 className="font-extrabold text-sm text-white">Modify Cell Allocation</h3>
              <button onClick={() => setActiveCell(null)} className="text-zinc-500 hover:text-white">✕</button>
            </div>

            <div className="text-xs text-zinc-300 space-y-1 font-mono">
              <p>Unit ID: <span className="text-orange-400 font-bold">{activeCell.unitId}</span></p>
              <p>Timeline Column: <span className="text-zinc-100 font-bold">{activeCell.day}</span></p>
            </div>

            <div className="space-y-2 pt-2">
              <p className="text-[11px] font-bold text-zinc-400 uppercase tracking-wider">Set Instant Unit Status:</p>
              <div className="grid grid-cols-2 gap-2 text-xs font-bold">
                <button
                  onClick={() => handleSetStatus('available')}
                  className="p-2.5 bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 rounded-xl hover:bg-emerald-500/20"
                >
                  🟢 Available
                </button>
                <button
                  onClick={() => handleSetStatus('reserved')}
                  className="p-2.5 bg-amber-500/10 border border-amber-500/30 text-amber-400 rounded-xl hover:bg-amber-500/20"
                >
                  🟡 Reserved
                </button>
                <button
                  onClick={() => handleSetStatus('occupied')}
                  className="p-2.5 bg-blue-500/10 border border-blue-500/30 text-blue-400 rounded-xl hover:bg-blue-500/20"
                >
                  🔵 Rented
                </button>
                <button
                  onClick={() => handleSetStatus('maintenance')}
                  className="p-2.5 bg-purple-500/10 border border-purple-500/30 text-purple-400 rounded-xl hover:bg-purple-500/20"
                >
                  🟣 Maintenance
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
