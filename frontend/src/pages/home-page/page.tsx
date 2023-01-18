import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import PostCard from "./post-card";
import { getPostsQuery } from "../../api/posts/queries";
import TopBar from "./top-bar";

function homePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getPostsQuery());
  };
}

function HomePage() {
  const { data: posts } = useQuery(getPostsQuery());

  if (!posts) return <div>No posts</div>;

  return (
    <div className="flex flex-col gap-2">
      <TopBar />

      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}

export default HomePage;
export { homePageLoader };
