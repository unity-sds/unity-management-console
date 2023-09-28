<script lang="ts">

  import type { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";

  export let product: MarketplaceMetadata | null;

</script>

{#if product}
  <div class="form-group mb-4">
    <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
    <input id="name" class="mt-1 p-2 w-full border rounded-md" bind:value={product.DisplayName} required />
  </div>

  <div class="form-group mb-4">
    <label for="version" class="block text-sm font-medium text-gray-700">Version</label>
    <input id="version" class="mt-1 p-2 w-full border rounded-md" bind:value={product.Version} />
  </div>

  <div class="form-group">
    <label for="branch" class="block text-sm font-medium text-gray-700">Branch</label>
    <input id="branch" class="mt-1 p-2 w-full border rounded-md" />
  </div>
{/if}