import { useNavigate } from "react-router";
import ApiService from "../services/api-service";
import { useEffect } from "react";

const Logout = () => {
  const navigate = useNavigate();
  useEffect(() => {
    const logout = async () => {
      await ApiService.logout();

      navigate("/login");
    }

    logout();
  }, [navigate]);

  return (
    <>
      <h1>Logging out...</h1>
    </>
  );
};

export default Logout;
