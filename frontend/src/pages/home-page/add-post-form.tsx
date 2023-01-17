import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useCreatePostMutation } from "../../api/posts/mutations";
import type { CreatePostInput } from "../../api/posts/functions";

function AddPostForm() {
  const { mutate: createComment } = useCreatePostMutation();

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<CreatePostInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<CreatePostInput> = (input) => {
    createComment({
      title: input.title,
      body: input.body,
    });
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)}>
      <label>title</label>
      <input {...register("title")} className="w-full border border-black" />
      <label>body</label>
      <textarea {...register("body")} className="w-full border border-black" />
      <input type="submit" className="rounded-md border border-black" />
    </form>
  );
}

export default AddPostForm;
