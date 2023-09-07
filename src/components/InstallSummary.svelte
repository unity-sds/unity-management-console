<script lang="ts">
  import type { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";
  import { getEntries } from "../data/utils";

  export let product: MarketplaceMetadata | undefined;
  export let installName: string;
</script>
<div>
  <h2>Package</h2>
  <ul class="list-group">
    <li class="list-group-item">Install Name: {installName} </li>
    <li class="list-group-item">Version: {product?.Version}</li>
    <li class="list-group-item">Branch:</li>
  </ul>
</div>
<div>
  <h2>Dependencies</h2>
  <ul class="list-group">

  </ul>
</div>
<div>
  <h2>Variables</h2>
  {#if product?.DefaultDeployment?.Variables}
    {#each Object.entries(product?.DefaultDeployment?.Variables) as [key, value], index}
      {#if key === 'AdvancedValues'}
        <div class="mt-12">
          {#each getEntries(value) as [advancedKey, advancedValue]}
            <legend>{advancedKey}</legend>
            {#each getEntries(advancedValue) as [groupKey, groupValue]}
              <ul class="list-group">
                <li class="list-group-item">{groupKey}:
                  {#if typeof groupValue === 'object' && !Array.isArray(groupValue)}
                    {#each getEntries(groupValue) as [subKey, subValue]}
                      <ul class="list-group">
                        <li class="list-group-item">{subKey}: {subValue}</li>
                      </ul>
                    {/each}
                  {:else if Array.isArray(groupValue)}
                    {#each groupValue as item, index}
                      <li class="list-group-item">{groupValue[index]}</li>
                    {/each}
                  {:else}
                    <li class="list-group-item">{groupValue}</li>
                  {/if}
                </li>
              </ul>
            {/each}
          {/each}
        </div>
      {:else if key === 'Values'}
        {#each Object.entries(value) as [valueKey, valueValue]}
          <ul class="list-group">
            <li class="list-group-item">{valueKey}:{valueValue}</li>
          </ul>
        {/each}
      {/if}
    {/each}
  {/if}
</div>