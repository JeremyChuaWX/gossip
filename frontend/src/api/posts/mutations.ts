import { useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  CreatePostInput,
  DeletePostInput,
  UpdatePostInput,
  UpdatePostScoreInput,
} from "./functions";
import {
  createPost,
  updatePost,
  updatePostScore,
  deletePost,
} from "./functions";

function useCreatePostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreatePostInput) => createPost(input),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.setQueryData(["get-post", data.id], data);
    },
  });
}

function useUpdatePostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdatePostInput) => updatePost(input),
    onSuccess: (variables) => {
      queryClient.invalidateQueries({ queryKey: ["get-posts"] });
      queryClient.invalidateQueries({ queryKey: ["get-post", variables.id] });
    },
  });
}

function useUpdatePostScoreMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdatePostScoreInput) => updatePostScore(input),
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

export {
  useCreatePostMutation,
  useUpdatePostMutation,
  useUpdatePostScoreMutation,
  useDeletePostMutation,
};
