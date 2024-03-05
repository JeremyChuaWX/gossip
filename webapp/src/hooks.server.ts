import { SESSION_ID_COOKIE } from "$lib/server/constants";
import { getCurrentUser } from "$lib/server/me";
import type { Handle } from "@sveltejs/kit";
import { redirect } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const sessionId = event.cookies.get(SESSION_ID_COOKIE);

    // https://github.com/sharu725/sveltekit-walkthrough-website/blob/master/src/hooks.server.js
    if (sessionId === undefined && event.route.id?.includes("(authed)")) {
        throw redirect(302, "/signin");
    }

    if (sessionId !== undefined) {
        const res = await getCurrentUser({ sessionId });
        if (res === undefined) {
            throw redirect(302, "/signin");
        }
        event.locals.user = res.user;
    }

    return await resolve(event);
};
