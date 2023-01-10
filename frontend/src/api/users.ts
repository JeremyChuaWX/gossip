import { User } from "../models/entities";
import { ServerResponse } from "../models/server";
import { axiosConfig } from "../utils/axios";

interface GetUserInput {
  id: string;
}

async function getUser(input: GetUserInput) {
  const res = (
    await axiosConfig.get<ServerResponse<User>>(`users/get-user/${input.id}`)
  ).data;

  if (res.error) {
    throw Error(`Error fetching user ${input.id}`);
  } else {
    return res.data;
  }
}

async function getMe() {
  const res = (await axiosConfig.get<ServerResponse<User>>("users/get-me"))
    .data;

  if (res.error) {
    throw Error("Error fetching user");
  } else {
    return res.data;
  }
}

interface UpdateMeInput {
  username?: string;
  email?: string;
  password?: string;
}

async function updateMe(input: UpdateMeInput) {
  const res = (
    await axiosConfig.put<ServerResponse<User>>("users/update-me", input)
  ).data;

  if (res.error) {
    throw Error("Error updating user");
  } else {
    return res.data;
  }
}

async function deleteMe() {
  const res = (
    await axiosConfig.delete<ServerResponse<User>>("users/delete-me")
  ).data;

  if (res.error) {
    throw Error("Error deleting user");
  } else {
    return res.data;
  }
}

export type { GetUserInput, UpdateMeInput };
export { getUser, getMe, updateMe, deleteMe };
