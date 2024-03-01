import { BASE_URL } from "./constants";
import type { Response, User } from "./types";
import { request } from "./utils";

type SigninParams = {
    username: string;
    password: string;
};

type SigninResponse = Response<{ user: User }>;

export async function signin(params: SigninParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    const res: SigninResponse = await request(`${BASE_URL}/auth/signin`, "post", body);
    if (res.error) {
        throw new Error(res.message);
    }
    return res.user;
}
