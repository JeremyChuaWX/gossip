<script lang="ts">
    import type { PageData } from "./$types";
    import SignupForm from "./signup-form.svelte";

    export let data: PageData;
</script>

<div class="flex h-full w-full flex-col items-center gap-8 pt-40">
    <h1 class="text-5xl font-extrabold">Gossip</h1>
    <SignupForm data={data.form} />
</div>
