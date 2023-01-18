import { Outlet } from "react-router-dom";
import NavBar from "../components/nav-bar";

function BaseLayout() {
  return (
    <div className="w-full h-full">
      <NavBar />

      <main className="mx-auto w-1/2 max-w-screen-xl">
        <Outlet />
      </main>
    </div>
  );
}

export default BaseLayout;
