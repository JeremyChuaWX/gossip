import { fail } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { Actions, PageServerLoad } from "./$types";
import { createRoomFormSchema } from "./schemas";

export const load: PageServerLoad = async (event) => {
    return {
        user: event.locals.user,
        createRoomForm: await superValidate(zod(createRoomFormSchema)),
    };
};

export const actions: Actions = {
    createRoom: async (event) => {
        const form = await superValidate(event, zod(createRoomFormSchema));
        if (!form.valid) {
            return fail(400, { form });
        }
        const res = await createRoom({});
        if (res === undefined || res.error) {
            return fail(500);
        }
        console.log("createRoom response", res);
    },
};
