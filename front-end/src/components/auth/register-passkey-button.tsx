"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { createPasskey } from "@/auth/webauthn-browser";

function base64urlToBuffer(base64url: string): ArrayBuffer {
     let base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
     while (base64.length % 4) base64 += '=';
     const str = atob(base64);
     const bytes = new Uint8Array(str.length);
     for (let i = 0; i < str.length; ++i) bytes[i] = str.charCodeAt(i);
     return bytes.buffer;
}

export function RegisterPasskeyButton({ userId }: { userId: number }) {
     const [status, setStatus] = useState<string>("");

     async function handleRegister() {
          setStatus("Requesting registration options...");
          const response = await fetch("/api/webauthn/begin-registration", {
               method: "POST",
               headers: { "Content-Type": "application/json" },
               body: JSON.stringify({ user_id: userId }),
          });
          if (!response.ok) {
               setStatus("Failed to get registration options");
               return;
          }
          const options = await response.json();

          setStatus("Triggering browser passkey dialog...");
          let credential: Record<string, unknown>;
          // Convert challenge and user.id to ArrayBuffer as required by WebAuthn
          options.challenge = base64urlToBuffer(options.challenge);
          options.user.id = base64urlToBuffer(options.user.id);
          try {
               credential = await createPasskey(options);
          } catch (err) {
               setStatus("Registration cancelled or failed");
               console.error("WebAuthn registration error:", err);
               let msg = "";
               if (err && typeof err === "object" && "message" in err) {
                    msg = (err as { message: string }).message;
               } else {
                    msg = String(err);
               }
               alert("WebAuthn registration error: " + msg);
               return;
          }

          setStatus("Sending credential to server...");
          const finishResp = await fetch("/api/webauthn/finish-registration", {
               method: "POST",
               headers: { "Content-Type": "application/json" },
               body: JSON.stringify({ credential, user_id: userId }),
          });
          if (!finishResp.ok) {
               setStatus("Failed to register passkey");
               return;
          }
          setStatus("Passkey registered successfully!");
     }

     return (
          <div className="mt-4">
               <Button type="button" variant="default" onClick={handleRegister}>
                    Register Passkey
               </Button>
               {status && <div className="mt-2 text-sm text-gray-600">{status}</div>}
          </div>
     );
} 