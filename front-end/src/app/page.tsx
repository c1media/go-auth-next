import { SignIn } from "@/components/auth/sign-in";
import { SignOut } from "@/components/auth/sign-out";
import { SessionDebug } from "@/components/auth/session-debug";
import { auth } from "@/auth";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default async function Home() {
  const session = await auth();

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-3xl">Auth Template</CardTitle>
              <CardDescription>
                Next.js 15 + Custom Auth + Go Backend
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {session ? (
                <div className="space-y-4">
                  <div className="text-center">
                    <p className="text-lg">Welcome back, <span className="font-semibold">{session.user.email}</span>!</p>
                    <p className="text-sm text-gray-600">Role: {session.user.role}</p>
                  </div>
                  <div className="flex flex-col sm:flex-row gap-2 justify-center">
                    <Link href="/dashboard">
                      <Button className="w-full sm:w-auto">Go to Dashboard</Button>
                    </Link>
                    <SignOut />
                  </div>
                </div>
              ) : (
                <div className="space-y-4 text-center">
                  <p className="text-lg">Please sign in to access the dashboard</p>
                  <SignIn />
                </div>
              )}

              <SessionDebug />
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  );
}
