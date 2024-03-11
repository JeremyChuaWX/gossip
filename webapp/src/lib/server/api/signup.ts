import { BASE_URL } from "$lib/server/constants";
import { request } from "$lib/server/utils";

type SignupParams = {
    username: string;
    password: string;
};

type SignupResponse = { sessionId: string };

export async function signup(params: SignupParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    return await request<SignupResponse>(`${BASE_URL}/auth/signup`, "post", { body });
}
