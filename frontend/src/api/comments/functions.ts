import type { Comment } from "../../models/entities";
import type { ServerResponse } from "../../models/server";
import { axiosConfig } from "../../utils/axios";

interface GetCommentInput {
  id: string;
}

async function getComment(input: GetCommentInput) {
  const res = (
    await axiosConfig.get<ServerResponse<Comment>>(
      `comments/get-comment/${input.id}`
    )
  ).data;

  return res.data;
}

interface CreateCommentInput {
  post_id: string;
  body: string;
  parent_id?: string;
}

async function createComment(input: CreateCommentInput) {
  const res = (
    await axiosConfig.put<ServerResponse<Comment>>("comments/create-comment/", {
      post_id: input.post_id,
      body: input.body,
      parent_id: input.parent_id,
    })
  ).data;

  return res.data;
}

interface UpdateCommentInput {
  id: string;
  body?: string;
}

async function updateComment(input: UpdateCommentInput) {
  const res = (
    await axiosConfig.put<ServerResponse<Comment>>(
      `comments/update-comment/${input.id}`,
      { body: input.body }
    )
  ).data;

  return res.data;
}

interface DeleteCommentInput {
  id: string;
}

async function deleteComment(input: DeleteCommentInput) {
  const res = (
    await axiosConfig.delete<ServerResponse<Comment>>(
      `comments/delete-comment${input.id}`
    )
  ).data;

  return res.data;
}

export type {
  GetCommentInput,
  CreateCommentInput,
  UpdateCommentInput,
  DeleteCommentInput,
};
export { getComment, createComment, updateComment, deleteComment };
