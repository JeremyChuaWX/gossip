import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useLocation, useNavigate } from "react-router-dom";
import type { SignInInput } from "../api/auth/functions";
import { useSignInMutation } from "../api/auth/mutations";

function SignInPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { mutate: signIn } = useSignInMutation();

  const from = (location.state?.from.pathname as string) || "/";

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<SignInInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<SignInInput> = (input) => {
    signIn(input, {
      onSuccess: () => {
        navigate(from);
      },
    });
  };

  return (
    <form
      onSubmit={handleSubmit(submitHandler)}
      className="flex flex-col gap-4 items-center p-4 mx-auto w-1/2 rounded-lg border border-black"
    >
      <h1 className="text-lg">Sign In</h1>
      <label>
        username:
        <input {...register("username")} className="ml-2 border border-black" />
      </label>
      <label>
        password:
        <input
          {...register("password")}
          type="password"
          className="ml-2 border border-black"
        />
      </label>
      <input type="submit" className="p-1 border border-black" />
    </form>
  );
}

export default SignInPage;
