import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import type { LoaderFunctionArgs } from "react-router-dom";
import { useParams } from "react-router-dom";
import { getUser as getUserApi } from "../api/users";

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

  const { data: getUserRes } = useQuery(getUserQuery(id));
  const user = getUserRes?.data;

  return <div>hello {user?.username}</div>;
}

export default UserPage;
export { userPageLoader };
