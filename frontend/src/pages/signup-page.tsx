import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useLocation, useNavigate } from "react-router-dom";
import type { SignUpInput } from "../api/auth/functions";
import { useSignUpMutation } from "../api/auth/mutations";

function SignUpPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { mutate: signUp } = useSignUpMutation();

  const from = (location.state?.from.pathname as string) || "/";

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<SignUpInput>();

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<SignUpInput> = (input) => {
    signUp(input, {
      onSuccess: () => {
        navigate(from);
      },
    });
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

export default SignUpPage;
