import { useState } from "react";
import AddPostForm from "./add-post-form";

function ToolBar() {
  const [open, setOpen] = useState<boolean>(true);
  const handleClick = () => setOpen((curr) => !curr);

  return (
    <div className="mb-4">
      <button
        className="px-2 w-max rounded-md border border-gray-300 duration-75 ease-in-out hover:bg-gray-300"
        onClick={handleClick}
      >
        New Post
      </button>
      {open ?? <AddPostForm />}
    </div>
  );
}

export default ToolBar;
