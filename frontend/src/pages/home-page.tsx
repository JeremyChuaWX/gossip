import { QueryClient, useQuery } from "@tanstack/react-query";
import { getPosts as getPostsApi } from "../api/posts";

function getPostsQuery() {
  return {
    queryKey: ["get-posts"],
    queryFn: () => getPostsApi(),
  };
}

function homePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getPostsQuery());
  };
}

function HomePage() {
  const { data: getPostsRes } = useQuery(getPostsQuery());
  const posts = getPostsRes?.data;

  if (!posts) return <div>no posts</div>;

  return (
    <div>
      <h1>home page</h1>
      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            {post.title}: {post.body}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default HomePage;
export { homePageLoader };
