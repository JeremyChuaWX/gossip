import { getUser } from "../../api/users";

const TWO_MINS = 1000 * 60 * 2;

function getUserQuery(id: string) {
  return {
    queryKey: ["get-user", id],
    queryFn: async () => getUser({ id }),
    staleTime: TWO_MINS,
  };
}

export { getUserQuery };
