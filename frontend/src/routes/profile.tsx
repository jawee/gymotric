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
import { toast } from "sonner";
import { Toaster } from "@/components/ui/sonner";
import WtDialog from "@/components/wt-dialog";

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
        <WtDialog
          openButtonTitle={<><RectangleEllipsis /> Change password</>}
          form={
            <>
              <Input type="password" value={oldPassword} onChange={e => setOldPassword(e.target.value)} placeholder="Old password" />
              <Input type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} placeholder="New password" />
              <Input type="password" value={confirmNewPassword} onChange={e => setConfirmNewPassword(e.target.value)} placeholder="Confirm new password" />
              {changePasswordError !== null ? <p className="text-red-500">{changePasswordError}</p> : null}
            </>
          }
          onSubmitButtonClick={changePassword}
          onSubmitButtonTitle="Change Password"
          shouldUseDefaultSubmit={false}
          dialogProps={{ open: passwordDialogOpen, onOpenChange: setPasswordDialogOpen }}
          topPercentage={"25"}
          title="Change Password" />
      </div>
      <div className="mt-2">
        <WtDialog
          openButtonTitle={<><Mail /> {user.email === null ? "Register email" : "Change email"}</>}
          form={
            <>
              <Input type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="Email" />
              {registerEmailError !== null ? <p className="text-red-500">{registerEmailError}</p> : null}
            </>
          }
          onSubmitButtonClick={registerEmail}
          onSubmitButtonTitle={user.email === null ? "Register Email" : "Change Email"}
          dialogProps={{ open: emailDialogOpen, onOpenChange: setEmailDialogOpen }}
          shouldUseDefaultSubmit={false}
          title={user.email === null ? "Register Email" : "Change Email"} />
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

