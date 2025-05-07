<script lang="ts">
  import type { MarketplaceMetadata } from "../store/stores";
  import { getEntries } from "../data/utils";

  export let product: MarketplaceMetadata | undefined;
  export let installName: string;
</script>
<div>
  <h2 class="text-xl font-bold">Package</h2>
  <ul class="divide-y divide-gray-200">
    <li class="py-2 px-4 bg-gray-50">Install Name: {installName}</li>
    <li class="py-2 px-4 bg-gray-50">Version: {product?.Version}</li>
    <li class="py-2 px-4 bg-gray-50">Branch:</li>
  </ul>
</div>
<div>
  <h2 class="text-xl font-bold mt-6">Dependencies</h2>
  <ul class="divide-y divide-gray-200">

  </ul>
</div>
<div>
  <h2 class="text-xl font-bold mt-6">Variables</h2>
  {#if product?.DefaultDeployment?.Variables}
    {#each Object.entries(product?.DefaultDeployment?.Variables) as [key, value], index}
      {#if key === 'AdvancedValues'}
        <div class="mt-12">
          {#each getEntries(value) as [advancedKey, advancedValue]}
            <legend class="font-medium">{advancedKey}</legend>
            {#each getEntries(advancedValue) as [groupKey, groupValue]}
              <ul class="divide-y divide-gray-200">
                <li class="py-2 px-4 bg-gray-50">{groupKey}:
                  {#if typeof groupValue === 'object' && !Array.isArray(groupValue)}
                    {#each getEntries(groupValue) as [subKey, subValue]}
                      <ul class="divide-y divide-gray-200">
                        <li class="py-2 px-4 bg-gray-50">{subKey}: {subValue}</li>
                      </ul>
                    {/each}
                  {:else if Array.isArray(groupValue)}
                    {#each groupValue as item, index}
                      <li class="py-2 px-4 bg-gray-50">{groupValue[index]}</li>
                    {/each}
                  {:else}
                    <li class="py-2 px-4 bg-gray-50">{groupValue}</li>
                  {/if}
                </li>
              </ul>
            {/each}
          {/each}
        </div>
      {:else if key === 'Values'}
        {#each Object.entries(value) as [valueKey, valueValue]}
          <ul class="divide-y divide-gray-200">
            <li class="py-2 px-4 bg-gray-50">{valueKey}:{valueValue}</li>
          </ul>
        {/each}
      {/if}
    {/each}
  {/if}
</div>