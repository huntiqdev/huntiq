'use client';
import { useCallback, useRef } from 'react';

/**
 * Plays a short notification beep using the Web Audio API.
 * AudioContext is created lazily on first play() call after a user gesture,
 * avoiding the Chrome autoplay policy warning.
 */
export function useNotificationSound() {
  const ctxRef = useRef<AudioContext | null>(null);

  const play = useCallback((type: 'alert' | 'success' | 'info' = 'alert') => {
    try {
      if (!ctxRef.current || ctxRef.current.state === 'closed') {
        ctxRef.current = new (window.AudioContext || (window as any).webkitAudioContext)();
      }
      const ctx = ctxRef.current;

      // If still suspended (no user gesture yet), skip silently — don't resume
      // forcefully, which is what triggers the Chrome warning.
      if (ctx.state === 'suspended') return;

      const configs: Record<string, { freqs: number[]; duration: number; gain: number }> = {
        alert:   { freqs: [880, 1100], duration: 0.12, gain: 0.35 },
        success: { freqs: [660, 880],  duration: 0.10, gain: 0.25 },
        info:    { freqs: [660],       duration: 0.15, gain: 0.20 },
      };

      const { freqs, duration, gain } = configs[type];
      const now = ctx.currentTime;

      freqs.forEach((freq, i) => {
        const osc = ctx.createOscillator();
        const gainNode = ctx.createGain();
        osc.connect(gainNode);
        gainNode.connect(ctx.destination);

        osc.type = 'sine';
        osc.frequency.setValueAtTime(freq, now + i * (duration + 0.02));

        gainNode.gain.setValueAtTime(0, now + i * (duration + 0.02));
        gainNode.gain.linearRampToValueAtTime(gain, now + i * (duration + 0.02) + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.001, now + i * (duration + 0.02) + duration);

        osc.start(now + i * (duration + 0.02));
        osc.stop(now + i * (duration + 0.02) + duration + 0.01);
      });
    } catch {
      // Web Audio not available (SSR, old browser) — silently ignore
    }
  }, []);

  // Call this once on any user click to unlock the AudioContext for future plays
  const unlock = useCallback(() => {
    try {
      if (!ctxRef.current || ctxRef.current.state === 'closed') {
        ctxRef.current = new (window.AudioContext || (window as any).webkitAudioContext)();
      }
      if (ctxRef.current.state === 'suspended') {
        ctxRef.current.resume().catch(() => {});
      }
    } catch {
      // ignore
    }
  }, []);

  return { play, unlock };
}
