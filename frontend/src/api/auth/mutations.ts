import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { SignInInput, SignUpInput } from "./functions";
import { signUp, signOut, signIn } from "./functions";

function useSignInMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: SignInInput) => signIn(input),
    onSuccess: (data) => {
      queryClient.setQueryData(["get-me"], data);
    },
  });
}

function useSignUpMutation() {
  return useMutation({
    mutationFn: (input: SignUpInput) => signUp(input),
  });
}

function useSignOutMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => signOut(),
    onSettled: () => {
      queryClient.resetQueries({ queryKey: ["get-me"] });
    },
  });
}

export { useSignOutMutation, useSignUpMutation, useSignInMutation };
