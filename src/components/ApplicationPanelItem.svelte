<script>
  import { uninstallApplication } from "../data/httpHandler";

  export let title = "";
  export let description = "";
  export let link = "";
  export let status = "";
  export let appPackage = "";
  export let deployment = "";
  const uninstallApp = () => {
    uninstallApplication(title, appPackage, deployment);
  };
</script>

<div class="lg:w-1/4 md:w-1/2 mb-4">
  <div class="bg-white border rounded shadow-md h-full">
    <div class="p-4">
      <h4 class="font-bold text-xl mb-2">{title}</h4>
      <p class="text-gray-700">{status} - {description}</p>
    </div>
    <div class="p-4 border-t">
      <a href="{link}" class="text-white bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded mr-2 inline-block">Explore</a>
      <a on:click={uninstallApp} class="text-white bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded inline-block">Uninstall</a>
    </div>
  </div>
</div>