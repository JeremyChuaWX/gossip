import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import PostCard from "../components/post-card";
import { getPostsQuery } from "../api/posts/queries";

function homePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getPostsQuery());
  };
}

function HomePage() {
  const { data: posts } = useQuery(getPostsQuery());

  if (!posts) return <div>no posts</div>;

  return (
    <div className="p-4">
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}

export default HomePage;
export { homePageLoader };
