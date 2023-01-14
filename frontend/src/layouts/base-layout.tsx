import { Outlet } from "react-router-dom";
import NavBar from "../components/nav-bar";

function BaseLayout() {
  return (
    <div className="w-full h-full">
      <NavBar />

      <main className="w-full h-full">
        <Outlet />
      </main>
    </div>
  );
}

export default BaseLayout;
