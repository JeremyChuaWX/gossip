<script lang="ts">
    import type { PageData } from "./$types";
    import SigninForm from "./signin-form.svelte";

    export let data: PageData;
</script>

<div class="flex h-full w-full flex-col items-center gap-8 pt-40">
    <h1 class="text-4xl font-extrabold">Welcome back</h1>
    <SigninForm data={data.form} />
</div>
