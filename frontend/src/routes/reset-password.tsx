import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import ApiService from "@/services/api-service";
import { Dumbbell } from "lucide-react";
import { useState } from "react";

const PasswordReset = () => {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState<string | null>(null);

  const submitPasswordReset = async () => {
    if (email === undefined || email === null || email.length === 0) {
      return;
    }
    const res = await ApiService.passwordReset(email);
    if (res.status === 204) {
      setMessage("Check your email for a link to reset your password");
      return;
    }

    setMessage("There was an error sending the password reset email. Please try again.");
  };

  if (message !== null) {
    return (
      <>
        <PasswordResetContent>
          <p>{message}</p>
        </PasswordResetContent>
      </>
    );
  }

  return (
    <>
      <PasswordResetContent>
        <><Input value={email} onChange={e => setEmail(e.target.value)} type="email" placeholder="Email" /> <Button onClick={submitPasswordReset}>Reset Password</Button></>
      </PasswordResetContent>
    </>
  );
};

export const PasswordResetContent = ({ ...props }: React.ComponentProps<"div">) => {
  return (
    <div className="flex min-h-svh flex-col items-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="/" className="flex items-center gap-2 self-center font-medium text-2xl">
          <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
            <Dumbbell className="size-6" />
          </div>
          Gymotric
        </a>
      </div>
      <div className="flex w-full max-w-sm flex-col gap-6" {...props} />
    </div>
  );
}

export default PasswordReset;
