'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { X, Bell, AlertTriangle, ShieldAlert, CheckCircle2, Clock, ArrowRight, Wrench, DollarSign } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';

export interface NotificationItem {
  id: string;
  title: string;
  description: string;
  timestamp: string;
  type: 'maintenance' | 'escrow' | 'moderation' | 'system';
  severity: 'high' | 'medium' | 'low';
  read: boolean;
  actionText?: string;
}

const INITIAL_NOTIFICATIONS: NotificationItem[] = [
  {
    id: 'NOTIF-101',
    title: 'Urgent Maintenance: Elevator Failure in Bole Tower B',
    description: 'Tenant reported elevator breakdown on 4th floor. Requires immediate contractor assignment.',
    timestamp: '10 minutes ago',
    type: 'maintenance',
    severity: 'high',
    read: false,
    actionText: 'Dispatch Technician',
  },
  {
    id: 'NOTIF-102',
    title: 'Escrow Release Hold: Lease #LA-8812',
    description: 'Security deposit payout of $1,500 pending admin approval following property checkout inspection.',
    timestamp: '25 minutes ago',
    type: 'escrow',
    severity: 'medium',
    read: false,
    actionText: 'Review Inspection',
  },
  {
    id: 'NOTIF-103',
    title: 'Flagged Listing: Unverified Seller Account',
    description: 'Product SKU #PRD-9912 flagged for counterfeit brand claim. Auto-suspended pending review.',
    timestamp: '1 hour ago',
    type: 'moderation',
    severity: 'high',
    read: false,
    actionText: 'Moderate Listing',
  },
  {
    id: 'NOTIF-104',
    title: 'New High-Value Rental Reservation',
    description: 'Villa #VIL-309 reserved for 6 months ($18,000 upfront escrow deposit received).',
    timestamp: '3 hours ago',
    type: 'system',
    severity: 'low',
    read: true,
  },
];

interface NotificationDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export function NotificationDrawer({ isOpen, onClose }: NotificationDrawerProps) {
  const { user } = useAuth();
  // Start empty — real notifications are loaded from the backend on open
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [filterType, setFilterType] = useState<'all' | 'maintenance' | 'escrow' | 'moderation'>('all');

  const mapBackendNotification = useCallback((n: any): NotificationItem => {
    let type: NotificationItem['type'] = 'system';
    if (n.type === 'maintenance') type = 'maintenance';
    else if (n.type === 'escrow') type = 'escrow';
    else if (n.type === 'moderation') type = 'moderation';
    
    let severity: NotificationItem['severity'] = 'low';
    if (n.severity === 'high') severity = 'high';
    else if (n.severity === 'medium') severity = 'medium';

    return {
      id: `NOTIF-${n.id}`,
      title: n.title || 'Platform Event',
      description: n.message || '',
      timestamp: n.created_at ? new Date(n.created_at).toLocaleTimeString() : 'Recent',
      type,
      severity,
      read: n.is_read || false,
      actionText: n.action_text || undefined,
    };
  }, []);

  const loadNotifications = useCallback(async () => {
    if (!user?.token) return;
    setLoading(true);
    try {
      const data = await api.adminGetNotifications(user.token);
      if (Array.isArray(data)) {
        const mapped = data.map(mapBackendNotification);
        if (mapped.length > 0) {
          setNotifications(mapped);
        } else {
          setNotifications([]);
        }
      }
    } catch (err) {
      console.error('Failed to load admin notifications:', err);
    } finally {
      setLoading(false);
    }
  }, [user, mapBackendNotification]);

  useEffect(() => {
    if (isOpen) {
      loadNotifications();
    }
  }, [isOpen, loadNotifications]);

  if (!isOpen) return null;

  const filtered = notifications.filter((n) => filterType === 'all' || n.type === filterType);
  const unreadCount = notifications.filter((n) => !n.read).length;

  const markAllRead = async () => {
    if (user?.token) {
      try {
        const promises = notifications.filter(n => !n.read).map(async (n) => {
          const numericId = parseInt(n.id.split('-')[1]);
          if (!isNaN(numericId)) {
            await api.adminMarkNotificationRead(numericId, user.token);
          }
        });
        await Promise.all(promises);
        loadNotifications();
      } catch (err) {
        console.error('Failed to mark all read:', err);
      }
    } else {
      setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    }
  };

