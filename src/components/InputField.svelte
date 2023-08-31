<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let label: string;
  export let id: string;
  export let isValid: boolean;
  export let subtext: string;

  export let value: string;

  const dispatch = createEventDispatcher();

  function handleInput(event: Event) {
    dispatch("input", event);
  }

</script>

<div class="form-group">
  <label for={id}>{label}</label>
  <input type="text" class="form-control" id={id} class:invalid={!isValid} on:input={handleInput} value={value} />
  {#if !isValid}
    <div class="invalid-feedback" style="display: block;">{label} should be alphanumeric.</div>
  {/if}
  <div class="form-text">{subtext}</div>

</div>