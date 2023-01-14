import { useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  CreateCommentInput,
  UpdateCommentInput,
  DeleteCommentInput,
  UpdateCommentScoreInput,
} from "./functions";
import {
  createComment,
  updateComment,
  updateCommentScore,
  deleteComment,
} from "./functions";

function useCreateCommentMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateCommentInput) => createComment(input),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", data.post_id] });
    },
  });
}

function useUpdateCommentMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateCommentInput) => updateComment(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", data.post_id] });
      queryClient.invalidateQueries({
        queryKey: ["get-comment", variables.id],
      });
    },
  });
}

function useUpdateCommentScoreMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateCommentScoreInput) => updateCommentScore(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", data.post_id] });
      queryClient.invalidateQueries({
        queryKey: ["get-comment", variables.id],
      });
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

export {
  useCreateCommentMutation,
  useUpdateCommentMutation,
  useUpdateCommentScoreMutation,
  useDeleteCommentMutation,
};
