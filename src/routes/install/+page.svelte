<script lang="ts">
  import { get } from 'svelte/store';
  import { config } from '../../store/stores';
  import type { NodeGroupType } from '../../data/entities';
  import { productInstall } from '../../store/stores';
  import SetupWizard from '../../components/SetupWizard.svelte';

  type StartApplicationInstallResponse = {
    deploymentID: string;
  };

  type ApplicationInstallStatus = [{ Status: string }];

  let nodeGroups: NodeGroupType[] = [];

  let product = get(productInstall);

  let deploymentID: string;

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

  let installInProgress = false;
  let installComplete = false;
  function startStatusPoller() {
    let poller = setInterval(async (_) => {
      const res = await fetch(`../api/install_application/status/${deploymentID}`);
      if (!res.ok) {
        console.warn("Couldn't get status!");
        return;
      }
      const json = (await res.json()) as ApplicationInstallStatus;
      if (json[0].Status === 'COMPLETE') {
        clearInterval(poller);
        installInProgress = false;
        installComplete = true;
      }
    }, 5000);
  }

  async function installApplication() {
    const outObj = { Name: product.Name, Version: product.Version, ...applicationMetadata };
    installInProgress = true;
    const url = '../api/install_application';
    const res = await fetch(url, { method: 'POST', body: JSON.stringify(outObj) });
    if (!res.ok) {
      console.log(res);
      return;
    }
    const json = (await res.json()) as StartApplicationInstallResponse;
    if (json.deploymentID) {
      deploymentID = json.deploymentID;
      startStatusPoller();
    } else {
      console.error('No deploymentID received in the response');
    }
  }

  let showLogs = false;
  let logs: string;
  async function toggleLogs() {
    showLogs = !showLogs;
    if (showLogs && !logs) {
      const res = await fetch(`../api/install_application/logs/${deploymentID}`);
      if (!res.ok) {
        console.warn('Unable to get logs!');
        return;
      }
      logs = await res.text();
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
          <div class="st-typography-label">Version:&nbsp;</div>
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
      {#if installInProgress}
        <button class="st-button" disabled on:click={installApplication}>Installing...</button>
      {:else if installComplete}
        <button class="st-button" disabled on:click={installApplication}>Install Complete</button>
      {:else if currentStepIndex === steps.length - 1}
        <button class="st-button" on:click={installApplication}>Install</button>
      {:else}
        <button class="st-button" on:click={(_) => currentStepIndex++}>Next</button>
      {/if}

      {#if installInProgress || installComplete}
        <button class="st-button" on:click={toggleLogs}>{showLogs ? 'Hide' : 'Show'} Logs</button>
      {/if}
    </div>
    {#if showLogs}
      <div style="margin-top:10px">
        <hr />
        <pre>
      {logs}
    </pre>
      </div>
    {/if}
  </div>
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
