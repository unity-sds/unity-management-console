<script lang="ts">
import { onMount } from 'svelte';
import { goto } from '$app/navigation';
	import { install } from '../../store/stores';
	import { HttpHandler } from '../../data/httpHandler';
	import type { Product } from '../../data/entities';
	let product: Product | null = null;
	let variables: Array<[string, string]> = [];
	let nodeGroups: NodeGroupType[] = [];
	onMount(async () => {
		const productString = localStorage.getItem('product');
		console.log('product:');
		console.log(productString);
		product = productString ? JSON.parse(productString) : null; // Retrieve product from local storage
		if (product != null) {
			variables = Object.entries(product.DefaultDeployment.Variables);
			nodeGroups = product.DefaultDeployment.EksSpec.NodeGroups.map((ng) => {
				let name = Object.keys(ng)[0];
				//return { name: name, ...ng[name] };
				return {
					name: name,
					settings: {
						MinNodes: Number(ng[name].MinNodes),
						MaxNodes: Number(ng[name].MaxNodes),
						DesiredNodes: Number(ng[name].DesiredNodes),
						InstanceType: ng[name].InstanceType
					}
				};
			});
			console.log(nodeGroups);
		}
	});

	let newVariable = { key: '', value: '' };

	function addVariable() {
		variables = [...variables, [newVariable.key, newVariable.value]];
		newVariable.key = '';
		newVariable.value = '';
	}
	function removeVariable(index: number) {
		variables = variables.filter((_, i) => i !== index);
	}
	const installSoftware = async () => {
		console.log('installing');
		const httpHandler = new HttpHandler();
		const id = await httpHandler.installSoftware($install);
		console.log(id);
      goto('/ui/progress', { replaceState: true });

	};
	interface NodeGroupType {
		name: string;
		settings: {
			MinNodes: number;
			MaxNodes: number;
			DesiredNodes: number;
			InstanceType: string;
		};
	}
	let newNodeGroup = {
		name: '',
		settings: {
			MinNodes: 0,
			MaxNodes: 0,
			DesiredNodes: 0,
			InstanceType: ''
		}
	};

	function addNodeGroup() {
		nodeGroups = [...nodeGroups, newNodeGroup];
		newNodeGroup = {
			name: '',
			settings: {
				MinNodes: 0,
				MaxNodes: 0,
				DesiredNodes: 0,
				InstanceType: ''
			}
		};
	}

	function removeNodeGroup(index: number) {
		nodeGroups = nodeGroups.filter((_, i) => i !== index);
	}
</script>

