import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import UpdateUserForm from "../components/update-user-form";
import { getMeQuery } from "../api/users/queries";

function profilePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getMeQuery());
  };
}

function ProfilePage() {
  const { data: user, isLoading } = useQuery(getMeQuery());

  if (!user || isLoading) return <div>loading...</div>;

  return (
    <>
      <div className="p-4">
        <h1 className="text-3xl">Hello {user.username}</h1>
        <p>email: {user.email}</p>
      </div>
      <UpdateUserForm />
    </>
  );
}

export default ProfilePage;
export { profilePageLoader };
