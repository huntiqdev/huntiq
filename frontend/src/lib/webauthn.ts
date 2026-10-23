/**
 * Enterprise WebAuthn / Passkey Client Utilities for HUNTIQ (SellBuy)
 *
 * Implements FIDO2 / WebAuthn browser ceremonies supporting compatible platform
 * authenticators (Touch ID, Face ID, Windows Hello, Android Biometrics) and
 * external hardware security keys (YubiKey, Titan).
 *
 * Biometric Privacy Guarantee:
 * HUNTIQ never collects or stores raw biometrics (fingerprints or face scans).
 * Biometric verification occurs strictly on the local hardware device.
 * Only standard public-key cryptographic assertions are exchanged with the server.
 */

// Base64URL decoding to Uint8Array/ArrayBuffer
export function base64UrlToBuffer(base64url: string): Uint8Array {
  // Pad if necessary
  let base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
  while (base64.length % 4 !== 0) {
    base64 += '=';
  }
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

// ArrayBuffer/Uint8Array encoding to Base64URL string
export function bufferToBase64Url(buffer: ArrayBuffer | Uint8Array): string {
  const bytes = buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

/**
 * Check if the browser environment supports WebAuthn APIs.
 */
export function isWebAuthnSupported(): boolean {
  return typeof window !== 'undefined' &&
    window.PublicKeyCredential !== undefined &&
    typeof window.PublicKeyCredential === 'function';
}

/**
 * Check if a platform authenticator (TouchID, Windows Hello, Android Biometric) is available.
 */
export async function isPlatformAuthenticatorAvailable(): Promise<boolean> {
  if (!isWebAuthnSupported()) return false;
  try {
    if (typeof PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable === 'function') {
      return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
    }
  } catch {
    return false;
  }
  return false;
}

/**
 * Convert server-provided CredentialCreationOptions JSON into PublicKeyCredentialCreationOptions
 * and execute navigator.credentials.create.
 */
export async function performPasskeyRegistration(optionsJSON: any): Promise<any> {
  if (!isWebAuthnSupported()) {
    throw new Error('WebAuthn / Passkey authentication is not supported on this browser or device.');
  }

  const publicKey = optionsJSON.publicKey || optionsJSON;

  // Transform challenge
  const challenge = base64UrlToBuffer(publicKey.challenge);

  // Transform user.id
  let userId: Uint8Array;
  if (typeof publicKey.user.id === 'string') {
    userId = base64UrlToBuffer(publicKey.user.id);
  } else {
    userId = new Uint8Array(publicKey.user.id);
  }

  // Transform excludeCredentials if present
  const excludeCredentials = (publicKey.excludeCredentials || []).map((cred: any) => ({
    ...cred,
    id: typeof cred.id === 'string' ? base64UrlToBuffer(cred.id) : new Uint8Array(cred.id),
  }));

  const creationOptions: CredentialCreationOptions = {
    publicKey: {
      ...publicKey,
      challenge: challenge.buffer as ArrayBuffer,
      user: {
        ...publicKey.user,
        id: userId.buffer as ArrayBuffer,
      },
      excludeCredentials: excludeCredentials.length > 0 ? excludeCredentials : undefined,
    },
  };

  try {
    const credential = (await navigator.credentials.create(creationOptions)) as PublicKeyCredential | null;
    if (!credential) {
      throw new Error('Authenticator did not return a credential.');
    }

    const response = credential.response as AuthenticatorAttestationResponse;

    return {
      id: credential.id,
      rawId: bufferToBase64Url(credential.rawId),
      type: credential.type,
      response: {
        clientDataJSON: bufferToBase64Url(response.clientDataJSON),
        attestationObject: bufferToBase64Url(response.attestationObject),
        transports: typeof (response as any).getTransports === 'function' ? (response as any).getTransports() : [],
      },
    };
  } catch (err: any) {
    if (err.name === 'NotAllowedError') {
      throw new Error('Registration was cancelled or timed out. Please try again.');
    }
    if (err.name === 'InvalidStateError') {
      throw new Error('This authenticator is already registered for this account.');
    }
    if (err.name === 'NotSupportedError') {
      throw new Error('The requested passkey algorithm or security level is not supported by your device.');
    }
    throw new Error(err.message || 'Passkey enrollment failed.');
  }
}

/**
 * Convert server-provided CredentialRequestOptions JSON into PublicKeyCredentialRequestOptions
 * and execute navigator.credentials.get.
 */
export async function performPasskeyAuthentication(optionsJSON: any): Promise<any> {
  if (!isWebAuthnSupported()) {
    throw new Error('WebAuthn / Passkey authentication is not supported on this browser or device.');
  }

  const publicKey = optionsJSON.publicKey || optionsJSON;

  // Transform challenge
  const challenge = base64UrlToBuffer(publicKey.challenge);

  // Transform allowCredentials if present
  const allowCredentials = (publicKey.allowCredentials || []).map((cred: any) => ({
    ...cred,
    id: typeof cred.id === 'string' ? base64UrlToBuffer(cred.id) : new Uint8Array(cred.id),
  }));

  const requestOptions: CredentialRequestOptions = {
    publicKey: {
      ...publicKey,
      challenge: challenge.buffer as ArrayBuffer,
      allowCredentials: allowCredentials.length > 0 ? allowCredentials : undefined,
    },
  };

  try {
    const assertion = (await navigator.credentials.get(requestOptions)) as PublicKeyCredential | null;
    if (!assertion) {
      throw new Error('Authenticator did not return an assertion.');
    }

    const response = assertion.response as AuthenticatorAssertionResponse;

    let userHandle: string | null = null;
    if (response.userHandle && response.userHandle.byteLength > 0) {
      userHandle = bufferToBase64Url(response.userHandle);
    }

    return {
      id: assertion.id,
      rawId: bufferToBase64Url(assertion.rawId),
      type: assertion.type,
      response: {
        clientDataJSON: bufferToBase64Url(response.clientDataJSON),
        authenticatorData: bufferToBase64Url(response.authenticatorData),
        signature: bufferToBase64Url(response.signature),
        userHandle: userHandle,
      },
    };
  } catch (err: any) {
    if (err.name === 'NotAllowedError') {
      throw new Error('Authentication cancelled, timed out, or biometric verification failed.');
    }
    if (err.name === 'SecurityError') {
      throw new Error('Security origin mismatch. Ensure HTTPS or localhost is configured.');
    }
    throw new Error(err.message || 'Passkey authentication failed.');
  }
}
