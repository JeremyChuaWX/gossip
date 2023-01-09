import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useLocation, useNavigate } from "react-router-dom";
import type { SignInInput } from "../api/auth";
import { signIn as signInApi } from "../api/auth";

function SignInPage() {
  const navigate = useNavigate();
  const location = useLocation();

  const from = (location.state?.from.pathname as string) || "/";

  const {
    reset,
    register,
    handleSubmit,
    formState: { isSubmitSuccessful },
  } = useForm<SignInInput>();

  const { mutate: signIn } = useMutation({
    mutationFn: (input: SignInInput) => signInApi(input),
    onSuccess: () => {
      navigate(from);
    },
  });

  useEffect(() => {
    if (isSubmitSuccessful) {
      reset();
    }
  }, [isSubmitSuccessful]);

  const submitHandler: SubmitHandler<SignInInput> = (input) => {
    signIn(input);
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)}>
      <label>
        username:
        <input {...register("username")} />
      </label>
      <label>
        password:
        <input {...register("password")} type="password" />
      </label>
      <input type="submit" />
    </form>
  );
}

export default SignInPage;
