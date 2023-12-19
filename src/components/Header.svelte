<script lang="ts">
  import { order } from "../store/stores";
  import { fade, fly } from "svelte/transition";
  import { quintOut } from "svelte/easing";

  $: count = $order.productCount || 0;
  $: displayText1 = count === 0 ? "(No Selection)" : `${count} product(s), `;
  $: displayText2 = count === 0 ? "" : `$${$order.total.toFixed(2)}`;
</script>

<div class="p-1 bg-gray-600 text-white text-right">
  {#key displayText1}
		<span class="inline-block" in:fade={{ duration: 500 }}>
			{displayText1}
		</span>
  {/key}
  {#key displayText2}
		<span class="inline-block" in:fly={{ y: -50, duration: 500, easing: quintOut }}>
			{displayText2}
		</span>
  {/key}
  <a href="/management/ui" class="st-button large my-1 mx-1"> Home </a>
</div>