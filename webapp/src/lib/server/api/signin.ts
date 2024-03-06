import { BASE_URL } from "$lib/server/constants";
import type { Response } from "$lib/server/types";
import { request } from "$lib/server/utils";

type SigninParams = {
    username: string;
    password: string;
};

type SigninResponse = Response<{ sessionId: string }>;

export async function signin(params: SigninParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    const res = await request<SigninResponse>(`${BASE_URL}/auth/signin`, "post", { body });
    if (res === undefined) {
        return undefined;
    }
    if (res.error) {
        console.error("signin error", res.message);
        return undefined;
    }
    return res;
}
