<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Application } from '../data/entities';
  import { goto } from '$app/navigation';
  import type { MarketplaceMetadata } from '../store/stores';

  const dispatch = createEventDispatcher<{ installApp: Application }>();

  export let product: MarketplaceMetadata;

  const handleInstallApp = () => {
    // Navigate with URL parameters only
    goto(`/management/ui/install?name=${encodeURIComponent(product.Name)}&version=${encodeURIComponent(product.Version)}`, { replaceState: true });
  };
</script>

<div>
  <div class="px-4 sm:px-0">
    <!-- <h2 class="font-semibold leading-7 text-gray-900 text-2xl">
      {product.Name}
    </h2> -->
    <div class="mt-6 border-t border-gray-100">
      <div class="divide-y divide-gray-100">
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Description</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
            {product.Description}
          </dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Package Location</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
            {product.Package}
          </dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Orchestration Type</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
            {product.Backend}
          </dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Tags</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">{product.Tags}</dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Managed Dependencies</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
            <pre class="bg-gray-100 p-2 rounded">{JSON.stringify(
                product.ManagedDependencies,
                null,
                2
              )}</pre>
          </dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
          <dt class="text-sm font-medium leading-6 text-gray-900">Minimum IAM Roles</dt>
          <dd class="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
            <pre class="bg-gray-100 p-2 rounded">{JSON.stringify(product.IamRoles, null, 2)}</pre>
          </dd>
        </div>
        <div class="px-4 py-6 sm:grid sm:grid-cols-1 sm:gap-4 sm:px-0">
          <!-- <button
            class="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded float-right"
            on:click={handleInstallApp}
            id={`${product.Name}-v${product.Version}`}
          >
            Install Application
          </button> -->
        </div>
      </div>
    </div>
  </div>
</div>
