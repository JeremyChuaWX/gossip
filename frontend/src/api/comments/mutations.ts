import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { UpdateCommentInput, DeleteCommentInput } from "./functions";
import { updateComment, deleteComment } from "./functions";

function useUpdateCommentMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateCommentInput) => updateComment(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", data.post_id] });
      queryClient.setQueryData(["get-comment", variables.id], data);
    },
  });
}

function useDeleteCommentMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: DeleteCommentInput) => deleteComment(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", data.post_id] });
      queryClient.resetQueries({ queryKey: ["get-comment", variables.id] });
    },
  });
}

export { useUpdateCommentMutation, useDeleteCommentMutation };
