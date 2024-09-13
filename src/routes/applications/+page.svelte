<script lang="ts">
  import { writable } from 'svelte/store';
  import { config, deploymentStore } from '../../store/stores';
  import ApplicationPanelItem from '../../components/ApplicationPanelItem.svelte';
  import { onDestroy, onMount } from 'svelte';
  import { fetchDeployedApplications } from '../../data/httpHandler';

  const page = writable('');

  let project = '';

  onMount(async () => {
    await fetchDeployedApplications();
  });

  type CardItem = {
    title: string;
    packageName: string;
    applicationName: string;
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

  const unsubscribe = deploymentStore.subscribe((value) => {
    console.log(value);

    value?.deployment.forEach((el) => {
      const dplName = el.name;
      el.application.forEach((ar) => {
        const newCardItem = {
          title: ar.displayName,
          source: ar.source,
          version: ar.version,
          status: ar.status,
          packageName: ar.packageName,
          link: '',
          deploymentName: dplName,
          applicationName: ar.displayName
        };
        cardData = [...cardData, newCardItem];
      });
    });
  });
  onDestroy(unsubscribe);
  let setuprun = true;
  $: {
    // If projectStore is not null or an empty string, set setuprun to false
    if ($config && $config?.applicationConfig?.Project) {
      setuprun = false;
      console.log('store set');
    } else {
      setuprun = true;
      console.log('store not set');
    }
  }
  $: cardData = [];
</script>

<h1>Installed Applications</h1>

<div class="container">
  <div class="row text-center mt-5">
    {#each cardData as card, index (card.title)}
      <ApplicationPanelItem
        title={card.title}
        description={card.source}
        status={card.status}
        link={card.link}
        appPackage={card.packageName}
        appName={card.applicationName}
        deployment={card.deploymentName}
        objectnumber={index + 1}
      />
    {/each}
  </div>
</div>
