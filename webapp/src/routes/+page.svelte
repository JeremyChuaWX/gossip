<script lang="ts">
    import { Button } from "$lib/components/ui/button";
</script>

<div class="mt-[-3rem] flex h-full w-full flex-col items-center justify-center gap-16">
    <div class="group flex flex-col items-center gap-4">
        <h1 class="text-9xl font-extrabold">Gossip</h1>
        <p class="animate-[widen-tracking_2s_ease-out_forwards] text-5xl">Speak your mind.</p>
    </div>
    <div class="flex w-fit gap-8">
        <Button href="/signup" variant="ghost" class="h-fit w-fit px-6 py-4 text-3xl"
            >Sign Up</Button
        >
        <Button href="/signin" variant="ghost" class="h-fit w-fit px-6 py-4 text-3xl"
            >Sign In</Button
        >
    </div>
</div>
