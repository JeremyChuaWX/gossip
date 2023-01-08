import { Outlet } from "react-router-dom";
import { NavLink } from "react-router-dom";

export default function BaseLayout() {
  return (
    <>
      <header>
        <h1>Gossip</h1>
        <nav>
          <NavLink to="/">Home</NavLink>
          <NavLink to="/auth/signin">Sign In</NavLink>
          <NavLink to="/auth/signup">Sign Up</NavLink>
          <NavLink to="/auth/signout">Sign Out</NavLink>
        </nav>
      </header>

      <main>
        <Outlet />
      </main>
    </>
  );
}
