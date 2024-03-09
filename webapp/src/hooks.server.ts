import { SESSION_ID_COOKIE } from "$lib/server/constants";
import { getCurrentUser } from "$lib/server/api/get-current-user";
import type { Handle } from "@sveltejs/kit";
import { redirect } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const sessionId = event.cookies.get(SESSION_ID_COOKIE);

    // https://github.com/sharu725/sveltekit-walkthrough-website/blob/master/src/hooks.server.js
    if (event.route.id?.includes("(authed)")) {
        if (sessionId === undefined) {
            throw redirect(302, "/signin");
        } else {
            const res = await getCurrentUser({ sessionId });
            if (res === undefined) {
                throw redirect(302, "/signin");
            }
            event.locals.user = res.user;
        }
    }

    return await resolve(event);
};
