<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { install } from "../../store/stores";
	import { HttpHandler } from "../../data/httpHandler";
	import type { Product } from "../../data/entities";
	import { Writer } from "protobufjs";
	import {Install} from "../../data/protobuf/extensions"
	import type {
		Install_Applications,
		Install_Extensions,
		Install_Extensions_Eks, Install_Extensions_Nodegroups
	} from "../../data/protobuf/extensions";
	import VariablesForm from "../../components/VariablesForm.svelte";
	import EKSSettings from "../../components/EKSSettings.svelte";
	import ProductForm from "../../components/ProductForm.svelte";

	interface NodeGroupType {
		name: string;
		settings: {
			MinNodes: number;
			MaxNodes: number;
			DesiredNodes: number;
			InstanceType: string;
		};
	}

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


	const generateNodeGroup = (nodegroup: NodeGroupType): Install_Extensions_Nodegroups => {
		return {
			name: nodegroup.name,
			nodecount: nodegroup.settings.DesiredNodes.toString(),
			instancetype: nodegroup.settings.InstanceType,
		}
	}

	const generateExtensions = (): Install_Extensions => {
		const eks: Install_Extensions_Eks = {
			nodegroups: nodeGroups.map(generateNodeGroup),
			owner: "tom",
			clustername: "cluster",
			projectname: "test"
		};

		return { eks };
	}

	const installSoftware = async () => {
		if (!product) {
			console.error('No product selected for installation');
			return;
		}

		const app: Install_Applications = {
			name: product.Name,
			version: product.Version,
			variables: {}
		};

		const extensions: Install_Extensions = generateExtensions();

		const inst = Install.create();
		inst.applications = app
		inst.extensions = extensions

		install.set(inst);

		const httpHandler = new HttpHandler();
		const writer = Writer.create();
		Install.encode(inst, writer)
		const protomessage = writer.finish();
		const id = await httpHandler.installSoftware(protomessage);
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
