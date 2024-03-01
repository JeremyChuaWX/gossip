<script lang="ts">
    import * as Form from "$lib/components/ui/form";
    import { Input } from "$lib/components/ui/input";
    import { formSchema, type FormSchema } from "./schema";
    import { type SuperValidated, type Infer, superForm } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";

    export let data: SuperValidated<Infer<FormSchema>>;

    const form = superForm(data, {
        validators: zodClient(formSchema),
    });

    const { form: formData, enhance } = form;
</script>

<form method="post" use:enhance class="flex w-full flex-col items-center gap-4 md:w-1/2 lg:w-1/3">
    <Form.Field {form} name="username" class="w-full">
        <Form.Control let:attrs>
            <Form.Label class="text-base">Username</Form.Label>
            <Input {...attrs} class="h-fit text-base" type="text" bind:value={$formData.username} />
        </Form.Control>
        <Form.FieldErrors />
    </Form.Field>
    <Form.Field {form} name="password" class="w-full">
        <Form.Control let:attrs>
            <Form.Label class="text-base">Password</Form.Label>
            <Input
                {...attrs}
                class="h-fit text-base"
                type="password"
                bind:value={$formData.password}
            />
        </Form.Control>
        <Form.FieldErrors />
    </Form.Field>
    <Form.Button class="h-fit w-full text-lg">Sign In</Form.Button>
</form>
