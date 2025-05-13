<script lang="ts">
  import { get } from 'svelte/store';
  import {
    config,
    marketplaceStore,
    type MarketplaceMetadata,
    createEmptyMarketplaceMetadata,
    refreshConfig
  } from '../../store/stores';
  import type { NodeGroupType } from '../../data/entities';
  import SetupWizard from '../../components/SetupWizard.svelte';
  import AdvancedVar from './advanced_var.svelte';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  type StartApplicationInstallResponse = {
    deploymentID: string;
  };

  // Load data properties from the +page.ts load function
  export let data;

  // Make sure marketplace data is loaded, but only once
  let dataLoadInitiated = false;

  onMount(async () => {
    console.log('mount!');
    // Get URL parameters once on mount
    name = $page.url.searchParams.get('name') || '';
    version = $page.url.searchParams.get('version') || '';

    // Immediate validation of parameters
    if (name && version) {
      // We'll find the product after we ensure marketplace data is loaded
    } else {
      product = null;
      paramError = true;
      errorMessage = 'Missing required URL parameters: name and version';
    }

    // Only attempt to load marketplace data if we don't have it
    // and we haven't already started loading it
    if (!dataLoadInitiated && get(marketplaceStore).length === 0) {
      dataLoadInitiated = true;
      try {
        // Instead of creating a new HttpHandler, use our store function directly
        await refreshConfig();

        // Now that we have the marketplace data, find the product
        if (name && version) {
          const foundProduct = findProduct(name, version);
          if (foundProduct) {
            product = foundProduct;
            paramError = false;
            errorMessage = '';
          } else {
            product = null;
            paramError = true;
            errorMessage = `Could not find product "${name}" with version "${version}"`;
          }
        }
      } catch (error) {
        console.error('Error loading marketplace data:', error);
      }
    } else if (name && version) {
      // If marketplace data is already loaded, find the product
      const foundProduct = findProduct(name, version);
      if (foundProduct) {
        product = foundProduct;
        paramError = false;
        errorMessage = '';
      } else {
        product = null;
        paramError = true;
        errorMessage = `Could not find product "${name}" with version "${version}"`;
      }
    }
  });

  type ApplicationInstallStatus = { Status: string };

  let nodeGroups: NodeGroupType[] = [];

  // Initialize with an empty product
  let product: MarketplaceMetadata | null = null;

  // Get name and version from URL parameters once, not reactively
  let name = '';
  let version = '';

  // We'll initialize these in onMount instead of using reactive declarations

  // Track if we have valid parameters
  let paramError = false;
  let errorMessage = '';

  // Function to find a product by name and version
  function findProduct(name: string, version: string): MarketplaceMetadata | undefined {
    return get(marketplaceStore).find((p) => p.Name === name && p.Version === version);
  }

  // We've removed the reactive product lookup block
  // and now handle this once in onMount to prevent refresh loops

  let deploymentID: string;

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
  }

  // Calculate this only when product changes (non-reactive)
  let managedDependenciesKeys: string[] = [];

  $: if (product && product.ManagedDependencies) {
    managedDependenciesKeys = getObjectKeys(product.ManagedDependencies);
  } else {
    managedDependenciesKeys = [];
  }

  const steps = ['deploymentDetails', 'variables', 'dependencies', 'summary'];
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

  // Use empty objects as fallbacks when product is null
  $: Variables = product ? product.DefaultDeployment?.Variables?.Values || {} : {};
  $: AdvancedValues = product ? product.DefaultDeployment?.Variables?.AdvancedValues || {} : {};

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
    // Only start polling if product is not null
    if (!product) return;

    // Capture product name and version to use in the interval
    const productName = product.Name;
    const productVersion = product.Version;

    let poller = setInterval(async (_) => {
      const res = await fetch(
        `../api/install_application/status/${productName}/${productVersion}/${applicationMetadata.DeploymentName}`
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
    // Only proceed if product is not null
    if (!product) {
      errMsg = 'Cannot install: product information is missing';
      installFailed = true;
      return;
    }

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
    // Only proceed if product is not null
    if (!product) {
      console.warn('Cannot get logs: product information is missing');
      return;
    }

    // Capture product name to use in the fetch call
    const productName = product.Name;

    const res = await fetch(
      `../api/install_application/logs/${productName}/${applicationMetadata.DeploymentName}`
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

  async function getDependencyCheck() {
    if (!product) return {};
    const res = await fetch(
      `../api/check_application_dependencies/${product.Name}/${product.Version}`
    );
    if (!res.ok) return;
    return await res.json();
  }

  // Commented out console logs that were causing refresh loops
  // $: console.log(product);
  // $: console.log($config);
</script>

<div class="container">
  {#if paramError}
    <div class="error-container">
      <div class="st-typography-header" style="color: red;">Error</div>
      <div class="st-typography-body">{errorMessage}</div>
      <div class="st-typography-body" style="margin-top: 1rem;">
        <a href="/management/ui/marketplace" class="st-button">Return to Marketplace</a>
      </div>
    </div>
  {:else if product}
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
      {:else if steps[currentStepIndex] === 'dependencies'}
        <div class="st-typography-displayBody">Dependencies</div>
        {#if !Object.keys(product.Dependencies).length}
          <span class="st-typography-label">This product has no dependencies.</span>
        {:else}
          {#await getDependencyCheck()}
            <span class="st-typography-label">Checking dependencies...</span>
          {:then depInfo}
            <span class="st-typography-label">This product does have dependencies</span>
          {/await}
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

  .error-container {
    padding: 2rem;
    border: 1px solid #f56565;
    border-radius: 0.5rem;
    background-color: #fff5f5;
    text-align: center;
    max-width: 600px;
    margin: 0 auto;
    margin-top: 2rem;
  }
</style>
