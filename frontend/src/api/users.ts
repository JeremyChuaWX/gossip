import axios from "axios";
import { User } from "../models/entities";
import { ServerResponse } from "../models/server";

interface GetUserInput {
  id: string;
}

async function getUser(input: GetUserInput) {
  const res = (
    await axios.get<ServerResponse<User>>(
      `http://localhost:3001/api/users/get-user/${input.id}`,
      { withCredentials: true }
    )
  ).data;

  if (res.error) {
    throw Error(`Error fetching user ${input.id}`);
  } else {
    return res;
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
    await axios.post<ServerResponse<User>>(
      `http://localhost:3001/api/users/update-user/${input.id}`,
      input,
      { withCredentials: true }
    )
  ).data;

  if (res.error) {
    throw Error("Error updating user");
  } else {
    return res;
  }
}

interface DeleteUserInput {
  id: string;
}

async function deleteUser(input: DeleteUserInput) {
  const res = (
    await axios.post<ServerResponse<User>>(
      `http://localhost:3001/api/users/delete-user/${input.id}`,
      { withCredentials: true }
    )
  ).data;

  if (res.error) {
    throw Error("Error updating user");
  } else {
    return res;
  }
}

export type { GetUserInput, UpdateUserInput, DeleteUserInput };
export { getUser, updateUser, deleteUser };
