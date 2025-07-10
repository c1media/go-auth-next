import { signInAction } from "@/auth"
import { redirect } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Checkbox } from "@/components/ui/checkbox"

interface VerifyPageProps {
  searchParams: Promise<{ email?: string; error?: string }>
}

export default async function VerifyPage({ searchParams }: VerifyPageProps) {
  const { email, error } = await searchParams

  if (!email) {
    redirect("/login")
  }

  async function handleVerifyCode(formData: FormData) {
    "use server"
    const code = formData.get("code") as string
    const rememberMe = formData.get("rememberMe") === "on"

    await signInAction(email!, code, rememberMe)
    // Redirect happens automatically in signInAction
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl text-center">Enter Verification Code</CardTitle>
          <CardDescription className="text-center">
            Code sent to {email}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Alert className="mb-4 border-blue-200 bg-blue-50">
            <AlertDescription className="text-blue-700 text-sm">
              💡 After signing in, you can set up a passkey for faster, more secure sign-ins in your dashboard.
            </AlertDescription>
          </Alert>
          <form action={handleVerifyCode} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="code">Verification Code</Label>
              <Input
                id="code"
                name="code"
                type="text"
                placeholder="Enter verification code"
                required
              />
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="rememberMe" name="rememberMe" />
              <Label htmlFor="rememberMe" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Remember me for 30 days
              </Label>
            </div>
            <Button type="submit" className="w-full">
              Sign In
            </Button>
            <Link href="/login">
              <Button type="button" variant="outline" className="w-full">
                Back
              </Button>
            </Link>
          </form>

          {error && (
            <Alert variant="destructive" className="mt-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>
    </div>
  )
}