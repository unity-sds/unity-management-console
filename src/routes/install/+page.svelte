<script lang="ts">
  import { get } from "svelte/store";
  import type { NodeGroupType } from "../../data/entities";
  import { productInstall } from "../../store/stores";
  import SetupWizard from "../../components/SetupWizard.svelte";

  let nodeGroups: NodeGroupType[] = [];

  let product = get(productInstall);

  function getObjectKeys(obj: object): string[] {
    return Object.keys(obj);
  }

  $: managedDependenciesKeys = product && product.ManagedDependencies ? getObjectKeys(product.ManagedDependencies) : [];
</script>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {#if product}
        <h1 class="my-4">{product.DisplayName} Installation</h1>
        <SetupWizard />
      {:else}
        <p>Loading product...</p>
      {/if}
    </div>
  </div>
</div>
