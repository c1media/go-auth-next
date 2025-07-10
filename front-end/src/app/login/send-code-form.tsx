import { handleSendCodeAction } from "./actions";
import { Button } from "@/components/ui/button"

interface SendCodeFormProps {
  email: string;
}

export function SendCodeForm({ email }: SendCodeFormProps) {
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await handleSendCodeAction(email);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input type="hidden" name="email" value={email} />
      <Button type="submit" className="w-full">
        Send Code
      </Button>
    </form>
  )
}