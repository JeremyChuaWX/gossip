import { BASE_URL } from "./constants";
import { request } from "./utils";

type SignupParams = {
    username: string;
    password: string;
};

export async function signup(params: SignupParams) {
    const body = {
        username: params.username,
        password: params.password,
    };
    return await request(`${BASE_URL}/auth/signup`, "post", body);
}
