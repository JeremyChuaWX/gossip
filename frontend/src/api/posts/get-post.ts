import { QueryClient } from "@tanstack/react-query";
import axios from "axios";
import { Post } from "../../models/entities";
import { ServerResponse } from "../../models/server";

const TWO_MINS = 1000 * 60 * 2;

// api call to the backend
async function apiCall(id: string) {
  const res = (
    await axios.get<ServerResponse<Post>>(
      `http://localhost:3001/api/posts/get-post/${id}`
    )
  ).data;

  if (res.error) {
    throw Error(`Error fetching post ${id}`);
  } else {
    return res;
  }
}

// returns react-query object
function queryObject(id: string) {
  return {
    queryKey: ["get-post", id],
    queryFn: async () => apiCall(id),
    staleTime: TWO_MINS,
  };
}

// returns react-router-dom loader function
export function loader(queryClient: QueryClient) {
  return async ({ params }: any) => {
    const query = queryObject(params.id);
    return queryClient.fetchQuery(query);
  };
}

// export type to be used for "useLoaderData" hook
export type loaderDataType = Awaited<ReturnType<ReturnType<typeof loader>>>;
