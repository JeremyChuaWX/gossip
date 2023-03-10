import { useEffect, useState } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdatePostInput } from "../api/posts/functions";
import { useUpdatePostMutation } from "../api/posts/mutations";

function UpdatePostForm({ id }: { id: string }) {
  const { mutate: updatePost } = useUpdatePostMutation();
  const [showForm, setShowForm] = useState<boolean>(false);

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<UpdatePostInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<UpdatePostInput> = (input) => {
    updatePost({
      id,
      body: input.body,
      title: input.title,
    });
  };

  const showFormOnClick = () => setShowForm((curr) => !curr);

  return (
    <div>
      <button onClick={showFormOnClick}>edit</button>
      {showForm && (
        <form
          onSubmit={handleSubmit(submitHandler)}
          className="flex flex-col gap-4 items-center p-4 mx-auto w-1/2 rounded-lg border border-black"
        >
          <h1 className="text-lg">Update Details</h1>
          <label>
            title
            <input
              {...register("title")}
              className="w-full border border-black"
            />
          </label>
          <label>
            body
            <textarea
              {...register("body")}
              className="w-full border border-black"
            />
          </label>
          <input type="submit" className="p-1 rounded-lg border border-black" />
        </form>
      )}
    </div>
  );
}

export default UpdatePostForm;
