<script lang="ts">
    import type { PageData } from "./$types";

    export let data: PageData;
</script>

<div>
    <h1>Home Page</h1>
    <p>Welcome {data.user.username}!</p>
</div>
