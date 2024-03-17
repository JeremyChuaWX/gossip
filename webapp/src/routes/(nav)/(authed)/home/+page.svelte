<script lang="ts">
    import type { PageData } from "./$types";

    export let data: PageData;
</script>

<div class="mx-auto flex h-full w-[min(80%,1000px)]">
    <div class="w-1/3 border-r border-zinc-800 p-4">
        <div class="flex justify-between">
            <h1>Rooms</h1>
            <button>+</button>
        </div>
    </div>
    <div class="flex h-full w-full flex-col border-r border-zinc-800">
        <div class="h-full w-full p-4">
            <p>Welcome {data.user.username}!</p>
        </div>
        <div class="flex h-12 w-full border-t border-zinc-800">
            <input class="h-full w-full bg-transparent p-4" />
        </div>
    </div>
</div>