<div class="container">
	<div class="row">
		<div class="col-md-12">
			{#if product}
				<h1 class="my-4">{product.Name} Installation</h1>
				<form on:submit|preventDefault={installSoftware}>
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
					<div class="form-group mt-4">
						<label for="ekscluster"
							>EKS Cluster(Leave blank to deploy new cluster with software)</label
						>
						<select id="ekscluster" class="form-control" />
					</div>
					<div class="accordion mt-4" id="accordionExample">
						<div class="accordion-item">
							<h2 class="accordion-header" id="headingOne">
								<button
									class="accordion-button collapsed"
									type="button"
									data-bs-toggle="collapse"
									data-bs-target="#collapseOne"
									aria-expanded="false"
									aria-controls="collapseOne"
								>
									Advanced EKS Settings
								</button>
							</h2>
							<div
								id="collapseOne"
								class="accordion-collapse collapse"
								aria-labelledby="headingOne"
								data-bs-parent="#accordionExample"
							>
								<div class="accordion-body">
									<h2>Node Groups</h2>
									{#each nodeGroups as nodeGroup, index (nodeGroup.name)}
										<div>
											<h3>{nodeGroup.name}</h3>
											<div class="form-group row">
												<label for={`minNodes${index}`} class="col-sm-2 col-form-label"
													>Min Nodes</label
												>
												<div class="col-sm-10">
													<input
														type="number"
														class="form-control"
														id={`minNodes${index}`}
														bind:value={nodeGroup.settings.MinNodes}
													/>
												</div>
											</div>
											<div class="form-group row">
												<label for={`maxNodes${index}`} class="col-sm-2 col-form-label"
													>Max Nodes</label
												>
												<div class="col-sm-10">
													<input
														type="number"
														class="form-control"
														id={`maxNodes${index}`}
														bind:value={nodeGroup.settings.MaxNodes}
													/>
												</div>
											</div>
											<div class="form-group row">
												<label for={`desiredNodes${index}`} class="col-sm-2 col-form-label"
													>Desired Nodes</label
												>
												<div class="col-sm-10">
													<input
														type="number"
														class="form-control"
														id={`desiredNodes${index}`}
														bind:value={nodeGroup.settings.DesiredNodes}
													/>
												</div>
											</div>
											<div class="form-group row">
												<label for={`instanceType${index}`} class="col-sm-2 col-form-label"
													>Instance Type</label
												>
												<div class="col-sm-10">
													<input
														type="text"
														class="form-control"
														id={`instanceType${index}`}
														bind:value={nodeGroup.settings.InstanceType}
													/>
												</div>
											</div>
											<button
												type="button"
												on:click={() => removeNodeGroup(index)}
												class="btn btn-danger">Remove Node Group</button
											>
										</div>
									{/each}
									<hr />
									<h3>Add Node Group</h3>
									<div class="form-group row">
										<label for="newNodeGroupName" class="col-sm-2 col-form-label">Name</label>
										<div class="col-sm-10">
											<input
												type="text"
												class="form-control"
												id="newNodeGroupName"
												bind:value={newNodeGroup.name}
											/>
										</div>
									</div>
									<div class="form-group row">
										<label for="newNodeGroupMinNodes" class="col-sm-2 col-form-label"
											>Min Nodes</label
										>
										<div class="col-sm-10">
											<input
												type="number"
												class="form-control"
												id="newNodeGroupMinNodes"
												bind:value={newNodeGroup.settings.MinNodes}
											/>
										</div>
									</div>
									<div class="form-group row">
										<label for="newNodeGroupMaxNodes" class="col-sm-2 col-form-label"
											>Max Nodes</label
										>
										<div class="col-sm-10">
											<input
												type="number"
												class="form-control"
												id="newNodeGroupMaxNodes"
												bind:value={newNodeGroup.settings.MaxNodes}
											/>
										</div>
									</div>
									<div class="form-group row">
										<label for="newNodeGroupDesiredNodes" class="col-sm-2 col-form-label"
											>Desired Nodes</label
										>
										<div class="col-sm-10">
											<input
												type="number"
												class="form-control"
												id="newNodeGroupDesiredNodes"
												bind:value={newNodeGroup.settings.DesiredNodes}
											/>
										</div>
									</div>
									<div class="form-group row">
										<label for="newNodeGroupInstanceType" class="col-sm-2 col-form-label"
											>Instance Type</label
										>
										<div class="col-sm-10">
											<input
												type="text"
												class="form-control"
												id="newNodeGroupInstanceType"
												bind:value={newNodeGroup.settings.InstanceType}
											/>
										</div>
									</div>
									<button type="button" on:click={addNodeGroup} class="btn btn-primary"
										>Add Node Group</button
									>
								</div>
							</div>
						</div>
					</div>

					<h2>Variables</h2>
					{#each variables as variable, index (variable[0])}
						<div class="form-group row mt-4">
							<label for={variable[0]} class="col-sm-2 col-form-label">{variable[0]}</label>
							<div class="col-sm-8">
								<input type="text" class="form-control" id={variable[0]} value={variable[1]} />
							</div>
							<div class="col-sm-2">
								<button type="button" on:click={() => removeVariable(index)} class="btn btn-danger st-button secondary
                    large">Remove</button>
							</div>
						</div>
					{/each}
					<div class="form-group row mt-4">
						<div class="col-sm-2">
							<input
								type="text"
								bind:value={newVariable.key}
								class="form-control"
								placeholder="Variable name"
							/>
						</div>
						<div class="col-sm-8">
							<input
								type="text"
								bind:value={newVariable.value}
								class="form-control"
								placeholder="Variable value"
							/>
						</div>
            <div class="col-sm-2">
					      <button type="button" on:click={addVariable} class="btn btn-secondary">Add Variable</button>
            </div>
					</div>
					<button class="btn btn-secondary btn-success mt-3" type="submit">Install</button>
				</form>
			{:else}
				<p>Loading product...</p>
			{/if}
		</div>
	</div>
</div>
