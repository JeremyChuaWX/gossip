<script lang="ts">
    import * as Form from "$lib/components/ui/form";
    import { Input } from "$lib/components/ui/input";
    import { createRoomFormSchema, type CreateRoomFormSchema } from "./schemas";
    import { type SuperValidated, type Infer, superForm } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";
    import * as Dialog from "$lib/components/ui/dialog";

    export let data: SuperValidated<Infer<CreateRoomFormSchema>>;

    const form = superForm(data, {
        validators: zodClient(createRoomFormSchema),
    });

    const { form: formData, enhance } = form;
</script>

<Dialog.Root>
    <Dialog.Trigger>+</Dialog.Trigger>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Create room</Dialog.Title>
        </Dialog.Header>

        <form method="post" action="?/createRoom" use:enhance class="flex flex-col items-center">
            <Form.Field {form} name="roomName" class="w-full">
                <Form.Control let:attrs>
                    <Form.Label class="text-base">Room name</Form.Label>
                    <Input
                        {...attrs}
                        class="h-fit text-base"
                        type="text"
                        bind:value={$formData.roomName}
                    />
                </Form.Control>
                <Form.FieldErrors />
            </Form.Field>
        </form>

        <Dialog.Footer>
            <button>Submit</button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
