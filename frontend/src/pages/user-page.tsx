import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import type { LoaderFunctionArgs } from "react-router-dom";
import { useParams } from "react-router-dom";
import { getUser as getUserApi } from "../api/users";
import UpdateUserForm from "../components/update-user-form";

function getUserQuery(id: string) {
  return {
    queryKey: ["get-user", id],
    queryFn: () => getUserApi({ id }),
  };
}

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

  return (
    <div>
      <h1>hello {user?.username}</h1>
      <UpdateUserForm id={id} />
    </div>
  );
}

export default UserPage;
export { userPageLoader };
