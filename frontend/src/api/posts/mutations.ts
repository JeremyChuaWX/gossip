import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { DeletePostInput, UpdatePostInput } from "./functions";
import { deletePost } from "./functions";
import { updatePost } from "./functions";

function useUpdatePostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdatePostInput) => updatePost(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.setQueryData(["get-post", variables.id], data);
    },
  });
}

function useDeletePostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: DeletePostInput) => deletePost(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.resetQueries(["get-post", variables.id]);
    },
  });
}

export { useUpdatePostMutation, useDeletePostMutation };
