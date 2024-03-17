import { BASE_URL } from "$lib/server/constants";
import { request } from "$lib/server/utils";

type SigninParams = {
    username: string;
    password: string;
};

type SigninResponse = { sessionId: string };

export async function signin(params: SigninParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    return await request<SigninResponse>(`${BASE_URL}/auth/signin`, "post", { body });
}
