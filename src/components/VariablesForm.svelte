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
<h2>Variables</h2>
<!--<div class="form-group row mt-4">-->
<!--<label class="col-sm-2 col-form-label">Variable Override</label>-->
<!--<input class="form-control" type="text">-->
<!--</div>-->
{#if product?.DefaultDeployment?.Variables}
  {#each Object.entries(product?.DefaultDeployment?.Variables) as [key, value], index}
    {#if key === 'AdvancedValues'}
      <div class="mt-12">
        {#each getEntries(value) as [advancedKey, advancedValue]}
          <legend>{advancedKey}</legend>
          {#each getEntries(advancedValue) as [groupKey, groupValue]}
            <div class="form-group mt-4">
              <label class="col-sm-2 col-form-label fw-bolder">{groupKey}:</label>
              {#if typeof groupValue === 'object' && !Array.isArray(groupValue)}
                {#each getEntries(groupValue) as [subKey, subValue]}
                  <div class="form-group mt-4">
                    <label class="col-sm-2 col-form-label">{subKey}:</label>
                    <input
                      class="form-control"
                      type="text"
                      value={subValue}
                      on:input={(e) => handleAdvancedInput(e, advancedKey, groupKey, subKey)}
                    />
                  </div>
                {/each}
              {:else if Array.isArray(groupValue)}
                {#each groupValue as item, index}
                  <input
                    class="form-control"
                    type="text"
                    value={item}
                    on:input={(e) => handleAdvancedInput(e, advancedKey, groupKey, null)}
                  />
                {/each}
              {:else}
                <input
                  class="form-control"
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
        <div class="form-group mt-4">
          <label class="col-sm-2 col-form-label">{valueKey}:</label>
          <input
            class="form-control"
            type="text"
            value={valueValue}
            on:input={(e) => handleInput(e, (value) => updateValue(valueKey, value))}
          />
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