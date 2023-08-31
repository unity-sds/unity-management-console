<script lang="ts">
  import { writable } from "svelte/store";
  import { config, parametersStore, projectStore, venueStore } from "../../store/stores";
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
    if (conf && conf.updatedby !== "") {
      setuprun = true;
    } else {
      setuprun = false;
    }
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
<header class="bg-primary text-white text-center py-5 mb-5">
  <h1>Welcome to the Unity Management Console</h1>
  <p class="lead">Control Panel</p>
</header>
<div class="container">
  <div class="row justify-content-md-center">
    <div class="col col-md-auto">
      {#if !setuprun}
        <div>
          <h5>Setup has not been run, please go to Core Management</h5>
        </div>
      {:else}
        <ul class="list-group">
          <li class="list-group-item">Project: {project}</li>
          <li class="list-group-item">Venue: {venue}</li>
        </ul>
      {/if}
    </div>
  </div>
  <div class="row text-center mt-5">
    {#each cardData as card (card.title)}
      <ControlPanelItem title={card.title} description={card.description} link={card.link} disabled={card.disabled} />
    {/each}
  </div>
</div>
