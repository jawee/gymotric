import { Link } from "react-router";
import {
  NavigationMenu,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu"
import { Dumbbell } from "lucide-react";

const Menu = () => {
  return (
    <div className="pb-2 mb-5">
      <NavigationMenu>
        <div className="flex h-10 w-10 items-center justify-center rounded-md text-primary m-1 ml-0 border-2 border-primary">
          <Dumbbell className="size-6" />
        </div>
        <NavigationMenuLink asChild><Link to="/app">Home</Link></NavigationMenuLink>
        <NavigationMenuLink asChild><Link to="/app/workouts">Workouts</Link></NavigationMenuLink>
        <NavigationMenuLink asChild><Link to="/app/exercise-types">Exercises</Link></NavigationMenuLink>
        <NavigationMenuLink asChild><Link to="/app/logout">Log Out</Link></NavigationMenuLink>
      </NavigationMenu>
    </div>
  );
};

export default Menu;
