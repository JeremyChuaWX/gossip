import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import { getMe as getMeApi } from "../api/users";
import UpdateUserForm from "../components/update-user-form";

function getMeQuery() {
  return {
    queryKey: ["get-me"],
    queryFn: () => getMeApi(),
  };
}

function profilePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getMeQuery());
  };
}

function ProfilePage() {
  const { data: user } = useQuery(getMeQuery());

  if (!user) throw Error("No user found");

  return (
    <div>
      <h1>Hello {user.username}</h1>
      <UpdateUserForm id={user.id} />
    </div>
  );
}

export default ProfilePage;
export { profilePageLoader };
