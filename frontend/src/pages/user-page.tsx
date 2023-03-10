import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import type { LoaderFunctionArgs } from "react-router-dom";
import { useParams } from "react-router-dom";
import { getUserQuery } from "../api/users/queries";

function userPageLoader(queryClient: QueryClient) {
  return async ({ params: { id } }: LoaderFunctionArgs) => {
    if (!id) throw Error("Invalid id");

    return queryClient.fetchQuery(getUserQuery(id));
  };
}

function UserPage() {
  const { id } = useParams();

  if (!id) throw Error("Invalid url params");

  const { data: user } = useQuery(getUserQuery(id));

  if (!user) throw Error("No user found");

  return (
    <div>
      <h1>{user.username}</h1>
      <h1>{user.email}</h1>
    </div>
  );
}

export default UserPage;
export { userPageLoader };
