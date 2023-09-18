<script lang="ts">
  import { writable } from "svelte/store";
  import { deploymentStore, projectStore, venueStore } from "../../store/stores";
  import ApplicationPanelItem from "../../components/ApplicationPanelItem.svelte";
  import { onDestroy, onMount } from "svelte";
  import { fetchDeployedApplications } from "../../data/httpHandler";

  const page = writable("");
  let project: string;
  let venue: string;
  projectStore.subscribe(value => {
    project = value;
  });
  venueStore.subscribe(value => {
    venue = value;
  });

  onMount(async () => {
    await fetchDeployedApplications();
  });

  type CardItem = {
    title: string;
    source: string;
    version: string;
    status: string;
    link: string;
    deploymentName: string;
  };

  let cardData: CardItem[] = [];

  $: {
    cardData;
  }

  const unsubscribe = deploymentStore.subscribe(value => {

    console.log(value);

    value?.deployment.forEach((el) => {
      const deploymentName = el.name;
      el.application.forEach(ar => {
        const newCardItem = {
          title: ar.applicationName,
          source: ar.source,
          version: ar.version,
          status: ar.status,
          link: "",
          deploymentName: deploymentName
        };
        cardData = [...cardData, newCardItem];

      });
    });
  });
  onDestroy(unsubscribe);
  let setuprun = true;
  $: {
    // If projectStore is not null or an empty string, set setuprun to false
    if ($projectStore && $projectStore.trim() !== "") {
      setuprun = false;
      console.log("store set");
    } else {
      setuprun = true;
      console.log("store not set");
    }
  }
  $: cardData = [];
</script>
<header class="bg-primary text-white text-center py-5 mb-5">
  <h1>Installed Applications</h1>
</header>
<div class="container">
  <div class="row text-center mt-5">
    {#each cardData as card, index (card.title)}
      <ApplicationPanelItem title={card.title} description={card.source} status={card.status} link={card.link}
                            appPackage={card.source} deployment={card.deploymentName} objectnumber={index + 1} />
    {/each}
  </div>
</div>
