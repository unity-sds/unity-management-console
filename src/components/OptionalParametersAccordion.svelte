<script lang="ts">
  export let list: { key: string; value: string }[] = [];
  export let key = "";
  export let value = "";

  function addToList() {
    list = [...list, { key: key, value: value }];
    key = "";
    value = "";
  }

  function removeFromList(index: number) {
    list = list.filter((_, i) => i !== index);
  }
</script>

<div class="accordion mt-4" id="optionalParametersAccordion">
  <div class="accordion-item">
    <h2 class="accordion-header" id="headingOne">
      <button
        class="accordion-button collapsed"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#collapseOne"
        aria-expanded="false"
        aria-controls="collapseOne"
      >
        Optional Parameters
      </button>
    </h2>
    <div
      id="collapseOne"
      class="accordion-collapse collapse"
      aria-labelledby="headingOne"
      data-bs-parent="#optionalParametersAccordion"
    >
      <div class="accordion-body">
        <div class="form-group">
          <label for="key">Key</label>
          <input type="text" class="form-control" id="key" bind:value={key} />
        </div>

        <div class="form-group">
          <label for="value">Value</label>
          <input type="text" class="form-control" id="value" bind:value={value} />
        </div>

        <button type="button" class="btn btn-primary mt-2" on:click={addToList}>
          Add
        </button>

        <div class="mt-3">
          <ul class="list-group">
            {#each list as item, index}
              <li class="list-group-item">
                {item.key}: {item.value}
                <button
                  type="button"
                  class="btn btn-danger btn-sm float-right"
                  on:click={() => removeFromList(index)}
                >
                  Remove
                </button>
              </li>
            {/each}
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
    /* You can add additional styles here */
</style>