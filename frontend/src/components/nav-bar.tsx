import { useNavigate } from "react-router-dom";
import { NavLink } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getMeQuery } from "../api/users/queries";
import { useSignOutMutation } from "../api/auth/mutations";
import { ReactNode } from "react";

function NavButton({ children }: { children: ReactNode }) {
  return (
    <div className="duration-75 ease-out hover:scale-110 group">
      {children}
      <span className="block max-w-0 h-0.5 bg-black transition-all duration-75 group-hover:max-w-full" />
    </div>
  );
}

function NavBar() {
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
    <header className="flex justify-between items-center p-4">
      <NavButton>
        <NavLink className="text-2xl font-bold" to="/">
          Gossip
        </NavLink>
      </NavButton>
      <nav className="flex gap-4">
        {user ? (
          <>
            <NavButton>
              <button onClick={() => signOutOnClick()}>Sign Out</button>
            </NavButton>
            <NavButton>
              <NavLink to="/profile">{user.username}</NavLink>
            </NavButton>
          </>
        ) : (
          <>
            <NavButton>
              <NavLink to="/auth/signin">Sign In</NavLink>
            </NavButton>
            <NavButton>
              <NavLink to="/auth/signup">Sign Up</NavLink>
            </NavButton>
          </>
        )}
      </nav>
    </header>
  );
}

export default NavBar;
