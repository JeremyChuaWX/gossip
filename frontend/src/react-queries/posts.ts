import { getPost, getPosts } from "../api/posts";

function getPostQuery(id: string) {
  return {
    queryKey: ["get-post", id],
    queryFn: () => getPost({ id }),
  };
}

function getPostsQuery() {
  return {
    queryKey: ["get-posts"],
    queryFn: () => getPosts(),
  };
}

export { getPostQuery, getPostsQuery };
