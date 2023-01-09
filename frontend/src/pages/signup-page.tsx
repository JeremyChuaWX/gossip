import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useLocation, useNavigate } from "react-router-dom";
import type { SignUpInput } from "../api/auth";
import { signUp as signUpApi } from "../api/auth";

function SignUpPage() {
  const navigate = useNavigate();
  const location = useLocation();

  const { mutate: signUp } = useMutation({
    mutationFn: (input: SignUpInput) => signUpApi(input),
    onSuccess: () => {
      navigate(from);
    },
  });

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
    signUp(input);
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
