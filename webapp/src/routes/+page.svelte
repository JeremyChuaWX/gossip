<script>
    import { Button } from "$lib/components/ui/button";
</script>

<div class="mt-[-3rem] flex h-full w-full flex-col items-center justify-center gap-10">
    <div class="group flex flex-col items-center gap-4">
        <h1 class="text-9xl font-extrabold">
            Gossip
            <span class="block h-2 w-0 bg-white transition-all duration-500 group-hover:w-full"
            ></span>
        </h1>
        <p
            class="text-5xl transition-all duration-500 group-hover:font-extralight group-hover:tracking-widest"
        >
            Speak your mind.
        </p>
    </div>
    <div class="flex w-fit gap-4">
        <Button href="/signup" variant="ghost" class="h-fit w-fit px-6 py-4 text-3xl"
            >Sign Up</Button
        >
        <Button href="/signin" variant="ghost" class="h-fit w-fit px-6 py-4 text-3xl"
            >Sign In</Button
        >
    </div>
</div>
