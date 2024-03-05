import { BASE_URL } from "./constants";
import type { Response, User } from "./types";
import { request } from "./utils";

type GetCurrentUserParams = {
    sessionId: string;
};

type GetCurrentUserResponse = Response<{ user: User }>;

export async function getCurrentUser(params: GetCurrentUserParams) {
    const res = await request<GetCurrentUserResponse>(`${BASE_URL}/auth/me`, "get", {
        sessionId: params.sessionId,
    });
    if (res === undefined) {
        return undefined;
    }
    if (res.error) {
        console.error("signin error", res.message);
        return undefined;
    }
    return res;
}
