<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Product,Application } from '../data/entities';
import { goto } from '$app/navigation';
	const dispatch = createEventDispatcher<{ installApp: Application }>();

	export let product: Product;
	let quantity: string;

	const handleInstallApp = () => {
      localStorage.setItem('product', JSON.stringify(product));  // Store product in local storage
      goto('/install', { replaceState: true });
	};
</script>

<div class="card m-1 p-1 bg-light">
	<h4>
		{product.Name}
		<span class="badge rounded-pill bg-primary float-end">
		</span>
	</h4>
	<div class="card-text bg-white p-1">
		  Description: {product.Description}<br/>
      Package Location: {product.Package}<br/>
      Orchestration Type: {product.Backend}<br/>
      Tags: {product.Tags}<br/>
      Managed Dependencies: <pre>{JSON.stringify(product.ManagedDependencies, null, 2)}</pre> <br/>
      Minimum Iam Roles: <pre>{JSON.stringify(product.IamRoles, null, 2)}</pre>
		<button class="st-button large float-end" on:click={handleInstallApp}>
			Install Application
		</button>
	</div>
</div>
