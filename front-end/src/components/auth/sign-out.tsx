import { signOutAction } from "@/auth"
import { Button } from "@/components/ui/button"

export function SignOut() {
  return (
    <form
      action={async () => {
        "use server"
        await signOutAction()
      }}
    >
      <Button type="submit" variant="outline">Sign Out</Button>
    </form>
  )
}