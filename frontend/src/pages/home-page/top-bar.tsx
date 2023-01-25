import AddPostForm from "./add-post-form";
import type { MouseEventHandler } from "react";
import { useState } from "react";

function TopBarButton({
  text,
  clickHandler,
}: {
  text: string;
  clickHandler: MouseEventHandler<HTMLButtonElement>;
}) {
  return (
    <button
      className="px-1 w-max rounded-md border border-black hover:text-white hover:bg-black"
      onClick={clickHandler}
    >
      {text}
    </button>
  );
}

function useButtonStates() {
  interface ButtonStates {
    newPost: boolean;
    filter: boolean;
  }

  const [buttonStates, setButtonStates] = useState<ButtonStates>({
    newPost: false,
    filter: false,
  });

  const flipButtonState = (button: keyof ButtonStates) => {
    const isOpen = buttonStates[button];
    setButtonStates({ newPost: false, filter: false });
    !isOpen &&
      setButtonStates((curr) => ({ ...curr, [button]: !curr[button] }));
  };

  return [buttonStates, flipButtonState] as const;
}

function TopBar() {
  const [buttonStates, flipButtonState] = useButtonStates();

  return (
    <div className="flex flex-col gap-1 p-2 rounded-md border border-black">
      <div className="flex gap-2">
        <TopBarButton
          text="New Post"
          clickHandler={() => flipButtonState("newPost")}
        />
        <TopBarButton
          text="Filter"
          clickHandler={() => flipButtonState("filter")}
        />
      </div>
      {buttonStates.newPost && <AddPostForm />}
      {buttonStates.filter && <div>filtering menu</div>}
    </div>
  );
}

export default TopBar;
