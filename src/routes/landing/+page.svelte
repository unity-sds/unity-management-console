<script lang="ts">
  import { config, parametersStore } from '../../store/stores';
  import type { Config } from '../../data/unity-cs-manager/protobuf/extensions';
  import ControlPanelItem from '../../components/ControlPanelItem.svelte';

  let conf: Config | null;

  config.subscribe((value) => {
    conf = value;
  });

  $: console.log({ conf });
  let setuprun: boolean;
  let bootstrapfailed: boolean;
  let bootstrapped: boolean;
  $: {
    if (conf && conf.bootstrap == 'complete') {
      bootstrapped = true;
      bootstrapfailed = false;
    } else if (conf && conf.bootstrap == 'failed') {
      bootstrapped = false;
      bootstrapfailed = true;
    } else if (conf && conf.bootstrap == '') {
      bootstrapped = false;
      bootstrapfailed = false;
    }
    setuprun = !!(conf && conf.updatedby !== '');
  }
  $: cardData = [
    {
      title: 'Unity Marketplace',
      description: 'Explore the Unity Marketplace.',
      link: '/management/ui/marketplace',
      disabled: !setuprun
    },
    {
      title: 'Application Management',
      description: 'Manage your applications.',
      link: '/management/ui/applications',
      disabled: !setuprun
    }
  ];

  $: project = $config ? $config.applicationConfig?.Project : '(loading)';
  $: venue = $config ? $config.applicationConfig?.Venue : '(loading)';
</script>

<header class="bg-blue-600 text-white text-center py-12 mb-12">
  <h1 class="text-4xl">Welcome to the Unity Management Console</h1>
  <p class="text-lg leading-6">Control Panel</p>
</header>

<div class="container mx-auto">
  <div class="flex justify-center">
    <div class="flex-initial">
      {#if bootstrapfailed}
        <div class="managementfeedback">
          <h5 class="text-xl">The Bootstrap Process Failed Please Check The Logs</h5>
        </div>
      {:else if !bootstrapped}
        <div class="managementfeedback">
          <h5 class="text-xl">Bootstrap is either in progress or has not been run</h5>
        </div>
      {:else}
        <div class="managementfeedback">
          <ul class="list-decimal pl-5">
            <li class="bg-gray-200 p-4 border-b border-gray-300">
              Project: {project}
            </li>
            <li class="bg-gray-200 p-4">Venue: {venue}</li>
          </ul>
        </div>
      {/if}
    </div>
  </div>

  <div class="flex justify-center mt-12 text-center">
    {#each cardData as card (card.title)}
      <ControlPanelItem
        title={card.title}
        description={card.description}
        link={card.link}
        disabled={card.disabled}
      />
    {/each}
  </div>
</div>
