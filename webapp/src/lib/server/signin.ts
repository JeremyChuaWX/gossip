import { BASE_URL } from "./constants";
import { request } from "./utils";

type SigninParams = {
    username: string;
    password: string;
};

export async function signin(params: SigninParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    return await request(`${BASE_URL}/auth/signin`, "post", body);
}
