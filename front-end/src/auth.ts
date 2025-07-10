"use server";
import {
  getSession,
  createSession,
  destroySession,
  requireAuth,
} from "@/lib/auth/session";
import { redirect } from "next/navigation";

const API_URL = process.env.API_URL || "http://localhost:8080";

// Server-side auth function (mirrors Auth.js v5 approach)
export const auth = getSession;

// Server action for sending code
export async function sendCodeAction(email: string, name?: string) {
  const response = await fetch(`${API_URL}/api/v1/auth/send-code`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Client-Type": "nextjs",
    },
    body: JSON.stringify({ email, name }),
  });

  let data;
  try {
    data = await response.json();
  } catch (error) {
    console.error("Failed to parse JSON response:", error);
    throw new Error("Server returned invalid response");
  }

  if (!response.ok) {
    throw new Error(data.error || "Failed to send code");
  }

  return data;
}

// Server action for sign in
export async function signInAction(
  email: string,
  code: string,
  rememberMe: boolean = false
) {
  // Verify with Go backend
  const response = await fetch(`${API_URL}/api/v1/auth/verify-code`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Client-Type": "nextjs",
    },
    body: JSON.stringify({ email, code }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Authentication failed");
  }

  // Create session with user data
  const user = {
    id: data.user.id, // keep as number
    email: data.user.email,
    name: data.user.name,
    role: data.user.role,
    is_active: data.user.is_active,
    created_at: data.user.created_at,
  };

  await createSession(user, rememberMe);
  redirect("/");
}

// Server action for WebAuthn authentication
export async function webAuthnLoginAction(userId: number, assertion: Record<string, unknown>) {
  const response = await fetch(`${API_URL}/api/v1/webauthn/finish-login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Client-Type": "nextjs",
    },
    body: JSON.stringify({
      user_id: userId,
      assertion
    }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "WebAuthn authentication failed");
  }

  // Create session with user data
  const user = {
    id: data.user.id,
    email: data.user.email,
    name: data.user.name,
    role: data.user.role,
    is_active: data.user.is_active,
    created_at: data.user.created_at,
  };

  await createSession(user, true); // Always remember for WebAuthn
  redirect("/dashboard");
}

export async function signOutAction() {
  await destroySession();
  redirect("/");
}

// Server-side utilities
export { requireAuth, createSession, destroySession };

// Types
export type { User, SessionData } from "@/lib/auth/session";
