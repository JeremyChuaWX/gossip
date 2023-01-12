import { getMe, getUser } from "../api/users";

function getMeQuery() {
  return {
    queryKey: ["get-me"],
    queryFn: () => getMe(),
  };
}

function getUserQuery(id: string) {
  return {
    queryKey: ["get-user", id],
    queryFn: () => getUser({ id }),
  };
}

export { getUserQuery, getMeQuery };
