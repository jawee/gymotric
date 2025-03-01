import { Link } from "react-router";

const Landing = () => {
    return (
        <>
            <h1>Welcome</h1>
            <p>Landing page</p>
            <Link to="/login">Log In</Link>
        </>
    );
}

export default Landing;
