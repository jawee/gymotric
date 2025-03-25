import Loading from "@/components/loading";
import { User } from "../models/user";
import { useEffect, useState } from "react";
import ApiService from "@/services/api-service";
import { useNavigate } from "react-router";
import StatisticsComponent from "@/components/statistics";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { LogOut } from "lucide-react";

const Profile = () => {
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchMe = async () => {
      const response = await ApiService.fetchMe();
      if (response.status === 200) {
        const user = await response.json();
        setUser(user);
        return;
      }

      navigate("/login");
    };

    fetchMe();

  }, []);

  if (user === null) {
    return <Loading />;
  }
  return (
    <>
      <h1 className="text-2xl">Profile for {user.username}</h1>
      <p>Member since: {new Date(user.created_on).toDateString()}</p>
      <h2 className="text-xl">Statistics</h2>
      <StatisticsComponent />
      <div className="mt-2">
        <Button
          onClick={() => { navigate("/app/logout");}}
          className={
            cn(buttonVariants({ variant: "default" }),
              "bg-red-500",
              "hover:bg-red-700"
            )
          }>
          <LogOut />
          Log Out
        </Button>
      </div>
    </>
  );
};


export default Profile;
