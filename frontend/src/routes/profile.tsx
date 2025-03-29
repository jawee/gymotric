import Loading from "@/components/loading";
import { User } from "../models/user";
import { useEffect, useState } from "react";
import ApiService from "@/services/api-service";
import { useNavigate } from "react-router";
import StatisticsComponent from "@/components/statistics";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { LogOut, Mail, RectangleEllipsis } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { toast } from "sonner";
import { Toaster } from "@/components/ui/sonner";

const Profile = () => {
  const [user, setUser] = useState<User | null>(null);
  const [email, setEmail] = useState<string>("");
  const [oldPassword, setOldPassword] = useState<string>("");
  const [newPassword, setNewPassword] = useState<string>("");
  const [confirmNewPassword, setConfirmNewPassword] = useState<string>("");
  const [changePasswordError, setChangePasswordError] = useState<string | null>(null);
  const [registerEmailError, setRegisterEmailError] = useState<string | null>(null);
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


  const registerEmail = async () => {
    if (email === "") {
      setRegisterEmailError("Email is required");
      return;
    }

    const response = await ApiService.registerEmail(email); 
    if (response.status === 204) {
      setEmailDialogOpen(false);
      toast("Email registered. A confirmation email has been sent");
      return;
    }

    setRegisterEmailError("Error registering email");
  };

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

    const response = await ApiService.changePassword(oldPassword, newPassword);
    if (response.status === 204) {
      setPasswordDialogOpen(false);
      toast("Password changed. You will need to log in again");
      navigate("/login", { state: { message: "Password changed. You need to log in again" } });
      return;
    }
    setChangePasswordError("Error changing password");
  };

  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [emailDialogOpen, setEmailDialogOpen] = useState(false);

  if (user === null) {
    return <Loading />;
  }
  return (
    <>
      <h1 className="text-2xl">Profile for {user.username}</h1>
      <p>Member since: {new Date(user.created_on).toDateString()}</p>
      <p>Registered email: {user.email}</p>
      <h2 className="text-xl mb-2">Statistics</h2>
      <StatisticsComponent />
      <div className="mt-2">
        <Dialog open={passwordDialogOpen} onOpenChange={setPasswordDialogOpen}>
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
        <Dialog open={emailDialogOpen} onOpenChange={setEmailDialogOpen}>
          <DialogTrigger className={buttonVariants({ variant: "default" })}><><Mail /> {user.email === null ? "Register email" : "Change email"}</></DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{user.email === null ? "Register email" : "Change email"}</DialogTitle>
              <DialogDescription>
              </DialogDescription>
            </DialogHeader>
            <>
              <Input type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="Email" />
              {registerEmailError !== null ? <p className="text-red-500">{registerEmailError}</p> : null}
            </>
            <DialogFooter>
              <DialogClose className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Cancel</DialogClose>
              <Button onClick={registerEmail}>Submit</Button>
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

