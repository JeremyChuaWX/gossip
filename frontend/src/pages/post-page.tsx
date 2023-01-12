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
      <h2>{post.title}</h2>
      <p>{post.body}</p>
      {user?.id === post.user_id && <UpdatePostForm id={id} />}
    </div>
  );
}

export default PostPage;
export { postPageLoader };
