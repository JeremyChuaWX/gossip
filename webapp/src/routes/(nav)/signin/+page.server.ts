import type { PageServerLoad, Actions } from "./$types";
import { fail, redirect } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema";
import { zod } from "sveltekit-superforms/adapters";
import { signin } from "$lib/server/signin";
import { SESSION_ID_COOKIE } from "$lib/server/constants";
import { setCookie } from "$lib/server/utils";

export const load: PageServerLoad = async () => {
    return {
        form: await superValidate(zod(formSchema)),
    };
};

export const actions: Actions = {
    default: async (event) => {
        const form = await superValidate(event, zod(formSchema));
        if (!form.valid) {
            return fail(400, { form });
        }
        const res = await signin({
            username: form.data.username,
            password: form.data.password,
        });
        if (res === undefined) {
            return fail(500);
        }
        console.log("signin response", res);
        setCookie(event.cookies, SESSION_ID_COOKIE, res.sessionId);
        throw redirect(302, "/home");
    },
};
