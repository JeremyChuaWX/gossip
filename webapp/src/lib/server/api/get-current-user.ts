import { BASE_URL } from "$lib/server/constants";
import type { User } from "$lib/server/types";
import { request } from "$lib/server/utils";

type GetCurrentUserParams = {
    sessionId: string;
};

type GetCurrentUserResponse = { user: User };

export async function getCurrentUser(params: GetCurrentUserParams) {
    return await request<GetCurrentUserResponse>(`${BASE_URL}/auth/me`, "get", {
        sessionId: params.sessionId,
    });
}
