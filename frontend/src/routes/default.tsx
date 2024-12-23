import { Outlet } from "react-router";
import Menu from "../components/menu";

const Default = () => {
    return (
    <>
        <Menu />
        <Outlet />
    </>
    );
};

export default Default;
