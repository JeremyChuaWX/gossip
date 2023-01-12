import type { User } from "../../models/entities";
import type { ServerResponse } from "../../models/server";
import { axiosConfig } from "../../utils/axios";

interface SignInInput {
  username: string;
  password: string;
}

async function signIn(input: SignInInput) {
  const res = (
    await axiosConfig.post<ServerResponse<User>>("auth/signin", input)
  ).data;

  return res.data;
}

interface SignUpInput {
  username: string;
  email?: string;
  password: string;
}

async function signUp(input: SignUpInput) {
  const res = (
    await axiosConfig.post<ServerResponse<User>>("auth/signup", input)
  ).data;

  return res.data;
}

async function signOut() {
  const res = (await axiosConfig.get<ServerResponse<any>>("auth/signout")).data;

  return res.data;
}

export type { SignInInput, SignUpInput };
export { signIn, signUp, signOut };
