import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import type { SignInInput } from "../api/auth";
import { signIn } from "../api/auth";

export default function SignInPage() {
  const { register, handleSubmit } = useForm<SignInInput>();
  const navigate = useNavigate();

  const { mutateAsync: signInMutate } = useMutation({
    mutationFn: (input: SignInInput) => {
      return signIn(input);
    },
  });

  const onSubmit = handleSubmit(async (input) => {
    await signInMutate(input);
    navigate("/");
  });

  return (
    <form onSubmit={onSubmit}>
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
