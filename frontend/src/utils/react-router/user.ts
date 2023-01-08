import type { QueryClient } from "@tanstack/react-query";
import type { LoaderFunctionArgs } from "react-router-dom";
import { getUserQuery } from "../react-query/user";

function getUserLoader(queryClient: QueryClient) {
  return async ({ params }: LoaderFunctionArgs) => {
    // input validation
    if (!params.id) throw Error("Invalid id");
    const id = params.id;

    const query = getUserQuery(id);

    return queryClient.fetchQuery(query);
  };
}

type getUserLoaderDataType = Awaited<
  ReturnType<ReturnType<typeof getUserLoader>>
>;

export type { getUserLoaderDataType };
export { getUserLoader };
