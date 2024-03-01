import { z } from "zod";

export const formSchema = z.object({
    username: z.string().trim().min(1),
    password: z.string().min(1),
});

export type FormSchema = typeof formSchema;
