<script lang="ts">
	import type { Product } from '../data/entities';

	export let product: Product | null;
</script>

{#if product}
	<div class="form-group">
		<label for="name">Name</label>
		<input id="name" class="form-control" bind:value={product.Name} required />
	</div>
	<div class="form-group mt-4">
		<label for="version">Version</label>
		<input id="version" class="form-control" bind:value={product.Version} />
	</div>
	<div class="form-group mt-4">
		<label for="branch">Branch</label>
		<input id="branch" class="form-control" bind:value={product.Branch} />
	</div>
{/if}
