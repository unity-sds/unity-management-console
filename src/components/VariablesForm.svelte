<script lang="ts">

  import type { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";
  import { getConfigValue, getEntries, getNestedValue } from "../data/utils";

  export let product: MarketplaceMetadata;

  function handleAdvancedInput(e: Event, advancedKey: string, groupKey: string, subKey: string | null) {
    const target = e.target as HTMLInputElement;
    updateAdvancedValue(advancedKey, groupKey, subKey, target.value);
  }

  function handleInput(e: Event, callback: (value: string) => void) {
    const target = e.target as HTMLInputElement;
    callback(target.value);
  }

  function updateAdvancedValue(advancedKey: string, groupKey: string, subKey: string | null, newValue: any) {
    if (
      !product ||
      !product.DefaultDeployment ||
      !product.DefaultDeployment.Variables ||
      !product.DefaultDeployment.Variables.AdvancedValues ||
      !product.DefaultDeployment.Variables.AdvancedValues[advancedKey] ||
      !product.DefaultDeployment.Variables.AdvancedValues[advancedKey][groupKey]
    ) {
      console.error("Invalid product structure");
      return;
    }

    // The update logic, depending on whether there's a subKey
    let updatedGroupValue = subKey
      ? {
        ...product.DefaultDeployment.Variables.AdvancedValues[advancedKey][groupKey],
        [subKey]: newValue
      }
      : newValue;

    product = {
      ...product,
      DefaultDeployment: {
        ...product.DefaultDeployment,
        Variables: {
          ...product.DefaultDeployment.Variables,
          AdvancedValues: {
            ...product.DefaultDeployment.Variables.AdvancedValues,
            [advancedKey]: {
              ...product.DefaultDeployment.Variables.AdvancedValues[advancedKey],
              [groupKey]: updatedGroupValue
            }
          }
        }
      }
    };
  }

  function updateValue(valueKey: string, newValue: any) {
    if (
      !product ||
      !product.DefaultDeployment ||
      !product.DefaultDeployment.Variables
    ) {
      console.error("Invalid product structure");
      return;
    }

    product = {
      ...product,
      DefaultDeployment: {
        ...product.DefaultDeployment,
        Variables: {
          ...product.DefaultDeployment.Variables,
          Values: {
            ...product.DefaultDeployment.Variables.Values,
            [valueKey]: newValue
          }
        }
      }
    };
  }
</script>
<h2 class="text-xl font-bold">Variables</h2>

{#if product?.DefaultDeployment?.Variables}
  {#each Object.entries(product?.DefaultDeployment?.Variables) as [key, value], index}
    {#if key === 'AdvancedValues'}
      <div class="mt-12">
        {#each getEntries(value) as [advancedKey, advancedValue]}
          <legend class="font-medium">{advancedKey}</legend>
          {#each getEntries(advancedValue) as [groupKey, groupValue]}
            <div class="mt-4">
              <h3 class="block text-sm font-medium">{groupKey}:</h3>
              {#if typeof groupValue === 'object' && !Array.isArray(groupValue)}
                {#each getEntries(groupValue) as [subKey, subValue]}
                  <div class="mt-4">
                    <label class="block text-sm font-medium">{subKey}:
                      <input
                        class="mt-1 p-2 w-full border rounded-md"
                        type="text"
                        value={subValue}
                        on:input={(e) => handleAdvancedInput(e, advancedKey, groupKey, subKey)}
                      />
                    </label>
                  </div>
                {/each}
              {:else if Array.isArray(groupValue)}
                {#each groupValue as item, index}
                  <input
                    class="mt-1 p-2 w-full border rounded-md"
                    type="text"
                    value={item}
                    on:input={(e) => handleAdvancedInput(e, advancedKey, groupKey, null)}
                  />
                {/each}
              {:else}
                <input
                  class="mt-1 p-2 w-full border rounded-md"
                  type="text"
                  value={groupValue}
                  on:input={(e) => handleAdvancedInput(e, advancedKey, groupKey, null)}
                />
              {/if}
            </div>
          {/each}
        {/each}
      </div>
    {:else if key === 'Values'}
      {#each Object.entries(value) as [valueKey, valueValue]}
        <div class="mt-4">
          <label class="block text-sm font-medium">{valueKey}:
            <input
              class="mt-1 p-2 w-full border rounded-md"
              type="text"
              value={valueValue}
              on:input={(e) => handleInput(e, (value) => updateValue(valueKey, value))}
            />
          </label>
        </div>
      {/each}
    {/if}
  {/each}
{/if}
<div class="form-group mt-4">
  <!--  <div class="col-sm-2">-->
  <!--    <input-->
  <!--      type="text"-->
  <!--      bind:value={newVariable.key}-->
  <!--      class="form-control"-->
  <!--      placeholder="Variable name"-->
  <!--    />-->
  <!--  </div>-->
  <!--  <div class="col-sm-8">-->
  <!--    <input-->
  <!--      type="text"-->
  <!--      bind:value={newVariable.value}-->
  <!--      class="form-control"-->
  <!--      placeholder="Variable value"-->
  <!--    />-->
  <!--  </div>-->
  <!--  <div class="col-sm-2">-->
  <!--    <button type="button" on:click={addVariable} class="btn btn-secondary"-->
  <!--    >Add Variable-->
  <!--    </button-->
  <!--    >-->
  <!--  </div>-->
</div>