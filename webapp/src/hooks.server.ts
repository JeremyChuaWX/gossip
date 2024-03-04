import { SESSION_ID_COOKIE } from "$lib/server/constants";
import { getCurrentUser } from "$lib/server/me";
import type { Handle } from "@sveltejs/kit";
import { redirect } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const sessionId = event.cookies.get(SESSION_ID_COOKIE);

    // https://github.com/sharu725/sveltekit-walkthrough-website/blob/master/src/hooks.server.js
    if (sessionId === undefined && event.route.id?.includes("(auth)")) {
        throw redirect(302, "/login");
    }

    // if sessionId !== undefined get the current logged in user info from "/auth/me" and details in locals
    if (sessionId !== undefined) {
        const res = await getCurrentUser({ sessionId });
        if (res === undefined) {
            throw redirect(302, "/login");
        }
        event.locals.user = res.user;
    }

    return await resolve(event);
};
