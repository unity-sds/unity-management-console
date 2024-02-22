<script lang="ts">
  import ProductForm from "./ProductForm.svelte";
  import { deploymentStore, installError, productInstall, installRunning } from "../store/stores";
  import VariablesForm from "./VariablesForm.svelte";
  import { fetchDeployedApplications, HttpHandler } from "../data/httpHandler";
  import { Deployments, Install_Applications, Install_Variables } from "../data/unity-cs-manager/protobuf/extensions";
  import { goto } from "$app/navigation";
  import InstallSummary from "./InstallSummary.svelte";
  import { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";
  import { onMount } from "svelte";
  import Deployment from "./Deployment.svelte";

  onMount(async () => {
    await fetchDeployedApplications();
  });
  let dependencyMap: { [key: string]: string } = {};

  let product: MarketplaceMetadata = MarketplaceMetadata.create();

  productInstall.subscribe(value => {
    product = value;
  });

  let step1Class = "";
  let step2Class = "";
  let step3Class = "";
  let step4Class = "";
  let step5Class = "";


  $: {
    step1Class = getStepStatus(1, currentStep);
    step2Class = getStepStatus(2, currentStep);
    step3Class = getStepStatus(3, currentStep);
    step4Class = getStepStatus(4, currentStep);
    step5Class = getStepStatus(5, currentStep);
  }

  let currentStep = 1;

  function nextStep() {
    currentStep++;
  }

  function prevStep() {
    currentStep--;
  }

  function getStepStatus(stepNumber: number, cstep: number) {
    if (cstep === stepNumber) {
      return "flex flex-col border-l-4 border-indigo-600 py-2 pl-4 md:border-l-0 md:border-t-4 md:pb-0 md:pl-0 md:pt-4";
    } else if (cstep > stepNumber) {
      return "flex flex-col border-l-4 border-indigo-600 py-2 pl-4 md:border-l-0 md:border-t-4 md:pb-0 md:pl-0 md:pt-4";
    }
    return "group flex flex-col border-l-4 border-gray-200 py-2 pl-4 hover:border-gray-300 md:border-l-0 md:border-t-4 md:pb-0 md:pl-0 md:pt-4";
  }

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
  }

  function resetInstallValues() {
    installError.set(false);
    installRunning.set(true);
  }

  let installName = "";
  const installSoftware = async () => {
    if (!product) {
      console.error("No product selected for installation");
      return;
    } else {
      console.log(product);
    }

    const httpHandler = new HttpHandler();

    type AnyObject = { [key: string]: any };

    const removeEmptyStrings = (obj?: AnyObject): void => {
      if (!obj) return;

      for (const key in obj) {
        if (typeof obj[key] === "string" && obj[key].length === 0) {
          delete obj[key];
        } else if (typeof obj[key] === "object" && obj[key] !== null) {
          // Recursively clean nested objects
          removeEmptyStrings(obj[key]);
        }
      }
    };

    type MergedType = {
      Values?: { [key: string]: string };
      AdvancedValues?: { [key: string]: any };
    };

    const merged: MergedType = {};

    removeEmptyStrings(product.DefaultDeployment?.Variables?.Values);
    removeEmptyStrings(product.DefaultDeployment?.Variables?.AdvancedValues);
    merged.Values = product.DefaultDeployment?.Variables?.Values;
    merged.AdvancedValues = product.DefaultDeployment?.Variables?.AdvancedValues;


    const vars = Install_Variables.fromJSON(merged);
    const a = Install_Applications.create({
      name: product.Name,
      version: product.Version,
      variables: vars,
      displayname: product.DisplayName,
      dependencies: dependencyMap
    } as any);
    console.log("installing");
    console.log(a);
    resetInstallValues();
    const id = await httpHandler.installSoftware(a, installName);
    console.log(id);
    goto("/management/ui/progress", { replaceState: true });
  };

  $: managedDependenciesKeys = product && product.ManagedDependencies ? getObjectKeys(product.ManagedDependencies) : [];

  let deployed: Deployments;

  deploymentStore.subscribe(value => {
    deployed = value;
  });

  function getVersionsForKey(key: string): string[] {
    let options: string[] = [];
    if (deployed && deployed.deployment) {
      for (let d of deployed.deployment) {
        for (let a of d.application) {
          if (key === a.packageName) {
            options.push(a.displayName);
          }
        }
      }
    }
    return options;
  }

  // Add a new function to handle the select change
  function handleDependencyChange(event: Event, key: string) {
    const selectElement = event.target as HTMLSelectElement;
    const selectedValue = selectElement.value;

    console.log("adding dependency: " + key + " " + selectedValue);
    if (selectedValue) {
      dependencyMap[key] = selectedValue;
    } else {
      delete dependencyMap[key];
    }
  }

  let i = 0;
