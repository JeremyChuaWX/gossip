import { Outlet } from "react-router-dom";
import { NavLink } from "react-router-dom";

export default function BaseLayout() {
  return (
    <>
      <header className="flex items-center justify-between p-4 relative">
        <NavLink to="/">
          <h1 className="font-bold text-2xl">Gossip</h1>
        </NavLink>
        <nav className="flex gap-4">
          <NavLink to="/auth/signin">Sign In</NavLink>
          <NavLink to="/auth/signup">Sign Up</NavLink>
        </nav>
      </header>

      <main>
        <Outlet />
      </main>
    </>
  );
}
