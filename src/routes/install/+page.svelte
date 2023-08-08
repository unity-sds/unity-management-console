<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { install } from "../../store/stores";
	import { HttpHandler } from "../../data/httpHandler";
	import type { Application, InstallationApplication, NodeGroupType, Product } from "../../data/entities";
	import { Writer } from "protobufjs";
	import {Install} from "../../data/unity-cs-manager/protobuf/extensions"
	import type {
		Install_Applications,
	} from "../../data/unity-cs-manager/protobuf/extensions";
	import VariablesForm from "../../components/VariablesForm.svelte";
	import EKSSettings from "../../components/EKSSettings.svelte";
	import ProductForm from "../../components/ProductForm.svelte";



	let product: Product | null = null;
	let variables: Array<[string, string]> = [];
	let hasEks = false;
	let nodeGroups: NodeGroupType[] = [];

	const productString = localStorage.getItem('product');

	onMount(async () => {
		if (productString) {
			product = JSON.parse(productString);
			if (product != null) {
				variables = Object.entries(product.DefaultDeployment.Variables);
				hasEks = product.ManagedDependencies.some((dependency) => dependency.Eks);
			}
		}
	});



	function generateRandomString(length = 10): string {
		let result = '';
		const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
		const charactersLength = characters.length;

		for (let i = 0; i < length; i++) {
			result += characters.charAt(Math.floor(Math.random() * charactersLength));
		}

		return result;
	}

	const installSoftware = async () => {
		if (!product) {
			console.error('No product selected for installation');
			return;
		}

		const httpHandler = new HttpHandler()
		const app: InstallationApplication = {
			name: product.Name,
			version: product.Version,
			variables: new Map([
				['key1', 'value1'],
				['key2', 'value2']
				// add more key-value pairs as needed
			])
		};


		const applications: InstallationApplication[] = [app];

		const id = await httpHandler.installSoftware(applications, "test deployment");
		console.log(id);

		goto('/ui/progress', { replaceState: true });
	};
</script>

<div class="container">
	<div class="row">
		<div class="col-md-12">
			{#if product}
				<h1 class="my-4">{product.Name} Installation</h1>
				<form on:submit|preventDefault={installSoftware}>
					<ProductForm bind:product />
					<!-- only show if uses eks managed dependency -->
					{#if hasEks}
						<EKSSettings bind:product bind:nodeGroups/>
					{/if}

					<VariablesForm bind:variables />
					<button class="btn btn-secondary btn-success mt-3" type="submit">Install</button>
				</form>
			{:else}
				<p>Loading product...</p>
			{/if}
		</div>
	</div>
</div>
