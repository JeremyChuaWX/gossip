import { NavLink } from "react-router-dom";
import type { Post } from "../models/entities";

function PostCard({ post }: { post: Post }) {
  return (
    <NavLink to={`post/${post.id}`}>
      <div className="border-b-black border pb-2">
        <h2 className="text-lg">{post.title}</h2>
        <p className="text-sm">{post.body}</p>
      </div>
    </NavLink>
  );
}

export default PostCard;
