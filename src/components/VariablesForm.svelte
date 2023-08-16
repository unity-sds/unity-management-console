<script lang="ts">

  import type { MarketplaceMetadata } from "../data/unity-cs-manager/protobuf/marketplace";

  interface NestedValue {
    Config: {
      [key: string]: ConfigValue;
    };
  }
  interface ConfigValue {
    Options: {
      type: string;
      default: string;
    };
  }
  export let product: MarketplaceMetadata | undefined;
  let newVariable = { key: '', value: '' };
  function getConfigValue(configValue: unknown): ConfigValue {
    return configValue as ConfigValue
  }

  function getNestedValue(nestedValue: unknown): NestedValue {
    return nestedValue as NestedValue;
  }


  console.log(product?.DefaultDeployment?.Variables)
  function addVariable() {
    // variables = [...variables, [newVariable.key, newVariable.value]];
    // newVariable.key = '';
    // newVariable.value = '';
  }
  function removeVariable(index: number) {
    //variables = variables.filter((_, i) => i !== index);
  }
</script>
<h2>Variables</h2>
<div class="form-group row mt-4">
<label class="col-sm-2 col-form-label">Variable Override</label>
<input class="form-control" type="text">
</div>
{#if product?.DefaultDeployment?.Variables}
{#each Object.entries(product?.DefaultDeployment?.Variables) as [key, value], index}
  {#if key === 'NestedValues'}
    <div class="row mt-12">
      {#each Object.entries(value) as [nestedKey, nestedValue]}
      <legend>{nestedKey}</legend>
        {#each Object.entries(getNestedValue(nestedValue).Config) as [configKey, configValue]}
          <div class="form-group row mt-4">
          <label class="col-form-label">          </label>

            {configKey}: <input type="text" class="form-control" value={getConfigValue(configValue).Options.default} />
          </div>
        {/each}
      {/each}
    </div>
  {:else if key === 'Values'}
    {#each Object.entries(value) as [valueKey, valueValue]}
      <div class="form-group row mt-4">
        <label class="col-sm-2 col-form-label">{valueKey}:</label>
        <input class="form-control" type="text" bind:value={valueValue} />
      </div>
    {/each}
  {/if}
{/each}
{/if}
<div class="form-group row mt-4">
  <div class="col-sm-2">
    <input
      type="text"
      bind:value={newVariable.key}
      class="form-control"
      placeholder="Variable name"
    />
  </div>
  <div class="col-sm-8">
    <input
      type="text"
      bind:value={newVariable.value}
      class="form-control"
      placeholder="Variable value"
    />
  </div>
  <div class="col-sm-2">
    <button type="button" on:click={addVariable} class="btn btn-secondary"
    >Add Variable</button
    >
  </div>
</div>