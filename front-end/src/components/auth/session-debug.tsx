import { auth } from "@/auth"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

export async function SessionDebug() {
  const session = await auth()

  return (
    <Card className="mt-6">
      <CardHeader>
        <CardTitle className="text-sm">Session Debug</CardTitle>
        <CardDescription>Current session information</CardDescription>
      </CardHeader>
      <CardContent>
        {session ? (
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <Badge variant="default">Authenticated</Badge>
              <span className="text-sm text-muted-foreground">Server-side session</span>
            </div>
            <details className="mt-4">
              <summary className="cursor-pointer text-sm font-medium hover:text-primary">
                View session data
              </summary>
              <pre className="mt-2 text-xs bg-muted p-2 rounded overflow-x-auto">
                {JSON.stringify(session, null, 2)}
              </pre>
            </details>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <Badge variant="outline">Not authenticated</Badge>
            <span className="text-sm text-muted-foreground">No session found</span>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
