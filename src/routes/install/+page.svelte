<script lang="ts">
  import { get } from 'svelte/store';
  import { config } from '../../store/stores';
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

  const steps = ['deploymentDetails', 'variables', 'summary'];
  let currentStepIndex = 0;

  let applicationMetadata = {
    DeploymentName: '',
    Variables: {} as { [key: string]: string }
  };

  $: {
    if ($config?.applicationConfig?.Project) {
      applicationMetadata.Variables.project = $config?.applicationConfig?.Project;
    }
    if ($config?.applicationConfig?.Venue) {
      applicationMetadata.Variables.venue = $config?.applicationConfig?.Venue;
    }
  }

  $: Variables = product?.DefaultDeployment?.Variables?.Values || {};
  $: {
    Object.entries(Variables).forEach(([key, value]) => {
      if (value) {
        applicationMetadata.Variables[key] = value;
      }
    });
  }

  async function installApplication() {
    const outObj = { name: product.Name, version: product.Version, ...applicationMetadata };
    console.log(outObj);
    const url =
      `http${window.location.protocol === 'https:' ? 's' : ''}://${window.location.host}` +
      '/management/install_application';
    const res = await fetch(url, { method: 'POST', body: JSON.stringify(outObj) });
    if (!res.ok) {
      console.log(res);
    }
  }

  $: console.log(product);
  $: console.log($config);
</script>

<div class="container">
  <div>
    <div class="st-typography-header">
      Installing Marketplace Application: <span class="st-typography-displayBody"
        >{product.DisplayName}</span
      >
    </div>
  </div>
  <hr />
  <div class="wizardContainer">
    {#if steps[currentStepIndex] === 'deploymentDetails'}
      <div class="st-typography-displayBody">Deployment Details</div>
      <div class="variablesForm">
        <div class="st-typography-label">Deployment Name</div>
        <input class="st-input" bind:value={applicationMetadata.DeploymentName} />
      </div>
    {:else if steps[currentStepIndex] === 'variables'}
      <div class="st-typography-small-caps">Variables</div>
      <div class="variablesForm">
        {#each Object.entries(Variables) as [key, value]}
          <div>
            <div class="st-typography-label">
              {key}
            </div>
            <input class="st-input" bind:value={applicationMetadata.Variables[key]} />
          </div>
        {/each}
      </div>
    {:else if steps[currentStepIndex] === 'summary'}
      <div class="st-typography-small-caps">Installation Summary</div>
      <div class="variablesForm">
        <div style="display: flex;">
          <div class="st-typography-label">Version</div>
          <div class="st-typography-bold">{product.Version}</div>
        </div>
        <hr />
        <div>
          <div class="st-typography-label">Variables</div>
          {#each Object.entries(applicationMetadata.Variables) as [key, value]}
            <div style="display: flex;">
              <div class="st-typography-label">{key}:&nbsp;</div>
              <div class="st-typography-bold">{value}</div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
    <hr />
    <div style="margin-top:10px;">
      {#if currentStepIndex > 0}
        <button class="st-button secondary" on:click={(_) => currentStepIndex--}>Back</button>
      {/if}
      {#if currentStepIndex === steps.length - 1}
        <button class="st-button" on:click={installApplication}>Install</button>
      {:else}
        <button class="st-button" on:click={(_) => currentStepIndex++}>Next</button>
      {/if}
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

<!-- <SetupWizard /> -->

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
    padding-top: 10px;
  }

  .variablesForm {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 10px;
  }

  .variablesForm input {
    width: 250px;
  }
</style>
