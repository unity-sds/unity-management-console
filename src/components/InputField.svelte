<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let label: string;
  export let id: string;
  export let isValid: boolean;
  export let subtext: string;

  export let value: string;

  const dispatch = createEventDispatcher();

  export let disabled: boolean;

  function handleInput(event: Event) {
    dispatch("input", event);
  }

</script>

<div class="flex flex-col space-y-2">
  <label for={id} class="text-sm font-medium text-gray-700">{label}</label>
  <input type="text" class={`form-input border rounded-md px-3 py-2 ${!isValid ? 'border-red-500' : 'border-gray-300'}`}
         id={id} on:input={handleInput} value={value} disabled={disabled} />
  {#if !isValid}
    <div class="text-sm text-red-500">{label} should be alphanumeric.</div>
  {/if}
  <div class="text-sm text-gray-500">{subtext}</div>
</div>