import { Outlet, useNavigate } from "react-router-dom";
import { NavLink } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getMeQuery } from "../api/users/queries";
import { useSignOutMutation } from "../api/auth/mutations";

function BaseLayout() {
  const navigate = useNavigate();

  const { mutate: signOut } = useSignOutMutation();
  const { data: user } = useQuery(getMeQuery());

  const signOutOnClick = () => {
    signOut(undefined, {
      onSuccess: () => {
        navigate("/");
      },
    });
  };

  return (
    <>
      <header className="flex relative justify-between items-center p-4">
        <NavLink to="/">
          <h1 className="text-2xl font-bold">Gossip</h1>
        </NavLink>
        <nav className="flex gap-4">
          {user ? (
            <>
              <button onClick={() => signOutOnClick()}>Sign Out</button>
              <NavLink to="/profile">{user.username}</NavLink>
            </>
          ) : (
            <>
              <NavLink to="/auth/signin">Sign In</NavLink>
              <NavLink to="/auth/signup">Sign Up</NavLink>
            </>
          )}
        </nav>
      </header>

      <main>
        <Outlet />
      </main>
    </>
  );
}

export default BaseLayout;
