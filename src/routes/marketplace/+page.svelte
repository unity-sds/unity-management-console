<script lang="ts">
  import ProductItem from "../../components/ProductItem.svelte";
  import CategoryList from "../../components/CategoryList.svelte";
  import Header from "../../components/Header.svelte";
  import { marketplaceStore, selectedCategory, order } from "../../store/stores";
  import type { OrderLine } from "../../data/entities";
  import { fade, slide } from "svelte/transition";

  $: categories = ["All", ...new Set($marketplaceStore.map((p) => p.Category))];
  $: filteredProducts = $marketplaceStore.filter(
    (p) => $selectedCategory === "All" || $selectedCategory === p.Category
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
</script>

<div>
  <Header />
  <div class="w-full mx-auto">
    <div class="flex">
      <div class="w-1/4 p-2" in:fade={{ duration: 500 }}>
        <CategoryList {categories} on:selectCategory={handleSelectCategory} />
      </div>
      <div class="w-3/4 p-2">
        {#each filteredProducts as product}
          <div transition:slide|local={{ duration: 500 }}>
            <ProductItem {product} on:addToCart={handelAddToCart} />
          </div>
        {/each}
      </div>
    </div>
  </div>
</div>