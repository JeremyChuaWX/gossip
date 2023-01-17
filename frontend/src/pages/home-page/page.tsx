import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import PostCard from "./post-card";
import { getPostsQuery } from "../../api/posts/queries";
// import ToolBar from "./tool-bar";
import AddPostForm from "./add-post-form";

function homePageLoader(queryClient: QueryClient) {
  return async () => {
    return queryClient.fetchQuery(getPostsQuery());
  };
}

function HomePage() {
  const { data: posts } = useQuery(getPostsQuery());

  if (!posts) return <div>no posts</div>;

  return (
    <div className="mx-auto mt-4 w-3/4">
      <AddPostForm />
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
