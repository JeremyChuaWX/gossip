import { NavLink } from "react-router-dom";
import type { Post } from "../../models/entities";

function PostCard({ post }: { post: Post }) {
  return (
    <NavLink to={`post/${post.id}`}>
      <div className="border-b border-gray-300">
        <h2 className="text-lg">{post.title}</h2>
        <NavLink
          className="w-max text-sm text-gray-500 duration-75 hover:text-black"
          to={`user/${post.user_id}`}
        >
          {post.user.username}
        </NavLink>
      </div>
    </NavLink>
  );
}

export default PostCard;
