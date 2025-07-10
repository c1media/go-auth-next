"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { SignInWithPasskeyButton } from "@/components/auth/sign-in-with-passkey-button";

type LoginState = "email" | "options" | "passkey";

interface UserCheckResponse {
     user_exists: boolean;
     has_passkeys: boolean;
     user_id?: number;
}

export function EmailLoginForm() {
     const [email, setEmail] = useState("");
     const [state, setState] = useState<LoginState>("email");
     const [userData, setUserData] = useState<UserCheckResponse | null>(null);
     const [error, setError] = useState<string>("");
     const [loading, setLoading] = useState(false);

     const handleEmailSubmit = async (e: React.FormEvent) => {
          e.preventDefault();
          setLoading(true);
          setError("");

          try {
               const response = await fetch("/api/auth/check-user", {
                    method: "POST",
                    headers: {
                         "Content-Type": "application/json",
                         "X-Client-Type": "nextjs"
                    },
                    body: JSON.stringify({ email }),
               });

               if (!response.ok) {
                    console.error("Failed to check user:", response.status, response.statusText);
                    throw new Error("Failed to check user");
               }

               const data: UserCheckResponse = await response.json();
               setUserData(data);

               if (!data.user_exists) {
                    // User doesn't exist, show sign up option
                    setError("User not found. Please sign up first.");
                    return;
               }

               // Always show options - let user choose their preferred method
               setState("options");
          } catch (error) {
               console.error("Error in handleEmailSubmit:", error);
               setError("Failed to check user. Please try again.");
          } finally {
               setLoading(false);
          }
     };

     const handlePasskeyLogin = () => {
          setState("passkey");
          // Trigger the passkey authentication immediately
          setTimeout(() => {
               const passkeyButton = document.querySelector('[data-passkey-button]') as HTMLButtonElement;
               if (passkeyButton) {
                    passkeyButton.click();
               }
          }, 100);
     };

     const handleCodeLogin = () => {
          // Send code and redirect to verification
          const sendCodeAndRedirect = async () => {
               try {
                    const sendCodeResponse = await fetch("/api/auth/send-code", {
                         method: "POST",
                         headers: {
                              "Content-Type": "application/json",
                              "X-Client-Type": "nextjs"
                         },
                         body: JSON.stringify({ email }),
                    });

                    if (sendCodeResponse.ok) {
                         window.location.href = `/login/verify?email=${encodeURIComponent(email)}`;
                    } else {
                         setError("Failed to send code. Please try again.");
                    }
               } catch {
                    setError("Failed to send code. Please try again.");
               }
          };
          sendCodeAndRedirect();
     };

     const handleBack = () => {
          setState("email");
          setError("");
          setUserData(null);
     };

     const handleBackToOptions = () => {
          setState("options");
     };

     if (state === "email") {
          return (
               <form onSubmit={handleEmailSubmit} className="space-y-4">
                    <div className="space-y-2">
                         <Label htmlFor="email">Email</Label>
                         <Input
                              id="email"
                              type="email"
                              placeholder="Enter your email"
                              value={email}
                              onChange={(e) => setEmail(e.target.value)}
                              required
                              disabled={loading}
                         />
                    </div>
                    <Button type="submit" className="w-full" disabled={loading}>
                         {loading ? "Checking..." : "Continue"}
                    </Button>
                    {error && (
                         <Alert variant="destructive">
                              <AlertDescription>{error}</AlertDescription>
                         </Alert>
                    )}
               </form>
          );
     }

     if (state === "options") {
          return (
               <div className="space-y-4">
                    <div className="text-center">
                         <p className="text-sm text-gray-600 mb-4">
                              Sign in as <strong>{email}</strong>
                         </p>
                    </div>
                    {userData?.has_passkeys ? (
                         <>
                              <Button
                                   onClick={handlePasskeyLogin}
                                   className="w-full"
                                   variant="default"
                              >
                                   🔐 Sign in with Passkey
                              </Button>
                              <Button
                                   onClick={handleCodeLogin}
                                   className="w-full"
                                   variant="outline"
                              >
                                   📧 Sign in with Code
                              </Button>
                         </>
                    ) : (
                         <Button
                              onClick={handleCodeLogin}
                              className="w-full"
                              variant="default"
                         >
                              📧 Send Verification Code
                         </Button>
                    )}
                    <Button
                         onClick={handleBack}
                         className="w-full"
                         variant="ghost"
                    >
                         Use different email
                    </Button>
               </div>
          );
     }

     if (state === "passkey") {
          return (
               <div className="space-y-4">
                    <div className="text-center">
                         <p className="text-sm text-gray-600 mb-4">
                              Sign in with passkey for <strong>{email}</strong>
                              <br />
                              <span className="text-xs text-gray-500">Click the button below to use your passkey</span>
                         </p>
                    </div>
                    <SignInWithPasskeyButton email={email} userId={userData?.user_id} />
                    <Button onClick={handleCodeLogin} className="w-full" variant="outline">
                         Try code instead
                    </Button>
                    <Button onClick={handleBackToOptions} className="w-full" variant="ghost">
                         Use different email
                    </Button>
               </div>
          );
     }


     return null;
} 