<script lang="ts">

  import type { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";

  export let product: MarketplaceMetadata | null;

  export let installName = "";
</script>

{#if product}
  <div class="form-group">
    <label for="name">Name</label>
    <input id="name" class="form-control" bind:value={installName} required />
  </div>
  <div class="form-group mt-4">
    <label for="version">Version</label>
    <input id="version" class="form-control" bind:value={product.Version} />
  </div>
  <div class="form-group mt-4">
    <label for="branch">Branch</label>
    <input id="branch" class="form-control" />
  </div>
{/if}
