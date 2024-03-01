import type { PageServerLoad, Actions } from "./$types";
import { fail, redirect } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema";
import { zod } from "sveltekit-superforms/adapters";
import { signup } from "$lib/server/signup";

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
        const res = await signup({
            username: form.data.username,
            password: form.data.password,
        });
        if (res === undefined) {
            return fail(500);
        }
        console.log("signup response", res);
        // set cookies
        throw redirect(302, "/home");
    },
};
