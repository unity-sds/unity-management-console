<script lang="ts">
  import ExtendedSemver from '../../lib/ExtendedSemver';
  import ProductItem from '../../components/ProductItem.svelte';
  import CategoryList from '../../components/CategoryList.svelte';
  import Header from '../../components/Header.svelte';
  import { marketplaceStore, selectedCategory, order, type MarketplaceMetadata } from '../../store/stores';
  import type { OrderLine } from '../../data/entities';
  import { fade, slide } from 'svelte/transition';
  import { goto } from '$app/navigation';

  const semver = new ExtendedSemver();

  $: categories = ['All', ...new Set($marketplaceStore.map((p) => p.Category))];
  $: filteredProducts = $marketplaceStore.filter(
    (p) => $selectedCategory === 'All' || $selectedCategory === p.Category
  );

  const handleSelectCategory = (event: { detail: string }) => {
    $selectedCategory = event.detail;
    categories = categories; // force update to trigger rerender in CategoryList
  };

  const handelAddToCart = (event: { detail: OrderLine }) => {
    const product = event.detail.product;
    const quantity = event.detail.quantity;
    $order.addProduct(product, quantity);
    $order = $order; // force update to trigger rerender in Header
  };

  type BinnedProducts = Record<string, MarketplaceMetadata[]>;
  let binnedProducts = <BinnedProducts>{};

  type SelectedVersionsForProducts = Record<string, MarketplaceMetadata>;
  let selectedVersionsForProducts = <SelectedVersionsForProducts>{};

  $: {
    binnedProducts = filteredProducts.reduce<BinnedProducts>((acc, product) => {
      acc[product.Name] = acc[product.Name] || [];
      acc[product.Name].push(product);
      acc[product.Name] = acc[product.Name].sort((a, b) =>
        semver.compare(semver.coerce(b.Version) || '', semver.coerce(a.Version) || '')
      );
      return acc;
    }, {});
    selectedVersionsForProducts = Object.keys(binnedProducts).reduce<SelectedVersionsForProducts>(
      (acc, name) => {
        acc[name] = binnedProducts[name][0];
        return acc;
      },
      {}
    );
  }

  function handleChangeVersion(name: string) {
    return function handleChange(event: Event) {
      const target = event.target as HTMLSelectElement;
      if (!target.value) return;

      const selectedProduct = binnedProducts[name].find((p) => p.Version === target.value);
      if (selectedProduct) {
        selectedVersionsForProducts[name] = selectedProduct;
      }
    };
  }

  function handleStartInstall(name: string) {
    return () => {
      const product = selectedVersionsForProducts[name];
      // Navigate to the install page with name and version as URL parameters only
      goto(`/management/ui/install?name=${encodeURIComponent(product.Name)}&version=${encodeURIComponent(product.Version)}`, { replaceState: true });
    };
  }
</script>

<div>
  <Header />
  <div class="w-full mx-auto">
    <div class="flex">
      <div class="w-1/4 p-2" in:fade={{ duration: 500 }}>
        <CategoryList {categories} on:selectCategory={handleSelectCategory} />
      </div>
      <div class="w-3/4 p-2">
        {#each Object.entries(binnedProducts) as [name, productList]}
          <div>
            <div class="px-4 sm:px-0" style="display: flex; gap: 10px; align-items: center;">
              <h2 class="font-semibold leading-7 text-gray-900 text-2xl">
                {name}
              </h2>
              <select
                value={selectedVersionsForProducts[name].Version}
                on:change={handleChangeVersion(name)}
              >
                {#each productList as product}
                  <option value={product.Version}>{product.Version}</option>
                {/each}
              </select>
              <button class="st-button" on:click={handleStartInstall(name)}>Install</button>
            </div>
            <ProductItem
              product={selectedVersionsForProducts[name]}
              on:addToCart={handelAddToCart}
            />
          </div>
          <hr style="margin-bottom: 25px;" />
        {/each}
      </div>
    </div>
  </div>
</div>
