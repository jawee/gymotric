import { useId, useState } from "react";
import ApiService from "../services/api-service";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const Login = () => {
  const [error, setError] = useState<string | null>(null);

  const usernameFieldId = useId();
  const passwordFieldId = useId();

  const login = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const target = event.target as typeof event.target & {
      username: { value: string };
      password: { value: string };
    };
    const username = target.username.value; // typechecks!
    const password = target.password.value; // typechecks!

    const res = await ApiService.login(username, password);
    if (res.status === 200) {
      window.location.href = "/app";
      return;
    }

    setError("Login failed");
  };

  return (
    <>
      <h1>Login</h1>
      <form onSubmit={login}>
        <Input id={usernameFieldId} name="username" type="text" placeholder="Username" />
        <Input id={passwordFieldId} name="password" type="password" placeholder="Password" />
        <Button type="submit">Login</Button>
        {error ?? <p>{error}</p>}

      </form>

    </>
  )
};

export default Login;
