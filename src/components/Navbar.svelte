<script>
  export let brand = 'My Website';
  export let logoutLink = '#';
</script>

<div class="bg-white p-8">
  <div class="md:flex md:items-center md:justify-between">
    <div class="min-w-0 flex-1">
      <h2
        class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight"
      >
        {brand}
      </h2>
    </div>
    <div class="mt-4 flex md:ml-4 md:mt-0">
      <a
        href="/management/ui"
        type="button"
        class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
      >
        Home
      </a>
      <a href="/management/ui/settings" type="button">Settings</a>
      <a
        href="/management/ui/progress"
        type="button"
        class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
      >
        Console Log
      </a>
      <button
        type="button"
        class="ml-3 inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
      >
        Logout
      </button>
    </div>
  </div>
</div>
