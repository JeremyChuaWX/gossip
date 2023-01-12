import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdatePostInput } from "../api/posts";
import { updatePost as updatePostApi } from "../api/posts";

function UpdatePostForm({ id }: { id: string }) {
  const queryClient = useQueryClient();

  const { mutate: updatePost } = useMutation({
    mutationFn: (input: UpdatePostInput) => updatePostApi(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", variables.id] });
    },
  });

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

  return (
    <form
      onSubmit={handleSubmit(submitHandler)}
      className="flex flex-col w-1/2 mx-auto gap-4 border border-black rounded-lg items-center p-4"
    >
      <h1 className="text-lg">Update Details</h1>
      <label>
        title
        <input {...register("title")} className="w-full border border-black" />
      </label>
      <label>
        body
        <textarea
          {...register("body")}
          className="w-full border border-black"
        />
      </label>
      <input type="submit" className="border border-black p-1 rounded-lg" />
    </form>
  );
}

export default UpdatePostForm;
