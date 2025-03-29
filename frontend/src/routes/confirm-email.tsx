import { Button } from "@/components/ui/button";
import ApiService from "../services/api-service";
import { Dumbbell } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import Loading from "../components/loading";

const ConfirmEmail = () => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");
  const [message, setMessage] = useState("Your email has been confirmed successfully!");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const confirmEmail = async () => {
      if (token === undefined || token === null || token.length === 0) {
        navigate("/login");
        return;
      }

      const res = await ApiService.confirmEmail(token);

      if (res.status === 204) {
        setLoading(false);
        return;
      }

      console.log(await res.text());
      setMessage("Your token is invalid or has expired. Please try again.");
      setLoading(false);
      return;
    };

    confirmEmail();

  }, []);

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="flex min-h-svh flex-col items-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium text-2xl">
          <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
            <Dumbbell className="size-6" />
          </div>
          Gymotric
        </a>
      </div>
      <div className="flex w-full max-w-sm flex-col gap-6">
        <p>{message}</p>
        <Button><Link to="/login">Log In</Link></Button>
      </div>
    </div>
  );
};

export default ConfirmEmail;
