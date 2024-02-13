<script lang="ts">
  import { installError, installRunning } from "../store/stores";

  let installRunningValue: boolean;
  const unsubscribeInstallRunning = installRunning.subscribe(value => {
    installRunningValue = value;
  });

  let installErrorValue: boolean;
  const unsubscribeErrorRunning = installError.subscribe(value => {
    installErrorValue = value;
  });


  // Clean up the subscription when the component is destroyed
  import { onDestroy } from "svelte";
  import SocketTerminal from "./SocketTerminal.svelte";

  onDestroy(() => {
    unsubscribeInstallRunning();
    unsubscribeErrorRunning();
  });
</script>

<div class="flex h-screen items-center">
  <div class="w-full">
    <div class="text-center">
      <h2>Installing Application Stack</h2>
    </div>
    <div class="flex">
      <div class="flex-grow w-full md:w-1/1">
        <SocketTerminal />
      </div>
    </div>
    <div class="flex mt-3">
      <div class="w-1/4 md:w-1/4"></div>
      <div class="w-1/4 md:w-1/4"></div>
      {#if !installRunningValue}
        {#if !installErrorValue}
          <a href="/management/ui/landing"
             class="btn btn-primary mt-3 bg-blue-500 text-white px-4 py-2 rounded installfinished">Installation
            Complete</a>
        {:else}
          <a href="/management/ui/landing"
             class="btn btn-primary mt-3 bg-red-500 text-white px-4 py-2 rounded installfinished">Installation
            Failed!</a>
        {/if}
      {/if}
    </div>
  </div>
</div>