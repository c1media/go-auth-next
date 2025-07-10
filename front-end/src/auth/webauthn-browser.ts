// Utility for browser-side WebAuthn logic

export async function createPasskey(
  options: PublicKeyCredentialCreationOptions
): Promise<Record<string, unknown>> {
  if (!window.PublicKeyCredential) {
    throw new Error("WebAuthn is not supported in this browser.");
  }
  const credential = (await navigator.credentials.create({
    publicKey: options,
  })) as PublicKeyCredential;
  if (!credential) throw new Error("Failed to create credential");
  // Convert to JSON for sending to backend
  const converted = credentialToJSON(credential);
  return converted;
}

export async function getPasskey(
  options: PublicKeyCredentialRequestOptions
): Promise<Record<string, unknown>> {
  if (!window.PublicKeyCredential) {
    throw new Error("WebAuthn is not supported in this browser.");
  }
  const assertion = (await navigator.credentials.get({
    publicKey: options,
  })) as PublicKeyCredential | null;
  if (!assertion) throw new Error("Failed to get credential");
  return assertionToJSON(assertion);
}

// Helper to convert credential to JSON for registration
function credentialToJSON(cred: PublicKeyCredential): Record<string, unknown> {
  const result: Record<string, unknown> = {
    id: cred.id,
    type: cred.type,
  };

  // Handle rawId (ArrayBuffer)
  if (cred.rawId) {
    result.rawId = bufferToBase64Url(cred.rawId);
  }

  // Handle response based on type
  if (cred.response) {
    if (cred.response instanceof AuthenticatorAttestationResponse) {
      // Registration response
      result.response = {
        attestationObject: bufferToBase64Url(cred.response.attestationObject),
        clientDataJSON: bufferToBase64Url(cred.response.clientDataJSON),
      };
    } else if (cred.response instanceof AuthenticatorAssertionResponse) {
      // Login response
      result.response = {
        authenticatorData: bufferToBase64Url(cred.response.authenticatorData),
        clientDataJSON: bufferToBase64Url(cred.response.clientDataJSON),
        signature: bufferToBase64Url(cred.response.signature),
        userHandle: cred.response.userHandle
          ? bufferToBase64Url(cred.response.userHandle)
          : null,
      };
    }
  }

  // Handle authenticatorAttachment if present
  if (cred.authenticatorAttachment) {
    result.authenticatorAttachment = cred.authenticatorAttachment;
  }

  return result;
}

// Helper to convert credential to JSON for login
function assertionToJSON(cred: PublicKeyCredential): Record<string, unknown> {
  const result: Record<string, unknown> = {
    id: cred.id,
    type: cred.type,
  };

  // Handle rawId (ArrayBuffer)
  if (cred.rawId) {
    result.rawId = bufferToBase64Url(cred.rawId);
  }

  // Handle response
  if (
    cred.response &&
    cred.response instanceof AuthenticatorAssertionResponse
  ) {
    result.response = {
      authenticatorData: bufferToBase64Url(cred.response.authenticatorData),
      clientDataJSON: bufferToBase64Url(cred.response.clientDataJSON),
      signature: bufferToBase64Url(cred.response.signature),
      userHandle: cred.response.userHandle
        ? bufferToBase64Url(cred.response.userHandle)
        : null,
    };
  }

  // Handle authenticatorAttachment if present
  if (cred.authenticatorAttachment) {
    result.authenticatorAttachment = cred.authenticatorAttachment;
  }

  return result;
}

export function bufferToBase64Url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let str = "";
  for (const byte of bytes) {
    str += String.fromCharCode(byte);
  }
  return btoa(str).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}
