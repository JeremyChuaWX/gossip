import axios from "axios";
import { Post } from "../models/entities";
import { ServerResponse } from "../models/server";

async function getPosts() {
  const res = (
    await axios.get<ServerResponse<Post[]>>(
      "http://localhost:3001/api/posts/get-posts"
    )
  ).data;

  if (res.error) {
    throw Error("Error fetching posts");
  } else {
    return res;
  }
}

interface GetPostInput {
  id: string;
}

async function getPost(input: GetPostInput) {
  const res = (
    await axios.get<ServerResponse<Post>>(
      `http://localhost:3001/api/posts/get-post/${input.id}`
    )
  ).data;

  if (res.error) {
    throw Error(`Error fetching post ${input.id}`);
  } else {
    return res;
  }
}

export type { GetPostInput };
export { getPost, getPosts };
