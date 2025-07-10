"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { getPasskey } from "@/auth/webauthn-browser";
import { webAuthnLoginAction } from "@/auth";

// Helper function to convert base64url string to ArrayBuffer
function base64UrlToBuffer(base64url: string): ArrayBuffer {
  // Add padding if needed
  const padding = '='.repeat((4 - (base64url.length % 4)) % 4);
  const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/') + padding;
  
  // Decode base64 to binary string
  const binaryString = atob(base64);
  
  // Convert binary string to ArrayBuffer
  const buffer = new ArrayBuffer(binaryString.length);
  const view = new Uint8Array(buffer);
  for (let i = 0; i < binaryString.length; i++) {
    view[i] = binaryString.charCodeAt(i);
  }
  
  return buffer;
}

interface SignInWithPasskeyButtonProps {
     email: string;
     userId?: number;
}

export function SignInWithPasskeyButton({ email, userId }: SignInWithPasskeyButtonProps) {
     const [status, setStatus] = useState<string>("");
     const [isLoading, setIsLoading] = useState(false);

     async function handleSignIn() {
          if (isLoading) return; // Prevent double-clicks
          setIsLoading(true);
          setStatus("Requesting login options...");
          const response = await fetch("/api/webauthn/begin-login", {
               method: "POST",
               headers: { 
                    "Content-Type": "application/json",
                    "X-Client-Type": "nextjs"
               },
               body: JSON.stringify({ email }),
          });
          if (!response.ok) {
               setStatus("Failed to get login options");
               setIsLoading(false);
               return;
          }
          const options = await response.json();

          setStatus("Triggering browser passkey dialog...");
          let assertion;
          try {
               // Convert base64url strings to ArrayBuffers for WebAuthn API
               const webAuthnOptions = {
                    ...options.publicKey,
                    challenge: base64UrlToBuffer(options.publicKey.challenge),
                    allowCredentials: options.publicKey.allowCredentials?.map((cred: { id: string; type: string }) => ({
                         ...cred,
                         id: base64UrlToBuffer(cred.id as string)
                    }))
               };
               assertion = await getPasskey(webAuthnOptions);
          } catch (error) {
               console.error("WebAuthn error:", error);
               setStatus("Login cancelled or failed: " + (error instanceof Error ? error.message : String(error)));
               setIsLoading(false);
               return;
          }

          setStatus("Authenticating with server...");
          try {
               if (!userId) {
                    throw new Error("User ID is required");
               }
               await webAuthnLoginAction(userId, assertion);
               setStatus("Signed in with passkey successfully!");
               // The server action will handle the redirect
          } catch (error) {
               console.error("WebAuthn login error:", error);
               setStatus("Failed to sign in with passkey: " + (error instanceof Error ? error.message : String(error)));
               setIsLoading(false);
          }
     }

     return (
          <div className="mb-4">
               <Button 
                    type="button" 
                    variant="default" 
                    className="w-full" 
                    onClick={handleSignIn}
                    disabled={isLoading}
                    data-passkey-button
               >
                    {isLoading ? "Authenticating..." : "Sign in with Passkey"}
               </Button>
               {status && <div className="mt-2 text-sm text-gray-600">{status}</div>}
          </div>
     );
} 