<script lang="ts">
  import { get } from 'svelte/store';
  import type { NodeGroupType } from '../../data/entities';
  import { productInstall } from '../../store/stores';
  import SetupWizard from '../../components/SetupWizard.svelte';

  let nodeGroups: NodeGroupType[] = [];

  let product = get(productInstall);

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
  }

  $: managedDependenciesKeys =
    product && product.ManagedDependencies ? getObjectKeys(product.ManagedDependencies) : [];

  const steps = ['deploymentDetails', 'variables'];
  let currentStepIndex = 0;

  let applicationMetadata = {
    deploymentName: ''
  };

  $: baseVariables = product?.DefaultDeployment?.Variables?.Values || {};

  $: console.log(product);
</script>

<div class="container">
  <div class="st-typography-header">Installing Marketplace Application:</div>
  {product.DisplayName}
  <hr />
  <div class="wizardContainer">
    {#if steps[currentStepIndex] === 'deploymentDetails'}
      <div class="st-typography-displayBody">Deployment Details</div>
      <div>
        <div class="st-typography-label">Deployment Name</div>
        <input class="st-input" bind:value={applicationMetadata.deploymentName} />
      </div>
    {:else if steps[currentStepIndex] === 'variables'}
      <div class="st-typography-displayBody">Variables</div>
      {#each Object.entries(baseVariables) as [key, value]}
        <div>{key}</div>
      {/each}
    {/if}
    <div>
      {#if currentStepIndex > 0}
        <button class="st-button" on:click={(_) => currentStepIndex--}>Back</button>
      {/if}
      <button class="st-button" on:click={(_) => currentStepIndex++}>Next</button>
    </div>
  </div>
  <!--   <div class="row">
    <div class="col-md-12">
      {#if product}

        <h1 class="my-4">{product.DisplayName} Installation</h1> -->

  <!--       {:else}
        <p>Loading product...</p>
      {/if}
    </div>
 -->
  <!-- </div> -->
</div>
<SetupWizard />

<style>
  .container {
    display: flex;
    margin: 10px;
    flex-direction: column;
    margin-left: 30px;
  }

  .wizardContainer {
    display: flex;
    flex-direction: column;
  }
</style>
