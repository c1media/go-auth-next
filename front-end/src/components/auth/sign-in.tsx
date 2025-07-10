import { redirect } from "next/navigation"
import { Button } from "@/components/ui/button"

export function SignIn() {
     return (
          <form
               action={async () => {
                    "use server"
                    redirect("/login")
               }}
          >
               <Button type="submit">Sign In</Button>
          </form>
     )
}