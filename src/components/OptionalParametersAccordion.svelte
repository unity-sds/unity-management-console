<script lang="ts">
  export let list: { key: string; value: string }[] = [];
  export let key = "";
  export let value = "";

  let searchTerm = "";

  function addToList() {
    list = [...list, { key: key, value: value }];
    key = "";
    value = "";
  }

  function removeFromList(index: number) {
    list = list.filter((_, i) => i !== index);
  }

  $: filteredList = list.filter(item => item.key.toLowerCase().includes(searchTerm.toLowerCase()));
</script>
<div class="space-y-4 mt-4" id="optionalParametersAccordion">
  <div>
    <h2 class="px-4 py-2 border rounded-t-md">
      <button
        class="w-full text-left"
        type="button"
        aria-expanded="false"
        aria-controls="collapseOne"
      >
        Optional Parameters
      </button>
    </h2>
    <div
      id="collapseOne"
      class="border-t-0 rounded-b-md border px-4 py-2"
      aria-labelledby="headingOne"
    >
      <div class="mt-2 mb-4">
        <input
          type="text"
          bind:value={searchTerm}
          placeholder="Search by key..."
          class="border rounded-md px-3 py-2 w-full text-gray-700"
        />
      </div>
      <div>
        <div class="space-y-2 mt-3">
          <ul>
            {#each filteredList as item, index}
              <li class="border rounded-md px-4 py-2 space-y-2">
                <div class="flex justify-between items-center">
                  <span>{item.key}:</span>
                  <button
                    type="button"
                    class="px-2 py-1 bg-red-600 text-white rounded-md hover:bg-red-700 focus:bg-red-800"
                    on:click={() => removeFromList(index)}
                  >
                    Remove
                  </button>
                </div>
                <input
                  type="text"
                  bind:value={list[index].value}
                  class="border rounded-md px-2 py-1 w-full text-gray-700" />
              </li>
            {/each}
          </ul>
        </div>
      </div>
      <div class="flex flex-col space-y-2 mt-2">
        <label for="key" class="text-sm font-medium text-gray-700">Key</label>
        <input type="text" id="key" bind:value={key} class="border rounded-md px-3 py-2 w-full text-gray-700" />
      </div>

      <div class="flex flex-col space-y-2 mt-2">
        <label for="value" class="text-sm font-medium text-gray-700">Value</label>
        <input type="text" id="value" bind:value={value} class="border rounded-md px-3 py-2 w-full text-gray-700" />
      </div>

      <button type="button"
              class="mt-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:bg-blue-800"
              on:click={addToList}>
        Add
      </button>

    </div>
  </div>
</div>