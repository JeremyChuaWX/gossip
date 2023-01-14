import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import PostCard from "../components/post-card";
import { getPostsQuery } from "../api/posts/queries";

function homePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getPostsQuery());
  };
}

function ToolBar() {
  return (
    <div className="mb-4">
      <button className="px-2 w-max rounded-md border border-gray-300 duration-75 ease-in-out hover:bg-gray-300">
        New Post
      </button>
    </div>
  );
}

function HomePage() {
  const { data: posts } = useQuery(getPostsQuery());

  if (!posts) return <div>no posts</div>;

  return (
    <div className="mx-auto mt-4 w-3/4">
      <ToolBar />
      <div className="flex flex-col gap-2">
        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>
    </div>
  );
}

export default HomePage;
export { homePageLoader };
