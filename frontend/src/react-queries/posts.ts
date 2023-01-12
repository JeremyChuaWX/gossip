import { getPost } from "../api/posts";

function getPostQuery(id: string) {
  return {
    queryKey: ["get-post", id],
    queryFn: () => getPost({ id }),
  };
}

export { getPostQuery };
