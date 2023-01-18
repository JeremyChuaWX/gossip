import { NavLink } from "react-router-dom";
import type { Post } from "../../models/entities";

function PostCard({ post }: { post: Post }) {
  return (
    <NavLink
      className="flex flex-col duration-75 hover:scale-105 group"
      to={`post/${post.id}`}
    >
      <h2 className="text-lg font-bold capitalize">{post.title}</h2>
      <NavLink
        className="w-max text-sm text-gray-500 duration-75 hover:text-black"
        to={`user/${post.user_id}`}
      >
        {post.user.username}
      </NavLink>
      <p className="line-clamp-2">{post.body}</p>
      <span className="block max-w-0 h-px bg-black transition-all duration-75 group-hover:max-w-full" />
    </NavLink>
  );
}

export default PostCard;
