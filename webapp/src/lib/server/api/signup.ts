import { BASE_URL } from "$lib/server/constants";
import type { Response } from "$lib/server/types";
import { request } from "$lib/server/utils";

type SignupParams = {
    username: string;
    password: string;
};

type SignupResponse = Response<{ sessionId: string }>;

export async function signup(params: SignupParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    const res = await request<SignupResponse>(`${BASE_URL}/auth/signup`, "post", { body });
    if (res === undefined) {
        return undefined;
    }
    if (res.error) {
        console.error("signup error", res.message);
        return undefined;
    }
    return res;
}
