import ApiService from "@/services/api-service";
import { useState } from "react";
import { useParams } from "react-router";
import { PasswordResetContent } from "./reset-password";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const PasswordResetConfirm = () => {
  const { token } = useParams() as { token: string };
  const [message, setMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [password, setPassword] = useState<string>("");
  const [confirmPassword, setConfirmPassword] = useState<string>("");

  if (token === undefined || token === null || token.length === 0) {
    setMessage("Invalid token. Please try again.");
    return;
  }

  const confirmPasswordReset = async () => {
    if (password === null || password.length === 0) {
      setErrorMessage("Please enter a new password.");
      return;
    }

    if (password !== confirmPassword) {
      setErrorMessage("Passwords do not match.");
      return;
    }

    const res = await ApiService.passwordResetConfirm(password, token);
    if (res.status === 204) {
      setMessage("Password confirmed successfully!");
      return;
    }
    const text = await res.text();
    console.log(text);
    setErrorMessage("Token expired or invalid.");
  };

  if (message !== null) {
    return (
      <PasswordResetContent>
        <p>{message}</p>
      </PasswordResetContent>
    );
  }

  return (
    <PasswordResetContent>
      <Input
        type="password"
        placeholder="New Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <Input
        type="password"
        placeholder="Confirm Password"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
      />
      {errorMessage && <p className="text-red-500">{errorMessage}</p>}
      <Button onClick={confirmPasswordReset}>Set Password</Button>
    </PasswordResetContent>
  );
};

export default PasswordResetConfirm;
