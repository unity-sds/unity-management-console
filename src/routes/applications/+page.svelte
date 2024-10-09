<script lang="ts">
  import { writable } from 'svelte/store';
  import { config, deploymentStore } from '../../store/stores';
  import ApplicationPanelItem from '../../components/ApplicationPanelItem.svelte';
  import { onDestroy, onMount } from 'svelte';
  import { fetchDeployedApplications } from '../../data/httpHandler';

  const page = writable('');

  let project = '';

  type InstalledMarketplaceApplication = {
    DeploymentName: string;
    PackageName: string;
    Name: string;
    Source: string;
    Version: string;
    Status: string;
  };

  let applications: InstalledMarketplaceApplication[] = [];

  async function getInstalledApplications() {
    const res = await fetch('../api/installed_applications');
    if (!res.ok) {
      console.warn('Unable to get application list!');
      return;
    }
    applications = await res.json();
    // json.forEach((app: InstalledMarketplaceApplication[] = []) => {
    //   const newCard: CardItem = {
    //     title: app.displayName,
    //     packageName: app.PackageName,
    //     applicationName: app.Name,
    //     source: app.Source,
    //     version: app.Version,
    //     status: app.Status,
    //     link: '',
    //     deploymentName: app.DisplayName
    //   };
    //   cardData = cardData.concat([newCard]);
    // });
    // console.log(json);
  }

  onMount(async () => {
    await getInstalledApplications();
    // await fetchDeployedApplications();
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
    if ($deploymentStore) {
      cardData = $deploymentStore.deployment.reduce<CardItem[]>((acc, el) => {
        const dplName = el.name;
        el.application.forEach((ar) => {
          const newCard: CardItem = {
            title: ar.displayName,
            source: ar.source,
            version: ar.version,
            status: ar.status,
            packageName: ar.packageName,
            link: '',
            deploymentName: dplName,
            applicationName: ar.displayName
          };
          acc.push(newCard);
        });
        return acc;
      }, []);
    }
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

<div style="margin-left: 20px">
  <div class="st-typography-displayH3">Installed Applications</div>
  <div
    style="
      width: 90%;
      display: flex;
      gap: 20px;
      margin-top: 10px;
      flex-wrap: wrap;
    "
  >
    {#each applications as card, index (card.DeploymentName)}
      <ApplicationPanelItem
        title={card.DeploymentName}
        description={card.Source}
        status={card.Status}
        appPackage={card.PackageName}
        appName={card.Name}
        deployment={card.DeploymentName}
        objectnumber={index + 1}
      />
    {/each}
  </div>
</div>
