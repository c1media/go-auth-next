import { auth } from "@/auth"
import { redirect } from "next/navigation"
import { SessionDebug } from "@/components/auth/session-debug"
import { SignOut } from "@/components/auth/sign-out"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RegisterPasskeyButton } from "@/components/auth/register-passkey-button";
import WebAuthnCredentialsManager from "@/components/auth/webauthn-credentials-manager"
import { PasskeySetupAlert } from "@/components/auth/passkey-setup-alert";


export default async function Dashboard() {
  const session = await auth()

  // Redirect to login if not authenticated
  if (!session) {
    redirect("/login")
  }

  // Session data is available for debugging if needed

  const user = session.user

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Header */}
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl">Dashboard</CardTitle>
              <CardDescription>Welcome to your dashboard</CardDescription>
            </CardHeader>
          </Card>

          {/* Passkey Setup Alert */}
          <PasskeySetupAlert email={user.email} />

          {/* User Info */}
          <Card>
            <CardHeader>
              <CardTitle>User Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-600">Email</p>
                  <p className="font-medium">{user.email}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Name</p>
                  <p className="font-medium">{user.name || "Not set"}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Role</p>
                  <Badge variant={user.role === "admin" ? "default" : user.role === "moderator" ? "secondary" : "outline"}>
                    {user.role}
                  </Badge>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Status</p>
                  <Badge variant={user.is_active ? "default" : "destructive"}>
                    {user.is_active ? "Active" : "Inactive"}
                  </Badge>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Created</p>
                  <p className="font-medium">{new Date(user.created_at).toLocaleDateString()}</p>
                </div>
              </div>
              <RegisterPasskeyButton userId={Number(user.id)} />
              <WebAuthnCredentialsManager userId={Number(user.id)} />
            </CardContent>
          </Card>

          {/* Role-based Content */}
          <Card>
            <CardHeader>
              <CardTitle>Role-based Features</CardTitle>
            </CardHeader>
            <CardContent>
              {user.role === "admin" && (
                <div className="border-l-4 border-red-500 bg-red-50 p-4 rounded-r-lg">
                  <h3 className="text-lg font-semibold text-red-800">Admin Panel</h3>
                  <p className="text-red-700 mb-4">🔧 You have admin access - you can manage users and system settings</p>
                  <div className="space-x-2">
                    <Button variant="outline" size="sm">Manage Users</Button>
                    <Button variant="outline" size="sm">System Settings</Button>
                  </div>
                </div>
              )}

              {user.role === "moderator" && (
                <div className="border-l-4 border-blue-500 bg-blue-50 p-4 rounded-r-lg">
                  <h3 className="text-lg font-semibold text-blue-800">Moderator Panel</h3>
                  <p className="text-blue-700 mb-4">📝 You have moderator access - you can moderate content</p>
                  <div className="space-x-2">
                    <Button variant="outline" size="sm">Moderate Content</Button>
                    <Button variant="outline" size="sm">View Reports</Button>
                  </div>
                </div>
              )}

              {user.role === "user" && (
                <div className="border-l-4 border-green-500 bg-green-50 p-4 rounded-r-lg">
                  <h3 className="text-lg font-semibold text-green-800">User Panel</h3>
                  <p className="text-green-700 mb-4">👤 You have user access - basic features available</p>
                  <div className="space-x-2">
                    <Button variant="outline" size="sm">View Profile</Button>
                    <Button variant="outline" size="sm">Edit Settings</Button>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Actions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col sm:flex-row gap-2">
                <Link href="/">
                  <Button variant="outline">Home</Button>
                </Link>
                <SignOut />
              </div>
            </CardContent>
          </Card>

          <SessionDebug />
        </div>
      </div>
    </div>
  )
}