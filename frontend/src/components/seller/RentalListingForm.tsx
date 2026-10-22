'use client';
import { useState } from 'react';
import { Clock, DollarSign, Calendar, Shield, AlertCircle } from 'lucide-react';

interface RentalFormData {
  rental_price: string;
  rental_duration_value: string;
  rental_duration_unit: 'hour' | 'day' | 'week' | 'month' | 'year';
  deposit_amount: string;
  available_from: string;
  available_until: string;
}

interface Props {
  data: RentalFormData;
  onChange: (field: keyof RentalFormData, value: string) => void;
  errors?: Partial<Record<keyof RentalFormData, string>>;
}

const DURATION_UNITS = [
  { value: 'hour', label: 'Hour(s)', icon: '⏰' },
  { value: 'day', label: 'Day(s)', icon: '📅' },
  { value: 'week', label: 'Week(s)', icon: '📆' },
  { value: 'month', label: 'Month(s)', icon: '🗓️' },
  { value: 'year', label: 'Year(s)', icon: '📋' },
] as const;

export function RentalListingForm({ data, onChange, errors }: Props) {
  const [showDepositInfo, setShowDepositInfo] = useState(false);

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl px-4 py-3">
        <div className="flex items-center gap-2 text-blue-700 dark:text-blue-400">
          <Clock size={16} />
          <p className="text-sm font-bold">Rental Listing</p>
        </div>
        <p className="text-xs text-blue-600/70 dark:text-blue-400/70 mt-1">
          Set your rental price, duration, and availability dates
        </p>
      </div>

      {/* Rental Price & Duration */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Rental Price */}
        <div>
          <label className="flex items-center gap-1.5 text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2">
            <DollarSign size={14} />
            Rental Price (ETB) *
          </label>
          <input
            type="number"
            min="1"
            step="0.01"
            value={data.rental_price}
            onChange={(e) => onChange('rental_price', e.target.value)}
            placeholder="e.g., 500"
            className={`w-full border rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 transition-all ${
              errors?.rental_price
                ? 'border-red-300 focus:ring-red-500/30 focus:border-red-500'
                : 'border-zinc-200 dark:border-zinc-700 focus:ring-blue-500/30 focus:border-blue-500'
            } bg-white dark:bg-zinc-800`}
          />
          {errors?.rental_price && (
            <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
              <AlertCircle size={10} /> {errors.rental_price}
            </p>
          )}
          <p className="text-xs text-zinc-400 mt-1">Price per rental period</p>
        </div>

        {/* Duration */}
        <div>
          <label className="flex items-center gap-1.5 text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2">
            <Clock size={14} />
            Rental Duration *
          </label>
          <div className="grid grid-cols-2 gap-2">
            <input
              type="number"
              min="1"
              value={data.rental_duration_value}
              onChange={(e) => onChange('rental_duration_value', e.target.value)}
              placeholder="e.g., 1"
              className={`border rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 transition-all ${
                errors?.rental_duration_value
                  ? 'border-red-300 focus:ring-red-500/30'
                  : 'border-zinc-200 dark:border-zinc-700 focus:ring-blue-500/30'
              } bg-white dark:bg-zinc-800`}
            />
            <select
              value={data.rental_duration_unit}
              onChange={(e) => onChange('rental_duration_unit', e.target.value)}
              className="border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500/30 bg-white dark:bg-zinc-800"
            >
              {DURATION_UNITS.map((unit) => (
                <option key={unit.value} value={unit.value}>
                  {unit.icon} {unit.label}
                </option>
              ))}
            </select>
          </div>
          {errors?.rental_duration_value && (
            <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
              <AlertCircle size={10} /> {errors.rental_duration_value}
            </p>
          )}
          <p className="text-xs text-zinc-400 mt-1">
            e.g., "1 Day" = rent for 1 day at the price above
          </p>
        </div>
      </div>

      {/* Deposit Amount */}
      <div>
        <label className="flex items-center gap-1.5 text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2">
          <Shield size={14} />
          Security Deposit (ETB)
          <button
            type="button"
            onClick={() => setShowDepositInfo(!showDepositInfo)}
            className="text-blue-500 hover:text-blue-600 text-xs"
          >
            ℹ️
          </button>
        </label>
        <input
          type="number"
          min="0"
          step="0.01"
          value={data.deposit_amount}
          onChange={(e) => onChange('deposit_amount', e.target.value)}
          placeholder="e.g., 1000 (optional)"
          className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500/30 focus:border-blue-500 bg-white dark:bg-zinc-800"
        />
        {showDepositInfo && (
          <div className="mt-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg px-3 py-2">
            <p className="text-xs text-blue-700 dark:text-blue-400">
              Security deposit is refundable and protects against damage. Leave at 0 if not required.
            </p>
          </div>
        )}
        <p className="text-xs text-zinc-400 mt-1">
          Refundable amount held during rental (optional)
        </p>
      </div>

      {/* Availability Dates */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Available From */}
        <div>
          <label className="flex items-center gap-1.5 text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2">
            <Calendar size={14} />
            Available From *
          </label>
          <input
            type="date"
            value={data.available_from}
            onChange={(e) => onChange('available_from', e.target.value)}
            min={new Date().toISOString().split('T')[0]}
            className={`w-full border rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 transition-all ${
              errors?.available_from
                ? 'border-red-300 focus:ring-red-500/30'
                : 'border-zinc-200 dark:border-zinc-700 focus:ring-blue-500/30'
            } bg-white dark:bg-zinc-800`}
          />
          {errors?.available_from && (
            <p className="text-xs text-red-500 mt-1 flex items-center gap-1">
              <AlertCircle size={10} /> {errors.available_from}
            </p>
          )}
          <p className="text-xs text-zinc-400 mt-1">When item becomes available</p>
        </div>

        {/* Available Until */}
        <div>
          <label className="flex items-center gap-1.5 text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2">
            <Calendar size={14} />
            Available Until (Optional)
          </label>
          <input
            type="date"
            value={data.available_until}
            onChange={(e) => onChange('available_until', e.target.value)}
            min={data.available_from || new Date().toISOString().split('T')[0]}
            className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500/30 focus:border-blue-500 bg-white dark:bg-zinc-800"
          />
          <p className="text-xs text-zinc-400 mt-1">
            Leave empty if available indefinitely
          </p>
        </div>
      </div>

      {/* Example Calculation */}
      {data.rental_price && data.rental_duration_value && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl px-4 py-3">
          <p className="text-xs font-bold text-green-700 dark:text-green-400 mb-1">
            Example Rental
          </p>
          <p className="text-xs text-green-600 dark:text-green-300">
            Customer rents for {data.rental_duration_value}{' '}
            {DURATION_UNITS.find((u) => u.value === data.rental_duration_unit)?.label.toLowerCase()} =
            ETB {parseFloat(data.rental_price || '0').toLocaleString()}
            {data.deposit_amount && parseFloat(data.deposit_amount) > 0 && (
              <span> + ETB {parseFloat(data.deposit_amount).toLocaleString()} deposit (refundable)</span>
            )}
          </p>
        </div>
      )}
    </div>
  );
}
