import axios from "axios";
import { User } from "../models/entities";
import { ServerResponse } from "../models/server";

interface SignInInput {
  username: string;
  password: string;
}

async function signIn(input: SignInInput) {
  const res = (
    await axios.post<ServerResponse<User>>(
      `http://localhost:3001/api/auth/signin`,
      input,
      { withCredentials: true }
    )
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
    await axios.post<ServerResponse<User>>(
      `http://localhost:3001/api/auth/signup`,
      input,
      { withCredentials: true }
    )
  ).data;

  if (res.error) {
    throw Error(`Error posting credentials`);
  } else {
    return res;
  }
}

async function signOut() {
  const res = (
    await axios.post<ServerResponse<any>>(
      `http://localhost:3001/api/auth/signout`,
      { withCredentials: true }
    )
  ).data;

  if (res.error) {
    throw Error(`Error posting credentials`);
  } else {
    return res;
  }
}

export type { SignInInput, SignUpInput };
export { signIn, signUp, signOut };
