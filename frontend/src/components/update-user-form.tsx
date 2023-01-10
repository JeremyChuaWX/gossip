import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdateMeInput } from "../api/users";
import { updateMe as updateMeApi } from "../api/users";

function UpdateUserForm() {
  const queryClient = useQueryClient();

  const { mutate: updateUser } = useMutation({
    mutationFn: (input: UpdateMeInput) => updateMeApi(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["get-me"] });
    },
  });

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<UpdateMeInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<UpdateMeInput> = (input) => {
    updateUser({
      username: input.username,
      email: input.email,
      password: input.password,
    });
  };

  return (
    <form
      onSubmit={handleSubmit(submitHandler)}
      className="flex flex-col w-1/2 mx-auto gap-4 border border-black rounded-lg items-center p-4"
    >
      <h1 className="text-lg">Update Details</h1>
      <label>
        username
        <input
          {...register("username")}
          className="w-full border border-black"
        />
      </label>
      <label>
        email
        <input {...register("email")} className="w-full border border-black" />
      </label>
      <label>
        password
        <input
          {...register("password")}
          type="password"
          className="border border-black w-full"
        />
      </label>
      <input type="submit" className="border border-black p-1 rounded-lg" />
    </form>
  );
}

export default UpdateUserForm;
