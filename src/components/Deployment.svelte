<script>
  export let installName = "";
</script>

<div class="form-group mb-4">
  <label for="name" class="block text-sm font-medium text-gray-700">Deployment Name</label>
  <input id="name" class="mt-1 p-2 w-full border rounded-md" bind:value={installName} required />
</div>