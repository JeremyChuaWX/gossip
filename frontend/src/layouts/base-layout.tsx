import { Outlet } from "react-router-dom";
import { NavLink } from "react-router-dom";
import { getMe as getMeApi } from "../api/users";
import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";

function getMeQuery() {
  return {
    queryKey: ["get-me"],
    queryFn: () => getMeApi(),
  };
}

function baseLayoutLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getMeQuery());
  };
}

function BaseLayout() {
  const { data: user } = useQuery(getMeQuery());

  return (
    <>
      <header className="flex items-center justify-between p-4 relative">
        <NavLink to="/">
          <h1 className="font-bold text-2xl">Gossip</h1>
        </NavLink>
        <nav className="flex gap-4">
          {user ? (
            <>
              <button>Sign Out</button>
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
export { baseLayoutLoader };
