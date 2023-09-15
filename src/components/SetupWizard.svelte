<script lang="ts">
  import ProductForm from "./ProductForm.svelte";
  import { deploymentStore, productInstall } from "../store/stores";
  import VariablesForm from "./VariablesForm.svelte";
  import { HttpHandler } from "../data/httpHandler";
  import { Deployments, Install_Applications, Install_Variables } from "../data/unity-cs-manager/protobuf/extensions";
  import { goto } from "$app/navigation";
  import InstallSummary from "./InstallSummary.svelte";
  import { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";

  let product: MarketplaceMetadata = MarketplaceMetadata.create();

  productInstall.subscribe(value => {
    product = value;
  });

  let currentStep = 1;

  function navigateToStep(step: number) {
    currentStep = step;
  }

  function nextStep() {
    currentStep++;
  }

  function prevStep() {
    currentStep--;
  }

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
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

    const removeEmptyStrings = (obj: AnyObject): void => {
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
      variables: vars
    } as any);
    const id = await httpHandler.installSoftware(a, installName);
    console.log(id);
    goto("/ui/progress", { replaceState: true });
  };

  $: managedDependenciesKeys = product && product.ManagedDependencies ? getObjectKeys(product.ManagedDependencies) : [];

  let deployed: Deployments;

  deploymentStore.subscribe(value => {
    deployed = value;
  });

  function getVersionsForKey(key): string[] {
    let options: string[] = [];
    for (let d of deployed.deployment) {
      for (let a of d.application) {
        if (key === a.packageName) {
          options.push(a.applicationName);
        }
      }
    }
    return options;
  }

</script>

