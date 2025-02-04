<script lang="ts">
	// import { TreeView } from 'carbon-components-svelte';

	interface Parameter {
		key: string;
		value: string;
	}

	interface ParameterResponse {
		parameterlist: Record<string, Parameter>;
	}

	async function getSSMParams(): Promise<ParameterResponse> {
		const res = await fetch('../api/ssm_params/current', { method: 'GET' });
		if (res.ok) {
			const res = await res.json();
			const nodes = Object.entries(res).reduce((acc, [key, param]) => {
				const base = key.split('/')[0];
				acc.push(base);
				return acc;
			}, []);
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
