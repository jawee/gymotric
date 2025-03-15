import { Outlet } from "react-router";
import Menu from "../components/menu";

const Default = () => {
    return (
    <div className="container mx-auto">
      <div className="w-full">
        <Menu />
        <Outlet />
      </div>
    </div>
    );
};

export default Default;
