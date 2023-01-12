import { Outlet, useNavigate } from "react-router-dom";
import { NavLink } from "react-router-dom";
import { getMe as getMeApi } from "../api/users";
import { signOut as signOutApi } from "../api/auth";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";

function BaseLayout() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { mutate: signOut } = useMutation({
    mutationFn: () => signOutApi(),
    onSettled: () => {
      queryClient.removeQueries({ queryKey: ["get-me"] });
      navigate("/");
    },
  });

  const { data: user } = useQuery({
    queryKey: ["get-me"],
    queryFn: () => getMeApi(),
  });

  return (
    <>
      <header className="flex items-center justify-between p-4 relative">
        <NavLink to="/">
          <h1 className="font-bold text-2xl">Gossip</h1>
        </NavLink>
        <nav className="flex gap-4">
          {user ? (
            <>
              <button onClick={() => signOut()}>Sign Out</button>
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
