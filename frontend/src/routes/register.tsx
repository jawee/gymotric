import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import ApiService from "@/services/api-service";
import { Dumbbell } from "lucide-react";
import { useState } from "react";

const Register = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  const submitRegistration = async () => {
    if (username === "") {
      setValidationMessage("Username is required");
      return;
    }
    if (email === "") {
      setValidationMessage("Email is required");
      return;
    }
    if (email.indexOf("@") === -1) {
      setValidationMessage("Email is invalid");
      return;
    }

    if (password === "") {
      setValidationMessage("Password is required");
      return;
    }
    if (password !== confirmPassword) {
      setValidationMessage("Passwords do not match");
      return;
    }

    const res = await ApiService.register(username, email, password);
    if (res.status === 201) {
      setMessage("Check your email for a link to confirm your account.");
      return;
    }

    setValidationMessage("There was an error registering your account. Email and username may already be in use.");
  };

  if (message !== null) {
    return (
      <>
        <RegisterContent>
          <p>{message}</p>
        </RegisterContent>
      </>
    );
  }

  return (
    <>
      <RegisterContent>
        <>
          <Input required value={username} onChange={e => setUsername(e.target.value)} type="text" placeholder="Username" />
          <Input required value={email} onChange={e => setEmail(e.target.value)} type="email" placeholder="Email" />
          <Input required value={password} onChange={e => setPassword(e.target.value)} type="password" placeholder="Password" />
          <Input required value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} type="password" placeholder="Password" />
          {validationMessage && <p className="text-red-500">{validationMessage}</p>}
          <Button onClick={submitRegistration}>Register</Button>
        </>
      </RegisterContent>
    </>
  );
};

export const RegisterContent = ({ ...props }: React.ComponentProps<"div">) => {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="/" className="flex items-center gap-2 self-center font-medium text-2xl">
          <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
            <Dumbbell className="size-6" />
          </div>
          Gymotric
        </a>
      </div>
      <div className={cn("flex w-full max-w-sm flex-col gap-6")}>
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Register Account</CardTitle>
            <CardDescription>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-6 w-full" {...props} />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default Register;
