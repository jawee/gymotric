import { Button } from "@/components/ui/button";
import { Dumbbell } from "lucide-react";
import { Link } from "react-router";

const Landing = () => {
  return (
    <div className="flex min-h-svh flex-col items-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium text-2xl">
          <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
            <Dumbbell className="size-6" />
          </div>
          Gymotric
        </a>
        <Button><Link to="/login">Log In</Link></Button>
      </div>
    </div>
  );
}

export default Landing;
