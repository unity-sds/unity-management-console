<script lang="ts">
  import { HttpHandler, reapplyApplication } from '../data/httpHandler';
  import { goto } from '$app/navigation';

  export let title = '';
  export let description = '';
  export let link = '';
  export let status = '';
  export let appPackage = '';
  export let appName = '';
  export let deployment = '';

  export let objectnumber = 0;
  const uninstallApp = () => {
    const httphandler = new HttpHandler();
    console.log('Uninstalling ' + appName);
    httphandler.uninstallSoftware(appName, appPackage, deployment);
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

<div class="lg:w-1/4 md:w-1/2 mb-4">
  <div class="bg-white border rounded shadow-md h-full">
    <div style="display: flex; flex-direction: column;">
      <span class="st-typography-header">{title}</span>
      <span class="st-typography-bold">Installation Status: {status}</span>
    </div>
    <div class="p-4 border-t">
      <a
        href={link}
        on:keydown={handleKeydown}
        class="text-white bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded mr-2 inline-block"
        >Explore</a
      >
      <button
        on:click={reapplyApp}
        on:keydown={handleKeydown}
        class="text-white bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded inline-block"
        >Reapply Installation
      </button>
      <button
        on:click={uninstallApp}
        on:keydown={handleKeydown}
        class="text-white bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded inline-block"
        >Uninstall
      </button>
    </div>
  </div>
</div>
