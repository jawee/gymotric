import ApiService from "@/services/api-service";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { PasswordResetContent } from "./reset-password";
import Loading from "@/components/loading";

const PasswordResetConfirm = () => {
  const { token } = useParams() as { token: string };
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    const confirmPasswordReset = async () => {
      if (token === undefined || token === null || token.length === 0) {
        setMessage("Invalid token. Please try again.");
        return;
      }

      const res = await ApiService.passwordResetConfirm(token);
      if (res.status === 204) {
        setMessage("Password confirmed successfully!");
        return;
      }
      const text = await res.text();
      console.log(text);
      setMessage("Invalid token. Please try again.");
    };

    confirmPasswordReset();
  }, []);

  if (message === null) {
    <PasswordResetContent>
      <Loading />
    </PasswordResetContent>;
  }

  return (
    <>
      <PasswordResetContent>
        <p>{message}</p>
      </PasswordResetContent>
    </>
  );
};

export default PasswordResetConfirm;
