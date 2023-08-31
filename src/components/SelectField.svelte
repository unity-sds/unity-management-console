<script lang="ts">
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let label = ""; // Label for the field
  export let id = ""; // HTML ID for the select field
  export let multiple = false; // Flag to indicate if multiple selections are allowed
  export let options: { value?: string; label?: string }[] = []; // Options for the select field
  export let value: string[] = []; // Array of selected values

  function handleChange(e: Event) {
    dispatch("change", (e.target as HTMLSelectElement).value);
  }

  export let subtext = "";
</script>

<div class="form-group">
  {#if multiple}
    <label for={`${id}-multiple`}>{label}</label>
    <select bind:value={value} id={`${id}-multiple`} class="form-control" on:change={handleChange} multiple>
      {#each options as option, index}
        <option value={option.value ? option.value : option}>
          {option.label ? option.label : option}
        </option>
      {/each}
    </select>
  {:else}
    <label for={`${id}-single`}>{label}</label>
    <select bind:value={value} id={`${id}-single`} on:change={handleChange} class="form-control">
      {#each options as option, index}
        <option value={option.value ? option.value : option}>
          {option.label ? option.label : option}
        </option>
      {/each}
    </select>
  {/if}
  <div class="form-text">{subtext}</div>

</div>

<style>
    /* Add any additional styles here */
</style>