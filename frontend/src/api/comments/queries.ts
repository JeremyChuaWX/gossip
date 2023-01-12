import { getComment } from "./functions";

function getCommentQuery(id: string) {
  return {
    queryKey: ["get-comment", id],
    queryFn: () => getComment({ id }),
  };
}

export { getCommentQuery };
