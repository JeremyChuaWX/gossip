import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { UpdateMeInput } from "./functions";
import { updateMe, deleteMe } from "./functions";

function useUpdateMeMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateMeInput) => updateMe(input),
    onSuccess: (data) => {
      queryClient.setQueryData(["get-me"], data);
    },
  });
}

function useDeleteMeMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => deleteMe(),
    onSuccess: () => {
      queryClient.resetQueries({ queryKey: ["get-me"] });
    },
  });
}

export { useUpdateMeMutation, useDeleteMeMutation };
