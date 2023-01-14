import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import type { LoaderFunctionArgs } from "react-router-dom";
import { useParams } from "react-router-dom";
import { getCommentQuery } from "../api/comments/queries";
import { getMeQuery } from "../api/users/queries";

function commentPageLoader(queryClient: QueryClient) {
  return async ({ params: { id } }: LoaderFunctionArgs) => {
    if (!id) throw Error("Invalid id");

    return queryClient.fetchQuery(getCommentQuery(id));
  };
}

function CommentPage() {
  const { id } = useParams();
  if (!id) throw Error("No such comment");

  const { data: user } = useQuery(getMeQuery());
  const { data: cmt, isLoading } = useQuery(getCommentQuery(id));
  const isAuthor = user?.id === cmt?.user_id;

  if (!cmt || isLoading) return <div>loading...</div>;

  return (
    <div className="p-4">
      <h3>{cmt.user.username}</h3>
      <h3>{cmt.body}</h3>
      <h3>{isAuthor ? "i am the author" : "i am a viewer"}</h3>
    </div>
  );
}

export default CommentPage;
export { commentPageLoader };
