<script lang="ts">
  import { get } from 'svelte/store';
  import { config } from '../../store/stores';
  import type { NodeGroupType } from '../../data/entities';
  import { productInstall } from '../../store/stores';
  import SetupWizard from '../../components/SetupWizard.svelte';
  import AdvancedVar from './advanced_var.svelte';

  type StartApplicationInstallResponse = {
    deploymentID: string;
  };

  type ApplicationInstallStatus = { Status: string };

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
  $: AdvancedValues = product?.DefaultDeployment?.Variables?.AdvancedValues || {};

  let varSetupDone = false;
  $: {
    if (!varSetupDone && Object.keys(Variables).length) {
      Object.entries(Variables).forEach(([key, value]) => {
        if (value) {
          applicationMetadata.Variables[key] = value;
        }
      });
      varSetupDone = true;
    }
  }

  let installInProgress = false;
  let installComplete = false;
  let installFailed = false;
  function startStatusPoller() {
    let poller = setInterval(async (_) => {
      const res = await fetch(
        `../api/install_application/status/${product.Name}/${product.Version}/${applicationMetadata.DeploymentName}`
      );
      if (!res.ok) {
        console.warn("Couldn't get status!");
        return;
      }
      const json = (await res.json()) as ApplicationInstallStatus;
      if (json.Status === 'COMPLETE') {
        clearInterval(poller);
        installInProgress = false;
        installComplete = true;
      } else if (json.Status.includes('FAILED')) {
        installInProgress = false;
        installComplete = true;
        installFailed = true;
      }
    }, 5000);
  }

  let errMsg = '';
  async function installApplication() {
    const outObj = {
      Name: product.Name,
      Version: product.Version,
      AdvancedValues,
      ...applicationMetadata
    };
    installInProgress = true;
    const url = '../api/install_application';
    const res = await fetch(url, { method: 'POST', body: JSON.stringify(outObj) });
    if (!res.ok) {
      try {
        const json = await res.json();
        if (json.error) {
          errMsg = json.error;
          installFailed = true;
          installInProgress = false;
        }
      } catch (e) {}
      return;
    }
    startStatusPoller();
  }

  let showLogs = false;
  let logs: string;
  let logsDiv: HTMLElement;

  function scrollLogsToBottom() {
    if (logsDiv) {
      logsDiv.scrollIntoView({ behavior: 'smooth', block: 'end' });
    }
  }

  $: if (logs) {
    // Wait for DOM update
    setTimeout(scrollLogsToBottom, 0);
  }

  let logInterval: any = null;

  async function getLogs() {
    const res = await fetch(
      `../api/install_application/logs/${product.Name}/${applicationMetadata.DeploymentName}`
    );
    if (!res.ok) {
      console.warn('Unable to get logs!');
      return;
    }
    const logStr = await res.text();
    if (logStr) {
      logs = logStr;
    }
  }

  $: {
    if (showLogs && !logInterval) {
      getLogs();
      logInterval = setInterval((_) => {
        getLogs();
      }, 5000);
    } else if (!showLogs) {
      clearInterval(logInterval);
      logInterval = null;
    }
  }

  let validationErrors = { deploymentName: '', variables: {} as { [key: string]: string } };
  function gotoNextStep() {
    let hasErrors = false;
    switch (steps[currentStepIndex]) {
      case 'deploymentDetails':
        if (!applicationMetadata.DeploymentName) {
          validationErrors.deploymentName = 'Please enter a deployment name.';
          hasErrors = true;
        } else {
          validationErrors.deploymentName = '';
        }
        break;
      case 'variables':
        Object.keys(Variables).forEach((key) => {
          if (!applicationMetadata.Variables[key]) {
            validationErrors.variables[key] = "This value can't be blank.";
            hasErrors = true;
          } else {
            validationErrors.variables[key] = '';
          }
        });
        break;
    }
    if (hasErrors) return;
    currentStepIndex = currentStepIndex + 1;
  }

  $: console.log(product);
  $: console.log($config);
</script>

<div class="container">
  <div>
    <div class="st-typography-header">
      Installing Marketplace Application: <span class="st-typography-displayBody"
        >{product.Name}</span
      >
    </div>
  </div>
  <hr />
  <div class="wizardContainer">
    {#if steps[currentStepIndex] === 'deploymentDetails'}
      <div class="st-typography-displayBody">Deployment Details</div>
      <div class="variablesForm">
        <div class="st-typography-label">
          Deployment Name (this should be a unique identifier for this installation of the
          Marketplace item)
        </div>
        {#if validationErrors.deploymentName}
          <span class="st-typography-label" style="color:red;"
            >{validationErrors.deploymentName}</span
          >
        {/if}
        <input class="st-input" bind:value={applicationMetadata.DeploymentName} maxlength="32" />
      </div>
    {:else if steps[currentStepIndex] === 'variables'}
      <div class="st-typography-small-caps">Variables</div>
      <div class="variablesForm">
        {#each Object.entries(Variables) as [key, value]}
          <div>
            <div class="st-typography-label">
              {key}
            </div>
            <div style="display: flex; flex-direction: column;">
              {#if validationErrors.variables[key]}
                <span class="st-typography-label" style="color:red;"
                  >{validationErrors.variables[key]}</span
                >
              {/if}
              <input class="st-input" bind:value={applicationMetadata.Variables[key]} />
            </div>
          </div>
        {/each}
      </div>
      {#if AdvancedValues}
        <hr style="margin-top:10px" />
        <div class="variablesForm">
          <AdvancedVar bind:json={AdvancedValues} editMode={true} />
        </div>
      {/if}
    {:else if steps[currentStepIndex] === 'summary'}
      <div class="st-typography-small-caps">Installation Summary</div>
      <div class="variablesForm">
        <div style="display: flex;">
          <div class="st-typography-label">Version:&nbsp;</div>
          <div class="st-typography-bold">{product.Version}</div>
        </div>
        <hr />
        <div>
          <div class="st-typography-bold">Variables</div>
          {#each Object.entries(applicationMetadata.Variables) as [key, value]}
            <div style="display: flex;">
              <div class="st-typography-label">{key}:&nbsp;</div>
              <div class="st-typography-bold">{value}</div>
            </div>
          {/each}
          {#if AdvancedValues}
            <AdvancedVar bind:json={AdvancedValues} />
          {/if}
        </div>
      </div>
    {/if}
    <hr style="margin-top:10px" />
    <div style="margin-top:10px;">
      {#if currentStepIndex > 0}
        <button class="st-button secondary" on:click={(_) => currentStepIndex--}>Back</button>
      {/if}
      {#if installInProgress}
        <button class="st-button" disabled>Installing...</button>
      {:else if installFailed}
        <button class="st-button" disabled style="color:red;">Install Failed!</button>
      {:else if installComplete}
        <button class="st-button" disabled on:click={installApplication}>Install Complete</button>
      {:else if currentStepIndex === steps.length - 1}
        <button class="st-button" on:click={installApplication}>Install</button>
      {:else}
        <button class="st-button" on:click={gotoNextStep}>Next</button>
      {/if}

      {#if installInProgress || installComplete}
        <button class="st-button" on:click={(_) => (showLogs = !showLogs)}
          >{showLogs ? 'Hide' : 'Show'} Logs</button
        >
      {/if}
    </div>
    {#if errMsg}
      <div class="st-typography-label" style="color:red;">{errMsg}</div>
    {/if}
    {#if showLogs && logs}
      <div style="margin-top:10px">
        <hr />
        <pre bind:this={logsDiv}>
      {logs}
    </pre>
      </div>
    {/if}
  </div>
</div>

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
