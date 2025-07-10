"use client";
import { useState, useEffect } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { RegisterPasskeyButton } from "@/components/auth/register-passkey-button";

interface PasskeySetupAlertProps {
     email: string;
}

interface UserCheckResponse {
     user_exists: boolean;
     has_passkeys: boolean;
     user_id?: number;
}

export function PasskeySetupAlert({ email }: PasskeySetupAlertProps) {
     const [hasPasskeys, setHasPasskeys] = useState<boolean | null>(null);
     const [userId, setUserId] = useState<number | null>(null);
     const [loading, setLoading] = useState(true);
     const [showSetup, setShowSetup] = useState(false);

     useEffect(() => {
          const checkPasskeys = async () => {
               try {
                    const response = await fetch("/api/auth/check-user", {
                         method: "POST",
                         headers: { "Content-Type": "application/json" },
                         body: JSON.stringify({ email }),
                    });

                    if (response.ok) {
                         const data: UserCheckResponse = await response.json();
                         setHasPasskeys(data.has_passkeys);
                         setUserId(data.user_id || null);
                    }
               } catch (error) {
                    console.error("Failed to check passkeys:", error);
               } finally {
                    setLoading(false);
               }
          };

          checkPasskeys();
     }, [email]);

     if (loading) {
          return null; // Don't show anything while loading
     }

     if (hasPasskeys) {
          return null; // Don't show alert if user has passkeys
     }

     if (showSetup) {
          return (
               <Alert className="border-blue-200 bg-blue-50">
                    <AlertDescription className="space-y-4">
                         <div>
                              <p className="font-medium text-blue-800 mb-2">
                                   Set up a passkey for faster, more secure sign-ins
                              </p>
                              <p className="text-blue-700 text-sm">
                                   Passkeys let you sign in with your fingerprint, face, or device PIN instead of typing a code.
                              </p>
                         </div>
                         {userId && <RegisterPasskeyButton userId={userId} />}
                         <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => setShowSetup(false)}
                              className="text-blue-600 hover:text-blue-800"
                         >
                              Maybe later
                         </Button>
                    </AlertDescription>
               </Alert>
          );
     }

     return (
          <Alert className="border-amber-200 bg-amber-50">
               <AlertDescription className="flex items-center justify-between">
                    <div>
                         <p className="font-medium text-amber-800">
                              🔐 Set up a passkey for faster sign-ins
                         </p>
                         <p className="text-amber-700 text-sm mt-1">
                              Skip typing codes and sign in with your fingerprint or face
                         </p>
                    </div>
                    <Button
                         variant="outline"
                         size="sm"
                         onClick={() => setShowSetup(true)}
                         className="ml-4"
                    >
                         Set up passkey
                    </Button>
               </AlertDescription>
          </Alert>
     );
} 