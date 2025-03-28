import Loading from "@/components/loading";
import { User } from "../models/user";
import { useEffect, useState } from "react";
import ApiService from "@/services/api-service";
import { useNavigate } from "react-router";
import StatisticsComponent from "@/components/statistics";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { LogOut, RectangleEllipsis } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { toast } from "sonner";
import { Toaster } from "@/components/ui/sonner";

const Profile = () => {
  const [user, setUser] = useState<User | null>(null);
  const [oldPassword, setOldPassword] = useState<string>("");
  const [newPassword, setNewPassword] = useState<string>("");
  const [confirmNewPassword, setConfirmNewPassword] = useState<string>("");
  const [changePasswordError, setChangePasswordError] = useState<string | null>(null);
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

  const changePassword = async () => {
    if (oldPassword === "" || newPassword === "" || confirmNewPassword === "") {
      setChangePasswordError("All fields are required");
      return;
    }

    if (newPassword.length < 4) {
      setChangePasswordError("Password must be at least 4 characters long");
      return;
    }

    if (newPassword !== confirmNewPassword) {
      setChangePasswordError("Passwords do not match");
      return;
    }
    console.log(oldPassword, newPassword, confirmNewPassword);

    const response = await ApiService.changePassword(oldPassword, newPassword);
    if (response.status === 204) {
      setDialogOpen(false);
      toast("Password changed. You will need to log in again");
      navigate("/login", { state: { message: "Password changed. You need to log in again" } });
      return;
    }
    setChangePasswordError("Error changing password");
  };

  const [dialogOpen, setDialogOpen] = useState(false);

  if (user === null) {
    return <Loading />;
  }
  return (
    <>
      <h1 className="text-2xl">Profile for {user.username}</h1>
      <p>Member since: {new Date(user.created_on).toDateString()}</p>
      <h2 className="text-xl mb-2">Statistics</h2>
      <StatisticsComponent />
      <div className="mt-2">
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger className={buttonVariants({ variant: "default" })}><><RectangleEllipsis /> Change password</></DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Change Password</DialogTitle>
              <DialogDescription>
              </DialogDescription>
            </DialogHeader>
            <>
              <Input type="password" value={oldPassword} onChange={e => setOldPassword(e.target.value)} placeholder="Old password" />
              <Input type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} placeholder="New password" />
              <Input type="password" value={confirmNewPassword} onChange={e => setConfirmNewPassword(e.target.value)} placeholder="Confirm new password" />
              {changePasswordError !== null ? <p className="text-red-500">{changePasswordError}</p> : null}
            </>
            <DialogFooter>
              <DialogClose className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Cancel</DialogClose>
              <Button onClick={changePassword}>Change Password</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
      <div className="mt-2">
        <Button
          onClick={() => { navigate("/app/logout"); }}
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
      <Toaster />
    </>
  );
};
 

export default Profile;

