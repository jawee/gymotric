import { Link } from "react-router";

const Menu = () => {
    return (
        <>
            <Link to="/">Home</Link>
            <Link to="/workouts">Workouts</Link>
            <Link to="/exercise-types">Exercise types</Link>
        </>
    );
};

export default Menu;
