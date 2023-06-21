<script lang="ts">
  export let variables
  let newVariable = { key: '', value: '' };

  function addVariable() {
    variables = [...variables, [newVariable.key, newVariable.value]];
    newVariable.key = '';
    newVariable.value = '';
  }
  function removeVariable(index: number) {
    variables = variables.filter((_, i) => i !== index);
  }
</script>
<h2>Variables</h2>
{#each variables as variable, index (variable[0])}
  <div class="form-group row mt-4">
    <label for={variable[0]} class="col-sm-2 col-form-label">{variable[0]}</label>
    <div class="col-sm-8">
      <input type="text" class="form-control" id={variable[0]} value={variable[1]} />
    </div>
    <div class="col-sm-2">
      <button
        type="button"
        on:click={() => removeVariable(index)}
        class="btn btn-danger st-button secondary
                    large">Remove</button
      >
    </div>
  </div>
{/each}
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