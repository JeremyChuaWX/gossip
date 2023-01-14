import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { CreateCommentInput } from "../api/comments/functions";
import { useCreateCommentMutation } from "../api/comments/mutations";

function AddCommentForm({ post_id }: { post_id: string }) {
  const { mutate: createComment } = useCreateCommentMutation();

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<CreateCommentInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<CreateCommentInput> = (input) => {
    createComment({
      post_id,
      body: input.body,
    });
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)}>
      <div className="flex gap-4">
        <label>new comment</label>
        <input type="submit" className="rounded-md border border-black" />
      </div>
      <textarea {...register("body")} className="w-full border border-black" />
    </form>
  );
}

export default AddCommentForm;
