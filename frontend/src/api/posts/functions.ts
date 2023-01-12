import type { Post } from "../../models/entities";
import type { ServerResponse } from "../../models/server";
import { axiosConfig } from "../../utils/axios";

async function getPosts() {
  const res = (await axiosConfig.get<ServerResponse<Post[]>>("posts/get-posts"))
    .data;

  return res.data;
}

interface GetPostInput {
  id: string;
}

async function getPost(input: GetPostInput) {
  const res = (
    await axiosConfig.get<ServerResponse<Post>>(`posts/get-post/${input.id}`)
  ).data;

  return res.data;
}

interface UpdatePostInput {
  id: string;
  title?: string;
  body?: string;
  post_score?: number;
}

async function updatePost(input: UpdatePostInput) {
  const res = (
    await axiosConfig.put<ServerResponse<Post>>(
      `posts/update-post/${input.id}`,
      {
        title: input.title,
        body: input.body,
        post_score: input.post_score,
      }
    )
  ).data;

  return res.data;
}

interface DeletePostInput {
  id: string;
}

async function deletePost(input: DeletePostInput) {
  const res = (
    await axiosConfig.delete<ServerResponse<Post>>(
      `posts/delete-post${input.id}`
    )
  ).data;

  return res.data;
}

export type { GetPostInput, UpdatePostInput, DeletePostInput };
export { getPost, getPosts, updatePost, deletePost };
