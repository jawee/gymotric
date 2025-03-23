import LoginForm from "@/components/login-form";
import ApiService from "../services/api-service";
import { Dumbbell } from "lucide-react";
import { useEffect } from "react";
import { useNavigate } from "react-router";

const Login = () => {
  const navigate = useNavigate();
  useEffect(() => {
    const refreshToken = async () => {
      const res = await ApiService.refreshToken();
      if (res.status === 200) {
        navigate("/app");
      }
    }

    refreshToken();
  }, [navigate]);

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium text-2xl">
          <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
            <Dumbbell className="size-6" />
          </div>
          Gymotric
        </a>
        <LoginForm />
      </div>
    </div>
  )
};

export default Login;
