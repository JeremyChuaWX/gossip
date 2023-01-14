import { QueryClient, useQuery } from "@tanstack/react-query";
import { LoaderFunctionArgs, useParams } from "react-router-dom";
import UpdatePostForm from "../components/update-post-form";
import { getPostQuery } from "../api/posts/queries";
import { getMeQuery } from "../api/users/queries";
import AddCommentForm from "../components/add-comment-form";
import { NavLink } from "react-router-dom";

function postPageLoader(queryClient: QueryClient) {
  return async ({ params: { id } }: LoaderFunctionArgs) => {
    if (!id) throw Error("Invalid id");

    return queryClient.fetchQuery(getPostQuery(id));
  };
}

function PostPage() {
  const { id } = useParams();
  if (!id) throw Error("No such post");

  const { data: user } = useQuery(getMeQuery());
  const { data: post, isLoading } = useQuery(getPostQuery(id));
  const isAuthor = user?.id === post?.user_id;

  if (!post || isLoading) return <div>loading...</div>;

  return (
    <div className="flex flex-col gap-4 p-4">
      <div className="border border-black">
        <h2>title: {post.title}</h2>
        <h3>author: {post.user.username}</h3>
        <p>body: {post.body}</p>
        {isAuthor && <UpdatePostForm id={id} />}
      </div>
      <AddCommentForm post_id={post.id} />
      {post.comments.map((cmt) => (
        <NavLink
          className="flex flex-col gap-1"
          key={cmt.id}
          to={`comment/${cmt.id}`}
        >
          <p>{cmt.user.username}</p>
          <p>{cmt.body}</p>
        </NavLink>
      ))}
    </div>
  );
}

export default PostPage;
export { postPageLoader };
