import { Link } from "react-router";

const Menu = () => {
    return (
        <>
            <Link to="/app">Home</Link><span> | </span>
            <Link to="/app/workouts">Workouts</Link><span> | </span>
            <Link to="/app/exercise-types">Exercise types</Link><span> | </span>
            <Link to="/app/logout">Log Out</Link>
        </>
    );
};

export default Menu;