  const toggleRead = async (id: string) => {
    const numericId = parseInt(id.split('-')[1]);
    if (isNaN(numericId)) return;

    if (user?.token) {
      try {
        await api.adminMarkNotificationRead(numericId, user.token);
        loadNotifications();
      } catch (err) {
        console.error('Failed to mark read:', err);
      }
    } else {
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, read: !n.read } : n))
      );
    }
  };

  return (
    <div className="fixed inset-0 z-50 overflow-hidden bg-black/60 backdrop-blur-sm transition-opacity">
      <div className="absolute inset-y-0 right-0 max-w-full flex pl-10">
        <div className="w-screen max-w-md bg-zinc-950 border-l border-zinc-800 text-zinc-100 flex flex-col shadow-2xl animate-in slide-in-from-right duration-200">
          
          {/* Header */}
          <div className="p-5 border-b border-zinc-800 flex items-center justify-between bg-zinc-900/60">
            <div className="flex items-center gap-3">
              <div className="p-2.5 bg-orange-600/10 border border-orange-500/20 text-orange-500 rounded-xl">
                <Bell size={20} />
              </div>
              <div>
                <h3 className="text-base font-bold text-white flex items-center gap-2">
                  Notifications & Queue
                  {unreadCount > 0 && (
                    <span className="px-2 py-0.5 text-[11px] font-extrabold bg-orange-600 text-white rounded-full">
                      {unreadCount} new
                    </span>
                  )}
                </h3>
                <p className="text-xs text-zinc-400">Enterprise operational alerts & maintenance queue</p>
              </div>
            </div>

            <button
              onClick={onClose}
              className="text-zinc-400 hover:text-white p-2 rounded-xl hover:bg-zinc-800 transition-colors"
            >
              <X size={20} />
            </button>
          </div>

          {/* Filter Bar */}
          <div className="flex items-center justify-between px-5 py-3 bg-zinc-900/30 border-b border-zinc-800 text-xs">
            <div className="flex items-center gap-1.5 overflow-x-auto py-1">
              {[
                { id: 'all', label: 'All' },
                { id: 'maintenance', label: 'Maintenance', icon: Wrench },
                { id: 'escrow', label: 'Escrow', icon: DollarSign },
                { id: 'moderation', label: 'Moderation', icon: ShieldAlert },
              ].map((tab) => {
                const Icon = tab.icon;
                const active = filterType === tab.id;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setFilterType(tab.id as any)}
                    className={`flex items-center gap-1 px-2.5 py-1 rounded-lg font-semibold transition-all ${
                      active
                        ? 'bg-orange-500/20 text-orange-400 border border-orange-500/30'
                        : 'text-zinc-400 hover:text-white hover:bg-zinc-800'
                    }`}
                  >
                    {Icon && <Icon size={13} />}
                    {tab.label}
                  </button>
                );
              })}
            </div>

            {unreadCount > 0 && (
              <button
                onClick={markAllRead}
                className="text-[11px] font-bold text-orange-400 hover:underline shrink-0 ml-2"
              >
                Mark all read
              </button>
            )}
          </div>

          {/* List */}
          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {loading ? (
              <div className="py-20 flex items-center justify-center">
                <Spinner />
              </div>
            ) : filtered.length === 0 ? (
              <div className="py-20 text-center text-zinc-500 text-sm">
                No notifications in this queue.
              </div>
            ) : (
              filtered.map((item) => {
                let SeverityIcon = CheckCircle2;
                let severityBg = 'border-zinc-800 bg-zinc-900/50';

                if (item.severity === 'high') {
                  SeverityIcon = AlertTriangle;
                  severityBg = 'border-red-500/30 bg-red-950/20';
                } else if (item.severity === 'medium') {
                  SeverityIcon = Clock;
                  severityBg = 'border-amber-500/30 bg-amber-950/20';
                }

                return (
                  <div
                    key={item.id}
                    className={`p-4 rounded-2xl border ${severityBg} ${
                      !item.read ? 'ring-1 ring-orange-500/40' : 'opacity-80'
                    } transition-all relative group`}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="flex items-center gap-2">
                        {item.type === 'maintenance' && <Wrench size={16} className="text-amber-400" />}
                        {item.type === 'escrow' && <DollarSign size={16} className="text-emerald-400" />}
                        {item.type === 'moderation' && <ShieldAlert size={16} className="text-red-400" />}
                        {item.type === 'system' && <CheckCircle2 size={16} className="text-blue-400" />}

                        <span className="font-bold text-sm text-zinc-100">{item.title}</span>
                      </div>

                      <button
                        onClick={() => toggleRead(item.id)}
                        className={`w-2.5 h-2.5 rounded-full shrink-0 ${
                          item.read ? 'bg-zinc-700' : 'bg-orange-500 ring-4 ring-orange-500/20'
                        }`}
                        title={item.read ? 'Mark unread' : 'Mark read'}
                      />
                    </div>

                    <p className="text-xs text-zinc-400 mt-2 leading-relaxed">{item.description}</p>

                    <div className="flex items-center justify-between mt-3 pt-3 border-t border-zinc-800/60">
                      <span className="text-[10px] text-zinc-500 font-mono">{item.timestamp}</span>

                      {item.actionText && (
                        <button className="flex items-center gap-1 text-xs font-bold text-orange-400 hover:text-orange-300 transition-colors">
                          <span>{item.actionText}</span>
                          <ArrowRight size={13} />
                        </button>
                      )}
                    </div>
                  </div>
                );
              })
            )}
          </div>

          {/* Footer */}
          <div className="p-4 bg-zinc-900 border-t border-zinc-800 text-center text-xs text-zinc-500">
            Automated Escalation System • Priority Queue Active
          </div>
        </div>
      </div>
    </div>
  );
}
