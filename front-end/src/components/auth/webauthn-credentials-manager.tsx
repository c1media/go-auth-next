"use client";
import { useEffect, useState } from "react";
import { listWebAuthnCredentials, deleteWebAuthnCredential } from "@/auth/webauthn-actions";

type WebAuthnCredential = {
     credential_id: string;
     name?: string;
     public_key?: string;
     counter?: number;
     user_id?: number;
};

export default function WebAuthnCredentialsManager({ userId }: { userId: number }) {
     const [credentials, setCredentials] = useState<WebAuthnCredential[]>([]);
     const [loading, setLoading] = useState(true);
     const [error, setError] = useState<string | null>(null);
     const [deleting, setDeleting] = useState<string | null>(null);

     useEffect(() => {
          setLoading(true);
          listWebAuthnCredentials(userId)
               .then(setCredentials)
               .catch((e) => setError(e.message))
               .finally(() => setLoading(false));
     }, [userId]);

     const handleDelete = async (credentialId: string) => {
          setDeleting(credentialId);
          try {
               await deleteWebAuthnCredential(userId, credentialId);
               setCredentials((creds) => creds.filter((c) => c.credential_id !== credentialId));
          } catch (e: unknown) {
               if (e instanceof Error) {
                    setError(e.message);
               } else {
                    setError("Unknown error");
               }
          } finally {
               setDeleting(null);
          }
     };

     if (loading) return <div>Loading credentials...</div>;
     if (error) return <div className="text-red-500">Error: {error}</div>;

     return (
          <div className="space-y-2">
               <h3 className="font-semibold">Passkeys / WebAuthn Credentials</h3>
               {credentials.length === 0 ? (
                    <div>No passkeys registered.</div>
               ) : (
                    <ul className="space-y-2">
                         {credentials.map((cred) => (
                              <li key={cred.credential_id} className="flex items-center justify-between border rounded p-2">
                                   <span>{cred.name || "Unnamed Device"} <span className="text-xs text-gray-500">({cred.credential_id.slice(0, 8)}...)</span></span>
                                   <button
                                        className="text-red-600 hover:underline disabled:opacity-50"
                                        onClick={() => handleDelete(cred.credential_id)}
                                        disabled={deleting === cred.credential_id}
                                   >
                                        {deleting === cred.credential_id ? "Deleting..." : "Delete"}
                                   </button>
                              </li>
                         ))}
                    </ul>
               )}
          </div>
     );
} 