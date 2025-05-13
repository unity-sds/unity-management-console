<script lang="ts">
  import { page } from '$app/stores';
  import { config } from '../../store/stores';
  import type { NodeGroupType } from '../../data/entities';
  import { productInstall, marketplaceStore, isLoading, createEmptyMarketplaceMetadata } from '../../store/stores';
  import type { MarketplaceMetadata } from '../../store/stores';
  import SetupWizard from '../../components/SetupWizard.svelte';
  import AdvancedVar from './advanced_var.svelte';

  type StartApplicationInstallResponse = {
    deploymentID: string;
  };

  type ApplicationInstallStatus = { Status: string };

  // Get data from the load function
  export let data;

  let nodeGroups: NodeGroupType[] = [];
  let errorMessage = '';
  let deploymentID: string;
  
  // Create a stable reference for the empty product state
  const emptyProduct: MarketplaceMetadata = createEmptyMarketplaceMetadata();
  
  import { onMount } from 'svelte';
  
  // Variables to track initialization state
  let initialized = false;
  let initialProductSet = false;
  
  // Main reactive block for handling product loading and state
  $: if (initialized) {
    const appName = data.name;
    const appVersion = data.version;

    // --- Add Logging Start ---
    console.log('--- Reactive Block Start ---');
    console.log(`URL Params: name=${appName}, version=${appVersion}`);
    console.log(`Marketplace Store Length: ${$marketplaceStore?.length ?? 'undefined'}`);
    console.log(`Current $isLoading: ${$isLoading}`);
    console.log(`Current $productInstall Name: ${$productInstall?.Name ?? 'undefined'}`);
    console.log(`Current errorMessage: "${errorMessage}"`);
    // --- Add Logging End ---

    let targetProduct: MarketplaceMetadata = emptyProduct; // Default to empty
    let targetIsLoading: boolean = false;
    let targetErrorMessage: string = '';
    let stateChanged = false; // Flag to track if any state needs updating

    if (appName && appVersion) {
      // Parameters are present
      if ($marketplaceStore && $marketplaceStore.length > 0) {
        // Marketplace data is available. Find the product.
        const foundProduct = $marketplaceStore.find(
          (p) => p.Name === appName && p.Version === appVersion
        );

        if (foundProduct) {
          targetProduct = foundProduct;
          targetIsLoading = false;
          targetErrorMessage = '';
        } else {
          // Product not found
          targetProduct = emptyProduct;
          targetIsLoading = false;
          targetErrorMessage = `Product "${appName}" (Version: "${appVersion}") not found in the marketplace.`;
        }
      } else {
        // Marketplace data not yet available or store is empty. Loading phase.
        targetProduct = emptyProduct;
        targetIsLoading = true;
        targetErrorMessage = 'Marketplace data is loading...';
      }
    } else {
      // No application name/version in URL parameters.
      targetIsLoading = false;
      // Keep existing product if valid, otherwise set to empty and show message.
      if ($productInstall && $productInstall.Name !== '') {
         targetProduct = $productInstall;
         targetErrorMessage = '';
      } else {
         targetProduct = emptyProduct;
         targetErrorMessage = 'No application specified. Please select an application from the marketplace.';
      }
    }

    // --- Compare target state with current state and update only if changed ---

    // Compare isLoading
    if ($isLoading !== targetIsLoading) {
      // --- Add Logging ---
      console.log(`>>> Updating isLoading: ${$isLoading} -> ${targetIsLoading}`);
      // --- End Logging ---
      isLoading.set(targetIsLoading);
      stateChanged = true;
    }

    // Compare productInstall (more robust comparison)
    const currentProduct = $productInstall ?? emptyProduct;
    // Check if Name or Version differs
    if (!initialProductSet || (targetProduct && (currentProduct.Name !== targetProduct.Name || currentProduct.Version !== targetProduct.Version))) {
       // --- Add Logging ---
       console.log(`>>> Updating productInstall: ${currentProduct.Name}/${currentProduct.Version} -> ${targetProduct.Name}/${targetProduct.Version}`);
       // --- End Logging ---
       productInstall.set(targetProduct);
       initialProductSet = true;
       stateChanged = true;
    }

    // Compare errorMessage
    if (errorMessage !== targetErrorMessage) {
      // --- Add Logging ---
      console.log(`>>> Updating errorMessage: "${errorMessage}" -> "${targetErrorMessage}"`);
      // --- End Logging ---
      errorMessage = targetErrorMessage;
      // No store update here, just a local variable
    }

    // --- Add Logging State Change ---
    if (stateChanged) {
        console.log('--- State Updated This Cycle ---');
    } else {
        console.log('--- State Unchanged This Cycle ---');
    }
    console.log('--- Reactive Block End ---');
    // --- End Logging State Change ---
  }
  
  // Use onMount to initialize the component once
  onMount(() => {
    // Set initialized to true after the component is mounted
    initialized = true;
  });

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
  }

  $: managedDependenciesKeys =
    $productInstall && $productInstall.ManagedDependencies ? getObjectKeys($productInstall.ManagedDependencies) : [];

  const steps = ['deploymentDetails', 'variables', 'summary'];
  let currentStepIndex = 0;

  let applicationMetadata = {
    DeploymentName: '',
    Variables: {} as { [key: string]: string }
  };

  $: if (initialized && $config) {
    let configProjectChanged = false;
    let configVenueChanged = false;
    if ($config?.applicationConfig?.Project) {
      if (applicationMetadata.Variables.project !== $config.applicationConfig.Project) {
        applicationMetadata.Variables.project = $config.applicationConfig.Project;
        configProjectChanged = true;
      }
    }
    if ($config?.applicationConfig?.Venue) {
       if (applicationMetadata.Variables.venue !== $config.applicationConfig.Venue) {
        applicationMetadata.Variables.venue = $config.applicationConfig.Venue;
        configVenueChanged = true;
      }
    }
    // Log only if something actually changed in this block
    if (configProjectChanged || configVenueChanged) {
        console.log(`--- Config Reactive Block: Updated project=${applicationMetadata.Variables.project}, venue=${applicationMetadata.Variables.venue} ---`);
    }
  }

  $: Variables = $productInstall?.DefaultDeployment?.Variables?.Values || {};
  $: AdvancedValues = $productInstall?.DefaultDeployment?.Variables?.AdvancedValues || {};

  let varSetupDone = false;
  $: if (initialized && !varSetupDone && Object.keys(Variables).length) {
    console.log('--- Variable Setup Block Running ---'); // Add log
    Object.entries(Variables).forEach(([key, value]) => {
      if (value) {
        // Avoid overwriting if already set, maybe? Or just log
        // console.log(`Setting applicationMetadata.Variables[${key}] from Variables`);
        applicationMetadata.Variables[key] = value;
      }
    });
    varSetupDone = true;
    console.log('--- Variable Setup Block Done ---'); // Add log
  }

  let installInProgress = false;
  let installComplete = false;
  let installFailed = false;
  function startStatusPoller() {
    // Add a check to ensure $productInstall is valid before starting
    if (!$productInstall || !$productInstall.Name || !$productInstall.Version || !applicationMetadata.DeploymentName) {
        console.warn("Cannot start status poller: Missing product details or deployment name.");
        return;
    }
    let poller = setInterval(async (_) => {
      const res = await fetch(
        `../api/install_application/status/${$productInstall.Name}/${$productInstall.Version}/${applicationMetadata.DeploymentName}`
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
    // Add a check for valid product
    if (!$productInstall || !$productInstall.Name) {
        errMsg = "Cannot install: No product selected.";
        installFailed = true;
        return;
    }
    const outObj = {
      Name: $productInstall.Name,
      Version: $productInstall.Version,
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

  // --- Remove or comment out final state logs if too noisy, keep initial ones in the main block ---
  // $: console.log("Final productInstall state:", $productInstall);
  // $: console.log("Final isLoading state:", $isLoading);
  // $: console.log("Final errorMessage:", errorMessage);
  // $: console.log("Final config state:", $config);
</script>

<div class="container">
  {#if $isLoading}
    <div class="st-typography-header">Loading product information...</div>
  {:else if errorMessage}
    <div class="st-typography-header" style="color: red;">{errorMessage}</div>
  {:else if $productInstall?.Name}
    <div>
      <div class="st-typography-header">
        Installing Marketplace Application: <span class="st-typography-displayBody"
          >{$productInstall.DisplayName || $productInstall.Name}</span
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
          <div class="st-typography-bold">{$productInstall.Version}</div>
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
  {/if}
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
