<script lang="ts">
  import {get} from 'svelte/store'
  import { goto } from "$app/navigation";
  import { HttpHandler } from "../../data/httpHandler";
  import type { NodeGroupType } from "../../data/entities";
  import VariablesForm from "../../components/VariablesForm.svelte";
  import ProductForm from "../../components/ProductForm.svelte";
  import { Install_Applications, Install_Variables } from "../../data/unity-cs-manager/protobuf/extensions";
  import { productInstall } from "../../store/stores";

  let nodeGroups: NodeGroupType[] = [];

  let product = get(productInstall);

  const installSoftware = async () => {
    if (!product) {
      console.error("No product selected for installation");
      return;
    } else{
      console.log(product)
    }

    const httpHandler = new HttpHandler();

    const merged = {"Values": product.DefaultDeployment?.Variables?.Values, "NestedValues": product.DefaultDeployment?.Variables?.NestedValues, "AdvancedValues": product.DefaultDeployment?.Variables?.AdvancedValues}
    const vars = Install_Variables.fromJSON(merged)
    const a = Install_Applications.create({
      name: product.Name,
      version: product.Version,
      variables: vars
    } as any)
    const id = await httpHandler.installSoftware(a, "test_deployment");
    console.log(id);
    goto("/ui/progress", { replaceState: true });
  };
</script>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {#if product}
        <h1 class="my-4">{product.Name} Installation</h1>
        <form on:submit|preventDefault={installSoftware}>
          <ProductForm bind:product />
          {#if product.ManagedDependencies}
            <h2>Dependencies</h2>
            {#each product.ManagedDependencies as dependency}
              {#each Object.keys(dependency) as key}
              <div class="form-group">
<!--                  <strong>{key}</strong>: Minimum Version - {dependency[key].MinimumVersion}-->
                  <label class="col-form-label">{key} <select class="form-control" ><option></option><option>test_deployment</option></select></label>
              </div>
            {/each}
              {/each}
          {/if}
          <VariablesForm bind:product />
          <button class="btn btn-secondary btn-success mt-3" type="submit">Install</button>
        </form>
      {:else}
        <p>Loading product...</p>
      {/if}
    </div>
  </div>
</div>
