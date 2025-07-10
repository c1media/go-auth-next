"use server";

import { sendCodeAction } from "@/auth";
import { redirect } from "next/navigation";

export async function handleSendCodeAction(email: string) {
  await sendCodeAction(email);
  // Redirect to code verification step
  redirect(`/login/verify?email=${encodeURIComponent(email)}`);
}
