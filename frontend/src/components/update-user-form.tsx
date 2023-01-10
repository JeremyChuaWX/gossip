import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdateMeInput } from "../api/users";
import { updateMe as updateMeApi } from "../api/users";

function UpdateUserForm() {
  const { mutate: updateUser } = useMutation({
    mutationFn: (input: UpdateMeInput) => updateMeApi(input),
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
    updateUser({ email: input.email, password: input.password });
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