</script>

<section class="py-8">

  <div class="mx-auto max-w-7xl px-4 py-24 sm:px-6 sm:py-32 lg:px-8">
    <div class="mx-auto max-w-2xl">
      <nav aria-label="Progress">
        <ol class="space-y-4 md:flex md:space-x-8 md:space-y-0">
          <li class="md:flex-1">
            <!-- Completed Step -->
            <a href="#"
               class="{step1Class}"
               aria-current="step">
              <span class="text-sm font-medium text-indigo-600">Step 1</span>
              <span class="text-sm font-medium">Deployment Details</span>
            </a>
          </li>
          <li class="md:flex-1">
            <!-- Completed Step -->
            <a href="#"
               class="{step2Class}"
               aria-current="step">
              <span class="text-sm font-medium text-indigo-600">Step 2</span>
              <span class="text-sm font-medium">Application Details</span>
            </a>
          </li>
          <li class="md:flex-1">
            <!-- Current Step -->
            <a href="#"
               class="{step3Class}"
               aria-current="step">
              <span class="text-sm font-medium text-indigo-600">Step 3</span>
              <span class="text-sm font-medium">Dependencies</span>
            </a>
          </li>
          <li class="md:flex-1">
            <!-- Upcoming Step -->
            <a href="#" class="{step4Class}">
              <span class="text-sm font-medium text-gray-500 group-hover:text-gray-700">Step 4</span>
              <span class="text-sm font-medium">Variables</span>
            </a>
          </li>
          <li class="md:flex-1">
            <!-- Upcoming Step -->
            <a href="#" class="{step5Class}">
              <span class="text-sm font-medium text-gray-500 group-hover:text-gray-700">Step 5</span>
              <span class="text-sm font-medium">Summary</span>
            </a>
          </li>
        </ol>
      </nav>

      <form on:submit|preventDefault={installSoftware}>
        {#if currentStep === 1}
          <Deployment bind:installName={installName} />
        {/if}
        {#if currentStep === 2}
          <ProductForm bind:product />
        {/if}
        {#if currentStep === 3}
          <!-- Collapse 2 content -->
          <div class="list-content">
            {#if product.ManagedDependencies}
              <h2>Dependencies</h2>
              <!--{#each product.ManagedDependencies as dependency}-->
              {#each managedDependenciesKeys as key}
                <div class="form-group">
                  <label class="col-form-label">
                    {key}
                    <select class="form-control" on:change="{(e) => handleDependencyChange(e, key)}">
                      <option></option>
                      {#each getVersionsForKey(key) as version}
                        <option>{version}</option>
                      {/each}
                    </select>
                  </label>
                </div>
              {/each}
            {/if}
          </div>
        {/if}

        {#if currentStep === 4}
          <!-- Collapse 3 content -->
          <div class="list-content">
            <VariablesForm bind:product />
          </div>
        {/if}

        {#if currentStep === 5}
          <div class="list-content">
            <h1>Installation Summary</h1>
            <InstallSummary bind:product bind:installName={installName} />
          </div>
        {/if}
        <div class="flex justify-end mt-4 space-x-4">
          {#if currentStep > 1}
            <button type="button" on:click={prevStep} class="btn btn-gray">Back</button>
          {/if}
          {#if currentStep < 5}
            <button type="button" on:click={nextStep} class="btn btn-gray">Next</button>
          {/if}
          {#if currentStep === 5}
            <button type="submit" class="btn btn-primary">Install Software</button>
          {/if}
        </div>
      </form>
    </div>
  </div>
</section>

<style>
    .completedStepClass {
        border-color: green;
        /* Add any additional styles for completed steps */
    }

    .currentStepClass {
        border-color: blue;
        /* Add any additional styles for the current step */
    }

    .upcomingStepClass {
        border-color: gray;
        /* Add any additional styles for upcoming steps */
    }
</style>