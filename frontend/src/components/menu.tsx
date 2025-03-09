import { Link } from "react-router";
import {
  NavigationMenu,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu"
import { Dumbbell } from "lucide-react";

const Menu = () => {
    return (
        <NavigationMenu>
            <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <Dumbbell className="size-4" />
            </div>
            <NavigationMenuLink asChild><Link to="/app">Home</Link></NavigationMenuLink>
            <NavigationMenuLink asChild><Link to="/app/workouts">Workouts</Link></NavigationMenuLink>
            <NavigationMenuLink asChild><Link to="/app/exercise-types">Exercises</Link></NavigationMenuLink>
            <NavigationMenuLink asChild><Link to="/app/logout">Log Out</Link></NavigationMenuLink>
        </NavigationMenu>
    );
};

export default Menu;
