import { User } from "../models/entities";
import { ServerResponse } from "../models/server";
import { axiosConfig } from "../utils/axios";

interface SignInInput {
  username: string;
  password: string;
}

async function signIn(input: SignInInput) {
  const res = (
    await axiosConfig.post<ServerResponse<User>>(`auth/signin`, input)
  ).data;

  if (res.error) {
    throw Error(`Error posting credentials`);
  } else {
    return res;
  }
}

interface SignUpInput {
  username: string;
  email?: string;
  password: string;
}

async function signUp(input: SignUpInput) {
  const res = (
    await axiosConfig.post<ServerResponse<User>>(`auth/signup`, input)
  ).data;

  if (res.error) {
    throw Error(`Error posting credentials`);
  } else {
    return res;
  }
}

async function signOut() {
  const res = (await axiosConfig.post<ServerResponse<any>>(`auth/signout`))
    .data;

  if (res.error) {
    throw Error(`Error posting credentials`);
  } else {
    return res;
  }
}

export type { SignInInput, SignUpInput };
export { signIn, signUp, signOut };
