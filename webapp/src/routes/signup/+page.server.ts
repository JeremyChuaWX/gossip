import type { PageServerLoad, Actions } from "./$types";
import { fail } from "@sveltejs/kit";
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
        try {
            const res = await signup({
                username: form.data.username,
                password: form.data.password,
            });
            console.log(res);
            return res;
        } catch (error) {
            console.error(error);
            return fail(500, { error });
        }
    },
};
