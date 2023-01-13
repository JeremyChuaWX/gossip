import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { UpdateMeInput } from "../api/users/functions";
import { useUpdateMeMutation } from "../api/users/mutations";

function UpdateUserForm() {
  const { mutate: updateMe } = useUpdateMeMutation();

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
    updateMe({
      username: input.username,
      email: input.email,
      password: input.password,
    });
  };

  return (
    <form
      onSubmit={handleSubmit(submitHandler)}
      className="flex flex-col gap-4 items-center p-4 mx-auto w-1/2 rounded-lg border border-black"
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
          className="w-full border border-black"
        />
      </label>
      <input type="submit" className="p-1 rounded-lg border border-black" />
    </form>
  );
}

export default UpdateUserForm;
