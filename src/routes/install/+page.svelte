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

  $: console.log(product);
</script>

<div class="container">
  <div class="st-typography-displayH3">
    Installing Marketplace Application: {product.DisplayName}
  </div>
  <!--   <div class="row">
    <div class="col-md-12">
      {#if product}

        <h1 class="my-4">{product.DisplayName} Installation</h1> -->
  <SetupWizard />
  <!--       {:else}
        <p>Loading product...</p>
      {/if}
    </div>
 -->
  <!-- </div> -->
</div>

<style>
  .container {
    display: flex;
    margin: 10px;
  }
</style>