<section class="section">
  <div class="container">
    <div class="row">
      <div class="col-md-12">
        <div class="wizard">
          <div class="wizard-inner">
            <div class="connecting-line"></div>
            <ul class="nav nav-tabs" role="tablist">
              <li role="presentation" class={currentStep === 1 ? 'active' : currentStep > 1 ? '' : 'disabled'}>
                <a on:click={() => navigateToStep(1)} href="#step1" aria-controls="step1" role="tab"><span
                  class="round-tab">1</span> <i>Package</i></a>
              </li>
              <li role="presentation" class={currentStep === 2 ? 'active' : currentStep > 2 ? '' : 'disabled'}>
                <a on:click={() => navigateToStep(2)} href="#step2" aria-controls="step2" role="tab"><span
                  class="round-tab">2</span> <i>Dependencies</i></a>
              </li>
              <li role="presentation" class={currentStep === 3 ? 'active' : currentStep > 3 ? '' : 'disabled'}>
                <a on:click={() => navigateToStep(3)} href="#step3" aria-controls="step3" role="tab"><span
                  class="round-tab">3</span> <i>Variables</i></a>
              </li>
              <li role="presentation" class={currentStep === 4 ? 'active' : 'disabled'}>
                <a on:click={() => navigateToStep(4)} href="#step4" aria-controls="step4" role="tab"><span
                  class="round-tab">4</span> <i>Overview</i></a>
              </li>
            </ul>
          </div>
          <div class="content">
            <form on:submit|preventDefault={installSoftware}>
              {#if currentStep === 1}
                <ProductForm bind:product bind:installName={installName} />
              {/if}

              {#if currentStep === 2}
                <!-- Collapse 2 content -->
                <div class="list-content">
                  {#if product.ManagedDependencies}
                    <h2>Dependencies</h2>
                    <!--{#each product.ManagedDependencies as dependency}-->
                    {#each managedDependenciesKeys as key}
                      <div class="form-group">
                        <!--                  <strong>{key}</strong>: Minimum Version - {dependency[key].MinimumVersion}-->
                        <label class="col-form-label">{key} <select class="form-control">
                          <option></option>
                          {#each getVersionsForKey(key) as version}
                            <option>{version}</option>
                          {/each}
                        </select></label>
                      </div>
                    {/each}
                    <!--{/each}-->
                  {/if}
                </div>
              {/if}

              {#if currentStep === 3}
                <!-- Collapse 3 content -->
                <div class="list-content">
                  <VariablesForm bind:product />
                </div>
              {/if}

              {#if currentStep === 4}
                <div class="list-content">
                  <h1>Installation Summary</h1>
                  <InstallSummary bind:product bind:installName={installName} />
                </div>
              {/if}
              <ul class="list-inline pull-right">
                {#if currentStep > 1}
                  <li>
                    <button on:click={prevStep} type="button" class="default-btn prev-step">Back</button>
                  </li>
                {/if}
                {#if currentStep < 4}
                  <li>
                    <button on:click={nextStep} type="button" class="default-btn next-step">Next</button>
                  </li>
                {/if}
                {#if currentStep === 4}
                  <li>
                    <button type="submit" class="btn btn-primary">Install Software</button>
                  </li>
                {/if}
              </ul>
            </form>
          </div>
          <div class="clearfix"></div>
        </div>
      </div>
    </div>
  </div>
</section>

<style>
    .section {
        margin-top: 50px;
    }

    .nav-tabs li {
        display: inline-block;
        width: 25%; /* Assuming you have 4 steps; adjust as necessary */
        text-align: center; /* Center the content of each list item */
    }

    .round-tab {
        display: inline-block;
        background-color: #e0e0e0;
        width: 30px; /* Or whatever size you want */
        height: 30px;
        line-height: 30px; /* Vertically center the number inside */
        border-radius: 50%; /* Makes it round */
        text-align: center; /* Horizontally center the number inside */
        margin: 0 auto; /* Center the round-tab in the li if there's extra space */
    }

    .nav-tabs {
        padding: 0;
        list-style: none; /* Remove bullet points */
    }

    .nav-tabs li {
        margin: 0; /* Remove any existing margins */
        padding: 0; /* Remove any existing padding */
    }

    /* Styling for active and disabled tabs */
    .nav-tabs li.active a {
        background-color: #f5f5f5;
        border-bottom: 2px solid #00bcd4;
    }

    .nav-tabs li.disabled a {
        pointer-events: none;
        cursor: default;
        color: #ccc;
    }

    /* Style for the connecting line */
    .connecting-line {
        width: 75%;
        margin: 0 auto;
        height: 2px;
        background: #e0e0e0;
        position: absolute;
        z-index: 1;
        top: 15px;
        left: 0;
        right: 0;
    }

    /*------------------------*/
    button:focus,
    .form-control:focus {
        outline: none;
        box-shadow: none;
    }

    .form-control:disabled {
        background-color: #fff;
    }


    .wizard .nav-tabs {
        position: relative;
        margin-bottom: 0;
        border-bottom-color: transparent;
    }

    .wizard > div.wizard-inner {
        position: relative;
        margin-bottom: 50px;
        text-align: center;
    }


    .wizard .nav-tabs > li.active > a, .wizard .nav-tabs > li.active > a:hover, .wizard .nav-tabs > li.active > a:focus {
        color: #555555;
        cursor: default;
        border: 0;
        border-bottom-color: transparent;
    }

    span.round-tab {
        width: 30px;
        height: 30px;
        line-height: 30px;
        display: inline-block;
        border-radius: 50%;
        background: #fff;
        z-index: 2;
        position: absolute;
        left: 0;
        text-align: center;
        font-size: 16px;
        color: #0e214b;
        font-weight: 500;
        border: 1px solid #ddd;
    }

    .wizard li.active span.round-tab {
        background: #0db02b;
        color: #fff;
        border-color: #0db02b;
    }

    .wizard li.active span.round-tab i {
        color: #5bc0de;
    }

    .wizard .nav-tabs > li.active > a i {
        color: #0db02b;
    }

    .wizard .nav-tabs > li {
        width: 25%;
    }

    .wizard li:after {
        content: " ";
        position: absolute;
        left: 46%;
        opacity: 0;
        margin: 0 auto;
        bottom: 0;
        border: 5px solid transparent;
        border-bottom-color: red;
        transition: 0.1s ease-in-out;
    }


    .wizard .nav-tabs > li a {
        width: 30px;
        height: 30px;
        margin: 20px auto;
        border-radius: 100%;
        padding: 0;
        background-color: transparent;
        position: relative;
        top: 0;
    }

    .wizard .nav-tabs > li a i {
        position: absolute;
        top: -15px;
        font-style: normal;
        white-space: nowrap;
        left: 50%;
        transform: translate(-50%, -50%);
        font-size: 12px;
        font-weight: 700;
        color: #000;
    }

    .wizard .nav-tabs > li a:hover {
        background: transparent;
    }

    .prev-step,
    .next-step {
        font-size: 13px;
        padding: 8px 24px;
        border: none;
        border-radius: 4px;
        margin-top: 30px;
    }

    .next-step {
        background-color: #0db02b;
    }

    .list-content {
        margin-bottom: 10px;
    }

    .list-content a i {
        text-align: right;
        position: absolute;
        top: 15px;
        right: 10px;
        transition: 0.5s;
    }

    .form-control {
        background-color: #fdfdfd;
    }

    .nav > li {
        padding: 0;
    }

    .list-inline li {
        display: inline-block;
    }

    .pull-right {
        float: right;
    }

    /*-----------custom-checkbox-----------*/
    /*----------Custom-Checkbox---------*/


    @media (max-width: 767px) {

        .wizard .nav-tabs > li a i {
            display: none;
        }

    }
</style>