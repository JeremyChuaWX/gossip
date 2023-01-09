import { Post } from "../models/entities";
import { ServerResponse } from "../models/server";
import { axiosConfig } from "../utils/axios";

async function getPosts() {
  const res = (await axiosConfig.get<ServerResponse<Post[]>>("posts/get-posts"))
    .data;

  if (res.error) {
    throw Error("Error fetching posts");
  } else {
    return res.data;
  }
}

interface GetPostInput {
  id: string;
}

async function getPost(input: GetPostInput) {
  const res = (
    await axiosConfig.get<ServerResponse<Post>>(`posts/get-post/${input.id}`)
  ).data;

  if (res.error) {
    throw Error(`Error fetching post ${input.id}`);
  } else {
    return res.data;
  }
}

export type { GetPostInput };
export { getPost, getPosts };
