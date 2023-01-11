import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import type { SubmitHandler } from "react-hook-form";
import { useForm } from "react-hook-form";
import { useLocation, useNavigate } from "react-router-dom";
import type { SignInInput } from "../api/auth";
import { signIn as signInApi } from "../api/auth";

function SignInPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();

  const from = (location.state?.from.pathname as string) || "/";

  const { mutate: signIn } = useMutation({
    mutationFn: (input: SignInInput) => signInApi(input),
    onSuccess: (data) => {
      queryClient.setQueryData(["get-me"], data);
      navigate(from);
    },
  });

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
    signIn(input);
  };

  return (
    <form
      onSubmit={handleSubmit(submitHandler)}
      className="flex flex-col w-1/2 mx-auto gap-4 border border-black rounded-lg items-center p-4"
    >
      <h1 className="text-lg">Sign In</h1>
      <label>
        username:
        <input {...register("username")} className="border border-black ml-2" />
      </label>
      <label>
        password:
        <input
          {...register("password")}
          type="password"
          className="border border-black ml-2"
        />
      </label>
      <input type="submit" className="border border-black p-1" />
    </form>
  );
}

export default SignInPage;
