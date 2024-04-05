<script lang="ts">
  import ProductItem from '../../components/ProductItem.svelte';
  import CategoryList from '../../components/CategoryList.svelte';
  import Header from '../../components/Header.svelte';
  import { marketplaceStore, selectedCategory, order } from '../../store/stores';
  import type { MarketplaceMetadata } from '../../data/unity-cs-manager/protobuf/marketplace';
  import type { OrderLine } from '../../data/entities';
  import { fade, slide } from 'svelte/transition';

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

  // type SelectedProductVersions = Record<string, string>;
  // let selectedProductVersions = <SelectedProductVersions>{};

  // type SelectedProducts = Record<string, MarketplaceMetadata>{};
  // let selectedProducts=<SelectedProducts>{}

  type BinnedProduct = Record<string, MarketplaceMetadata[]>;
  $: binnedProducts = filteredProducts.reduce<BinnedProduct>((acc, product) => {
    acc[product.Name] = acc[product.Name] || [];
    acc[product.Name].push(product);
    return acc;
  }, {});

  // function getSelectedVersion(name: string): MarketplaceMetadata | undefined {
  //   return filteredProducts.find(
  //     (p) => p.Name === name && p.Version === selectedProductVersions[name]
  //   );
  // }
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
              <select>
                {#each productList as product}
                  <option value={product.Version}>{product.Version}</option>
                {/each}
              </select>
            </div>
            <!-- {#if selectedProductVersions[name]}
              <ProductItem product={selectedProductVersions[name]} on:addToCart={handelAddToCart} />
            {/if} -->
          </div>
        {/each}
        <!--         {#each filteredProducts as product}
          <div transition:slide|local={{ duration: 500 }}>
            <ProductItem {product} on:addToCart={handelAddToCart} />
          </div>
        {/each} -->
      </div>
    </div>
  </div>
</div>
