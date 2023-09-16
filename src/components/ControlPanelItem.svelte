<script>
  export let title = "";
  export let description = "";
  export let link = "";
  export let disabled = false;

</script>

<div class="lg:w-1/3 md:w-1/2 mb-4">
  <div class="bg-white shadow rounded overflow-hidden h-full">
    <div class="p-6">
      <h4 class="text-xl font-bold">{title}</h4>
      <p class="mt-2">{description}</p>
    </div>
    <div class="p-4 border-t border-gray-200">
      {#if !disabled}
        <a href="{link}" class="px-4 py-2 text-white bg-blue-600 hover:bg-blue-700 rounded">Go!</a>
      {:else}
        <span class="px-4 py-2 text-gray-400 bg-gray-300 cursor-not-allowed rounded">Go!</span>
      {/if}
    </div>
  </div>
</div>
