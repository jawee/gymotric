import { Link } from "react-router";
import {
  NavigationMenu,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu"
import { Dumbbell } from "lucide-react";

const Menu = () => {
  return (
    <div className="flex">
      <NavigationMenu className="w-full justify-evenly max-w-none text-center md:max-w-xs">
        <div className="flex h-10 w-14 items-center justify-center rounded-none text-primary m-1 ml-0 border-r-2 border-primary">
          <Dumbbell className="size-6" />
        </div>
        <NavigationMenuLink className="border-r-2 rounded-none p-2 grow border-primary" asChild><Link to="/app">Home</Link></NavigationMenuLink>
        <NavigationMenuLink className="border-r-2 rounded-none p-2 grow border-primary" asChild><Link to="/app/exercise-types">Exercises</Link></NavigationMenuLink>
        <NavigationMenuLink className="rounded-none p-2 grow" asChild><Link to="/app/profile">Profile</Link></NavigationMenuLink>
      </NavigationMenu>
    </div>
  );
};

export default Menu;
