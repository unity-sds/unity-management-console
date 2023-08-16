<script lang="ts">
    import { installError, installRunning, messageStore } from "../store/stores";
    import { tick } from 'svelte';

    let installRunningValue : boolean;
    const unsubscribeInstallRunning = installRunning.subscribe(value => {
        installRunningValue = value;
    });

    let installErrorValue : boolean;
    const unsubscribeErrorRunning = installError.subscribe(value => {
        installErrorValue = value;
    });
    let textarea : HTMLTextAreaElement;

    const unsubscribe = messageStore.subscribe((value) => {
        // This code runs whenever messageStore changes
        if (textarea) {
            // Use nextTick to ensure that the DOM has been updated
            tick().then(() => {
                textarea.scrollTop = textarea.scrollHeight;
            });
        }
    });

    // Clean up the subscription when the component is destroyed
    import { onDestroy } from 'svelte';
    onDestroy(() => {
        unsubscribe();
        unsubscribeInstallRunning();
        unsubscribeErrorRunning();
    });
</script>

<div class="container d-flex align-items-center justify-content-center vh-100">
    <div class="text-center w-100">
        <div class="row text-center">
            <h2>Installing Application Stack</h2>
        </div>
        <div class="row">
            <div class="form-group col-md-12">
                <textarea bind:this={textarea} class="form-control" id="console" rows="30" bind:value={$messageStore} readonly></textarea>
            </div>
        </div>
        <div class="row mt-3">
            <div class="col-md-3"></div>
            <div class="col-md-3"></div>
            {#if !installRunningValue}
                {#if !installErrorValue}
              	  <a href="/ui/landing" class="btn btn-primary mt-3">Installation Complete</a>
                {:else}
                    <a href="/ui/landing" class="btn btn-primary mt-3">Installation Failed!</a>
                {/if}
            {/if}
        </div>
    </div>
</div>
