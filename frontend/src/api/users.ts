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

interface UpdateUserInput {
  id: string;
  username?: string;
  email?: string;
  password?: string;
}

async function updateUser(input: UpdateUserInput) {
  const res = (
    await axiosConfig.put<ServerResponse<User>>(
      `users/update-user/${input.id}`,
      input
    )
  ).data;

  if (res.error) {
    throw Error("Error updating user");
  } else {
    return res.data;
  }
}

interface DeleteUserInput {
  id: string;
}

async function deleteUser(input: DeleteUserInput) {
  const res = (
    await axiosConfig.delete<ServerResponse<User>>(
      `users/delete-user/${input.id}`
    )
  ).data;

  if (res.error) {
    throw Error("Error updating user");
  } else {
    return res.data;
  }
}

export type { GetUserInput, UpdateUserInput, DeleteUserInput };
export { getUser, updateUser, deleteUser };
