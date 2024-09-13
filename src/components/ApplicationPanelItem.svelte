<script lang="ts">
  import ScaleOut from './common/ScaleOut.svelte';
  import { HttpHandler, reapplyApplication } from '../data/httpHandler';
  import { goto } from '$app/navigation';

  import checkIcon from '../icons/check.svg';

  export let title = '';
  export let description = '';
  export let link = '';
  export let status = '';
  export let appPackage = '';
  export let appName = '';
  export let deployment = '';

  export let objectnumber = 0;

  let isUninstalling = false;
  const uninstallApp = () => {
    isUninstalling = true;
    // const httphandler = new HttpHandler();
    // console.log('Uninstalling ' + appName);
    // httphandler.uninstallSoftware(appName, appPackage, deployment);
  };

  const handleKeydown = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.key === objectnumber.toString()) {
      uninstallApp();
    } else if (event.key === objectnumber.toString()) {
      goto(link);
    }
  };

  const reapplyApp = () => {
    console.log('Reapplying: ' + title);
    reapplyApplication(title, appName, deployment);
  };
</script>

<div class="lg:w-1/3 md:w-1/2 mb-4">
  <div class="bg-white border rounded shadow-md h-full">
    <div style="display: flex; flex-direction: column; align-items: center; padding: 5px;">
      <span class="st-typography-header">{title}</span>
      <div style="display:flex; gap: 10px; margin: 10px; justify-content: center;">
        <span class="st-typography-bold">Installation Status:</span>
        {#if status === 'COMPLETE'}
          <span class="st-typography-small-caps" style="color: green;">Done</span>
        {:else}
          <span class="st-typography-small-caps" style="color:red;">{status}</span>
        {/if}
      </div>
    </div>
    <div class="p-4 border-t" style="text-align: center;">
      {#if isUninstalling}
        <div style="display: flex; gap: 5px;">
          <span class="st-typography-medium">Uninstalling....</span>
          <ScaleOut />
        </div>
      {:else}
        <a href={link} on:keydown={handleKeydown} class="st-button">Explore</a>
        <button on:click={reapplyApp} on:keydown={handleKeydown} class="st-button"
          >Reapply Installation
        </button>
        <button
          on:click={uninstallApp}
          on:keydown={handleKeydown}
          class="st-button tertiary"
          style="color: red; margin-top: 5px;"
          >Uninstall
        </button>
      {/if}
    </div>
  </div>
</div>
