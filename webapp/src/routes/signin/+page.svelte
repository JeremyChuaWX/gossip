<script lang="ts">
    import type { PageData } from "./$types";
    import SigninForm from "./signin-form.svelte";

    export let data: PageData;
</script>

<div class="container mt-[-3rem] flex h-full w-full flex-col items-center justify-center gap-8">
    <h1 class="text-5xl font-extrabold">Gossip</h1>
    <SigninForm data={data.form} />
</div>
