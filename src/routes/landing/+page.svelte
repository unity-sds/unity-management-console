<script lang="ts">
  import { config, parametersStore } from "../../store/stores";
  import type { Config } from "../../data/unity-cs-manager/protobuf/extensions";
  import ControlPanelItem from "../../components/ControlPanelItem.svelte";

  let project: string;
  let venue: string;
  let conf: Config | null;

  config.subscribe(value => {
    conf = value;
  });

  parametersStore.subscribe(value => {
    project = value?.parameterlist?.["/unity/core/project"]?.value ?? "Unknown";
    venue = value?.parameterlist?.["/unity/core/venue"]?.value ?? "Unknown";
  });

  let setuprun: boolean;
  $: {
    setuprun = !!(conf && conf.updatedby !== "");
  }
  $: cardData = [
    {
      title: "Core Management",
      description: "Manage your core settings and features.",
      link: "/ui/setup",
      disabled: false
    },
    {
      title: "Unity Marketplace",
      description: "Explore the Unity Marketplace.",
      link: "/ui/marketplace",
      disabled: !setuprun
    },
    {
      title: "Application Management",
      description: "Manage your applications.",
      link: "/ui/applications",
      disabled: !setuprun
    }
    // {title: 'Extension Management', description: 'Manage your hosted extensions.', link: '#', disabled: setuprun}
  ];
</script>
<header class="bg-blue-600 text-white text-center py-12 mb-12">
  <h1 class="text-4xl">Welcome to the Unity Management Console</h1>
  <p class="text-lg leading-6">Control Panel</p>
</header>

<div class="container mx-auto">
  <div class="flex justify-center">
    <div class="flex-initial">
      {#if !setuprun}
        <div>
          <h5 class="text-xl">Setup has not been run, please go to Core Management</h5>
        </div>
      {:else}
        <ul class="list-decimal pl-5">
          <li class="bg-gray-200 p-4 border-b border-gray-300">Project: {project}</li>
          <li class="bg-gray-200 p-4">Venue: {venue}</li>
        </ul>
      {/if}
    </div>
  </div>

  <div class="flex justify-center mt-12 text-center">
    {#each cardData as card (card.title)}
      <ControlPanelItem title={card.title} description={card.description} link={card.link} disabled={card.disabled} />
    {/each}
  </div>
</div>