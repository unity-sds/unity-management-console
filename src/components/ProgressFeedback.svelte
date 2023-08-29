<script lang="ts">
    import { installError, installRunning } from "../store/stores";

    let installRunningValue : boolean;
    const unsubscribeInstallRunning = installRunning.subscribe(value => {
        installRunningValue = value;
    });

    let installErrorValue : boolean;
    const unsubscribeErrorRunning = installError.subscribe(value => {
        installErrorValue = value;
    });


    // Clean up the subscription when the component is destroyed
    import { onDestroy } from 'svelte';
    import SocketTerminal from "./SocketTerminal.svelte";
    onDestroy(() => {
        unsubscribeInstallRunning();
        unsubscribeErrorRunning();
    });
</script>

<div class="container d-flex vh-100">
    <div class="w-100">
        <div class="row text-center">
            <h2>Installing Application Stack</h2>
        </div>
        <div class="row">
            <div class="form-group col-md-12">
                <SocketTerminal/>
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
