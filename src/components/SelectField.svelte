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

<div class="flex flex-col space-y-2">
  {#if multiple}
    <label for={`${id}-multiple`} class="text-sm font-medium text-gray-700">{label}</label>
    <select bind:value={value} id={`${id}-multiple`} on:change={handleChange} multiple
            class="border rounded-md px-3 py-2 w-full text-gray-700">
      {#each options as option, index}
        <option value={option.value ? option.value : option}>
          {option.label ? option.label : option}
        </option>
      {/each}
    </select>
  {:else}
    <label for={`${id}-single`} class="text-sm font-medium text-gray-700">{label}</label>
    <select bind:value={value} id={`${id}-single`} on:change={handleChange}
            class="border rounded-md px-3 py-2 w-full text-gray-700">
      {#each options as option, index}
        <option value={option.value ? option.value : option}>
          {option.label ? option.label : option}
        </option>
      {/each}
    </select>
  {/if}
  <div class="text-sm text-gray-500 mt-1">{subtext}</div>
</div>
<style>
    /* Add any additional styles here */
</style>