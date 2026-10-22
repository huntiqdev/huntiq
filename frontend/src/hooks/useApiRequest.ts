'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { ApiError } from '@/lib/api/core';

interface UseApiRequestOptions {
  onSuccess?: (data: any) => void;
  onError?: (error: ApiError) => void;
  showSuccessToast?: boolean;
  showErrorToast?: boolean;
  successMessage?: string;
}

/**
 * Hook for handling API requests with loading states and automatic delays
 * - Shows loading for minimum 1 second
 * - Success/error messages visible for 3 seconds
 * - Automatically cancels in-flight requests on unmount
 */
export function useApiRequest<T = any>(options: UseApiRequestOptions = {}) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<T | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const isMountedRef = useRef(true);

  // Cleanup on unmount - cancel any in-flight requests
  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  const execute = useCallback(
    async (apiCall: () => Promise<T>) => {
      // Cancel previous request if still in progress
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      // Create new abort controller for this request
      abortControllerRef.current = new AbortController();
      const signal = abortControllerRef.current.signal;

      setLoading(true);
      setError(null);
      
      try {
        // API call already has minimum 1-second delay built in
        const result = await apiCall();
        
        // Only update state if component is still mounted and request wasn't aborted
        if (!signal.aborted && isMountedRef.current) {
          setData(result);
          
          if (options.onSuccess) {
            options.onSuccess(result);
          }
        }
        
        return result;
      } catch (err: any) {
        // Ignore abort errors - they're expected when cancelling
        if (err.name === 'AbortError' || signal.aborted) {
          return null as T;
        }

        const errorMessage = err instanceof ApiError 
          ? err.message 
          : 'An unexpected error occurred';
        
        // Only update state if component is still mounted
        if (isMountedRef.current) {
          setError(errorMessage);
          
          if (options.onError && err instanceof ApiError) {
            options.onError(err);
          }
          
          // Keep error visible for 3 seconds
          setTimeout(() => {
            if (isMountedRef.current) {
              setError(null);
            }
          }, 3000);
        }
        
        throw err;
      } finally {
        if (isMountedRef.current) {
          setLoading(false);
        }
      }
    },
    [options]
  );

  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
    setData(null);
  }, []);

  return {
    loading,
    error,
    data,
    execute,
    reset,
  };
}
