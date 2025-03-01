import { Link } from "react-router";

const Menu = () => {
    return (
        <>
            <Link to="/app">Home</Link>
            <Link to="/app/workouts">Workouts</Link>
            <Link to="/app/exercise-types">Exercise types</Link>
            <Link to="/login">Log Out</Link>
        </>
    );
};

export default Menu;
