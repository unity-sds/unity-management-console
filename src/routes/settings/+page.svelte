<script lang="ts">
	// import { TreeView } from 'carbon-components-svelte';

	interface Parameter {
		key: string;
		value: string;
	}

	interface ParameterResponse {
		parameterlist: Record<string, Parameter>;
	}

	interface ParamNode {
		text: string;
		nodes: ParamNode[];
	}

	function parseNodes(rootNode) {
		return Object.entries(rootNode).reduce((acc, [key, param]) => {
			const base = key.split('/')[1];
			const prefixIndex = acc.findIndex((node) => node.text === base);
			if (prefixIndex >= 0) {
				acc[prefixIndex].nodes.push({ text: key, nodes: [] });
			} else {
				acc.push({ text: base, nodes: [{ text: key, nodes: [] }] });
			}
			return acc;
		}, [] as ParamNode[]);
	}

	async function getSSMParams(): Promise<ParameterResponse> {
		const res = await fetch('../api/ssm_params/current', { method: 'GET' });
		if (res.ok) {
			const json = await res.json();
			const nodes = parseNodes(json.parameterlist);
			console.log(nodes);
		}
		return <ParameterResponse>{};
	}
</script>

<div class="container mx-auto px-4">
	<h4>SSM Params Navigator</h4>
	{#await getSSMParams()}
		<strong>Loading</strong>
	{:then res}
		Done!
	{/await}
</div>
