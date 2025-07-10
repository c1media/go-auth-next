import jwt from "jsonwebtoken";
import { cookies } from "next/headers";

export interface User {
  id: string;
  email: string;
  name?: string;
  role: string;
  is_active: boolean;
  created_at: string;
}

export interface SessionData {
  user: User;
  exp: number;
  iat: number;
}

const JWT_SECRET =
  process.env.JWT_SECRET || "your-secret-key-change-in-production";

export async function createSession(user: User, rememberMe: boolean = false): Promise<void> {
  const expirationTime = rememberMe ? 30 * 24 * 60 * 60 : 24 * 60 * 60; // 30 days or 24 hours
  const payload = {
    user,
    exp: Math.floor(Date.now() / 1000) + expirationTime,
    iat: Math.floor(Date.now() / 1000),
  };

  const token = jwt.sign(payload, JWT_SECRET);

  const cookieStore = await cookies();
  cookieStore.set("session", token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    maxAge: expirationTime * 1000, // Convert to milliseconds
    path: "/",
  });
}

export async function getSession(): Promise<SessionData | null> {
  const cookieStore = await cookies();
  const token = cookieStore.get("session")?.value;

  if (!token) return null;

  try {
    const decoded = jwt.verify(token, JWT_SECRET) as SessionData;

    // Check if token is expired
    if (decoded.exp < Math.floor(Date.now() / 1000)) {
      await destroySession();
      return null;
    }

    return decoded;
  } catch {
    await destroySession();
    return null;
  }
}

export async function destroySession(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete("session");
}

export async function requireAuth(): Promise<SessionData> {
  const session = await getSession();
  if (!session) {
    throw new Error("Authentication required");
  }
  return session;
}
