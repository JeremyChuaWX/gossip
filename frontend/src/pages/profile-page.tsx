import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import UpdateUserForm from "../components/update-user-form";
import { getMeQuery } from "../react-queries/users";

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
      <div className="p-4">
        <h1 className="text-3xl">Hello {user.username}</h1>
        <p>email: {user.email}</p>
      </div>
      <UpdateUserForm />
    </div>
  );
}

export default ProfilePage;
export { profilePageLoader };
