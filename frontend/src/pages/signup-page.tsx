import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import type { SignUpInput } from "../api/auth";
import { signUp } from "../api/auth";

export default function SignInPage() {
  const { register, handleSubmit } = useForm<SignUpInput>();
  const navigate = useNavigate();

  const { mutateAsync: signUpMutate } = useMutation({
    mutationFn: (input: SignUpInput) => {
      return signUp(input);
    },
  });

  const onSubmit = handleSubmit(async (input) => {
    await signUpMutate(input);
    navigate("/");
  });

  return (
    <form onSubmit={onSubmit}>
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
