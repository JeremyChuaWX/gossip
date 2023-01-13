import { QueryClient, useQuery } from "@tanstack/react-query";
import { LoaderFunctionArgs, useParams } from "react-router-dom";
import UpdatePostForm from "../components/update-post-form";
import { getPostQuery } from "../api/posts/queries";
import { getMeQuery } from "../api/users/queries";

function postPageLoader(queryClient: QueryClient) {
  return async ({ params: { id } }: LoaderFunctionArgs) => {
    if (!id) throw Error("Invalid id");

    return queryClient.fetchQuery(getPostQuery(id));
  };
}

function PostPage() {
  const { id } = useParams();

  const { data: user } = useQuery(getMeQuery());
  if (!id) throw Error("Invalid url params");

  const { data: post } = useQuery(getPostQuery(id));
  if (!post) throw Error("No such post");

  return (
    <div>
      <div className="border border-black">
        <h2>title: {post.title}</h2>
        <h3>author: {post.user.username}</h3>
        <p>body: {post.body}</p>
        {user?.id === post.user_id && <UpdatePostForm id={id} />}
      </div>
      {post.comments.map((cmt) => (
        <p key={cmt.id}>{cmt.body}</p>
      ))}
    </div>
  );
}

export default PostPage;
export { postPageLoader };
