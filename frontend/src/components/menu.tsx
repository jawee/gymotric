import { Link } from "react-router";
import {
  NavigationMenu,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu"

const Menu = () => {
    return (
        <NavigationMenu>
            <NavigationMenuLink asChild><Link to="/app">Home</Link></NavigationMenuLink>
            <NavigationMenuLink><Link to="/app/workouts">Workouts</Link></NavigationMenuLink>
            <NavigationMenuLink><Link to="/app/exercise-types">Exercise types</Link></NavigationMenuLink>
            <NavigationMenuLink><Link to="/app/logout">Log Out</Link></NavigationMenuLink>

        </NavigationMenu>
    );
};

export default Menu;
