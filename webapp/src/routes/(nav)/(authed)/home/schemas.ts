import { z } from "zod";

export const createRoomFormSchema = z.object({
    roomName: z.string().trim().min(1),
});

export type CreateRoomFormSchema = typeof createRoomFormSchema;
