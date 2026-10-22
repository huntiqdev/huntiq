/**
 * Security configuration for the frontend application.
 * These values control timing, retries, and security-sensitive behaviors.
 * 
 * SECURITY NOTE: These are defaults. Some can be overridden by backend settings.
 */

export interface SecurityConfig {
  // JWT & Token Management
  jwtGracePeriodMs: number;        // Grace period for token expiry (clock skew tolerance)
  sessionCheckIntervalMs: number;   // How often to check admin permissions
  redirectThrottleMs: number;       // Minimum time between 401 redirects
  
  // OTP
  otpLength: number;                // Expected OTP code length
  otpExpirySeconds: number;         // OTP validity period (from backend)
  otpCountdownWarningSeconds: number; // Show warning when countdown reaches this
  
  // Network & Retries
  retryMaxAttempts: number;         // Max retry attempts for failed requests
  retryDelayMs: number;             // Base delay between retries (exponential backoff)
  retryableStatusCodes: number[];   // HTTP status codes that trigger retry
  minLoadingTimeMs: number;         // Minimum loading time for UX consistency
  
  // File Upload
  maxImageSizeMB: number;           // Max image file size (from backend)
  acceptedImageMimes: string[];     // Allowed image MIME types
  
  // Referral
  referralCodeLength: number;       // Expected referral code length
  referralCodePattern: RegExp;      // Validation pattern for referral codes
  
  // Debounce Timings (UX)
  searchDebounceMs: number;         // Debounce for search input
  locationDebounceMs: number;       // Debounce for location search
  formSaveDebounceMs: number;       // Debounce for auto-save forms
}

/**
 * Default security configuration.
 * These values are used as fallbacks if backend settings are unavailable.
 */
export const DEFAULT_SECURITY_CONFIG: SecurityConfig = {
  // Token management
  jwtGracePeriodMs: 10_000,        // 10 seconds (matches backend)
  sessionCheckIntervalMs: 30_000,   // 30 seconds
  redirectThrottleMs: 5_000,        // 5 seconds
  
  // OTP
  otpLength: 6,
  otpExpirySeconds: 300,            // 5 minutes default, override from backend
  otpCountdownWarningSeconds: 60,   // Show warning at 1 minute remaining
  
  // Network
  retryMaxAttempts: 2,
  retryDelayMs: 1_000,              // 1 second base delay
  retryableStatusCodes: [408, 429, 500, 502, 503, 504],
  minLoadingTimeMs: 1_000,          // 1 second minimum loading
  
  // File upload
  maxImageSizeMB: 20,               // 20 MB default, override from backend
  acceptedImageMimes: ['image/jpeg', 'image/png', 'image/webp'],
  
  // Referral
  referralCodeLength: 11,           // REF + 8 chars = 11
  referralCodePattern: /^REF[A-HJ-NP-Z2-9]{8}$/,
  
  // Debounce
  searchDebounceMs: 500,
  locationDebounceMs: 300,
  formSaveDebounceMs: 400,
};

/**
 * Create a security config with backend settings merged in.
 * 
 * @param backendSettings - Settings from /public/payment-settings endpoint
 * @returns Merged security configuration
 */
export function createSecurityConfig(backendSettings?: Record<string, string> | null): SecurityConfig {
  if (!backendSettings) {
    return DEFAULT_SECURITY_CONFIG;
  }

  return {
    ...DEFAULT_SECURITY_CONFIG,
    
    // Override OTP expiry from backend
    otpExpirySeconds: backendSettings.otp_expiry_minutes 
      ? parseInt(backendSettings.otp_expiry_minutes) * 60
      : DEFAULT_SECURITY_CONFIG.otpExpirySeconds,
    
    // Override max image size from backend
    maxImageSizeMB: backendSettings.max_image_size_mb
      ? parseInt(backendSettings.max_image_size_mb)
      : DEFAULT_SECURITY_CONFIG.maxImageSizeMB,
  };
}

/**
 * Validate security config values are within safe ranges.
 * Prevents misconfiguration from causing security issues.
 */
export function validateSecurityConfig(config: SecurityConfig): { valid: boolean; errors: string[] } {
  const errors: string[] = [];

  // JWT grace period should be reasonable (1-60 seconds)
  if (config.jwtGracePeriodMs < 1_000 || config.jwtGracePeriodMs > 60_000) {
    errors.push('JWT grace period must be between 1-60 seconds');
  }

  // Session check shouldn't be too frequent or too rare
  if (config.sessionCheckIntervalMs < 5_000 || config.sessionCheckIntervalMs > 300_000) {
    errors.push('Session check interval must be between 5s-5min');
  }

  // OTP length should be reasonable
  if (config.otpLength < 4 || config.otpLength > 12) {
    errors.push('OTP length must be between 4-12 digits');
  }

  // Max retries should be limited
  if (config.retryMaxAttempts < 0 || config.retryMaxAttempts > 5) {
    errors.push('Max retry attempts must be between 0-5');
  }

  // File size should be reasonable (1MB - 100MB)
  if (config.maxImageSizeMB < 1 || config.maxImageSizeMB > 100) {
    errors.push('Max image size must be between 1-100 MB');
  }

  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Log security config for debugging (redacts sensitive values).
 */
export function logSecurityConfig(config: SecurityConfig): void {
  if (process.env.NODE_ENV === 'development') {
    console.log('[SecurityConfig]', {
      jwtGracePeriod: `${config.jwtGracePeriodMs}ms`,
      sessionCheckInterval: `${config.sessionCheckIntervalMs}ms`,
      otpExpiry: `${config.otpExpirySeconds}s`,
      maxImageSize: `${config.maxImageSizeMB}MB`,
      retryAttempts: config.retryMaxAttempts,
    });
  }
}
