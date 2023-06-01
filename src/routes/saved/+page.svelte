<script>
import { onMount } from 'svelte';

let message = 'Save complete';
let buttonText = 'Press okay to continue';

onMount(() => {
    // This will be run when the component first renders
});
</script>

<style>
.center-screen {
    height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}
</style>

<div class="center-screen text-center">
    <h1 class="mb-4">{message}</h1>
    <a href="/ui/landing" class="btn btn-primary">{buttonText}</a>
</div>
