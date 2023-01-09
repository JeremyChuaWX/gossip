import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdateUserInput } from "../api/users";
import { updateUser as updateUserApi } from "../api/users";

function UpdateUserForm({ id }: { id: string }) {
  const { mutate: updateUser } = useMutation({
    mutationFn: (input: UpdateUserInput) => updateUserApi(input),
  });

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<UpdateUserInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<UpdateUserInput> = (input) => {
    updateUser({ id, email: input.email, password: input.password });
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)}>
      <label>
        username:
        <input {...register("username")} />
      </label>
      <label>
        email:
        <input {...register("email")} />
      </label>
      <label>
        password:
        <input {...register("password")} type="password" />
      </label>
      <input type="submit" />
    </form>
  );
}

export default UpdateUserForm;
