import { BASE_URL } from "./constants";
import type { Response, User } from "./types";
import { request } from "./utils";

type SignupParams = {
    username: string;
    password: string;
};

type SignupResponse = Response<{ user: User }>;

export async function signup(params: SignupParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    const res = await request<SignupResponse>(`${BASE_URL}/auth/signup`, "post", body);
    if (res === undefined) {
        return undefined;
    }
    if (res.error) {
        console.error("signup error", res.message);
        return undefined;
    }
    return res.user;
}
