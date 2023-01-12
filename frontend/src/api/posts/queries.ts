import { getPost, getPosts } from "./functions";

function getPostsQuery() {
  return {
    queryKey: ["get-posts"],
    queryFn: () => getPosts(),
  };
}

function getPostQuery(id: string) {
  return {
    queryKey: ["get-post", id],
    queryFn: () => getPost({ id }),
  };
}

export { getPostQuery, getPostsQuery };
