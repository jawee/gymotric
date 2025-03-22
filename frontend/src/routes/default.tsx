import { Outlet } from "react-router";
import Menu from "../components/menu";

const Default = () => {
    return (
    <div className="container mx-auto pl-4 pr-4 md:pl-0 md:pr-0">
      <div className="w-full">
        <Menu />
        <Outlet />
      </div>
    </div>
    );
};

export default Default;
