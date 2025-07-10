"use server";
import { bufferToBase64Url } from "./webauthn-browser";

export async function listWebAuthnCredentials(userId: number) {
  const apiUrl = process.env.API_URL || "http://localhost:8080";
  const res = await fetch(
    `${apiUrl}/api/v1/webauthn/list-credentials`,
    {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
        "X-Client-Type": "nextjs"
      },
      body: JSON.stringify({ user_id: userId }),
      credentials: "include",
    }
  );
  if (!res.ok) {
    throw new Error("Failed to fetch credentials");
  }
  const data = await res.json();
  return data.credentials;
}

export async function deleteWebAuthnCredential(
  userId: number,
  credentialId: string | ArrayBuffer
) {
  // If credentialId is ArrayBuffer, convert to base64url
  let credIdStr: string;
  if (credentialId instanceof ArrayBuffer) {
    credIdStr = bufferToBase64Url(credentialId);
  } else {
    credIdStr = credentialId;
  }
  const apiUrl = process.env.API_URL || "http://localhost:8080";
  const res = await fetch(
    `${apiUrl}/api/v1/webauthn/delete-credential`,
    {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
        "X-Client-Type": "nextjs"
      },
      body: JSON.stringify({ user_id: userId, credential_id: credIdStr }),
      credentials: "include",
    }
  );
  if (!res.ok) {
    throw new Error("Failed to delete credential");
  }
  return await res.json();
}
