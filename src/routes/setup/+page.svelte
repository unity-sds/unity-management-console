<script lang="ts">
	import { projectStore, venueStore } from '../../store/stores';
	import ProgressFeedback from '../../components/ProgressFeedback.svelte';
	import { config, parametersStore } from '../../store/stores';
  import { get } from "svelte/store";
	import { HttpHandler } from "../../data/httpHandler";
	let project = '';
  let venue = '';

	let httpHandler = new HttpHandler()

	let parameters = get(parametersStore); // Get the current value of the store

	parametersStore.subscribe((value) => {
		parameters = value; // Update parameters whenever the store changes

		const venuekey = "/unity/core/venue"
		const projectkey = "/unity/core/project"

		// Check if the value is not null and has the key
		if (parameters && parameters.parameterlist && Object.prototype.hasOwnProperty.call(parameters?.parameterlist, projectkey)) {
			project = parameters.parameterlist[projectkey].value;
		} else {
			console.log("Key does not exist or parameters is null/undefined.");
		}

		// Check if the value is not null and has the key
		if (parameters && parameters.parameterlist && Object.prototype.hasOwnProperty.call(parameters?.parameterlist, venuekey)) {
			venue = parameters.parameterlist[venuekey].value;
		} else {
			console.log("Key does not exist or parameters is null/undefined.");
		}

		for (const key in parameters.parameterlist) {
			if (key !== '/unity/core/venue' && key !== '/unity/core/project') {
				list = [...list, { key: key, value: parameters.parameterlist[key].value }];
			}
		}
	});

  let privateSubnets: string[] = [];
	let publicSubnets: string[] = [];
	let key = '';
	let value = '';
	let list: { key: string; value: string }[] = [];
	$: list;

	let isValid = true;
	let isValidVenue = true;
	// Reactive statement to check if 'project' is alphanumeric
	$: {
		const alphanumeric = /^[a-z0-9]+$/i;
		isValid = alphanumeric.test(project);
	}

	$: {
		const alphanumeric = /^[a-z0-9]+$/i;
		isValidVenue = alphanumeric.test(venue);
	}
	function addToList() {
		list = [...list, { key: key, value: value }];
		key = '';
		value = '';
	}
	function removeFromList(index: number) {
		list = list.filter((_, i) => i !== index);
	}
	function handleSubmit() {
		// Do your save operation here
		// After saving, navigate to /saved
		projectStore.set(project);
		venueStore.set(venue);
		httpHandler.updateParameters()
		//goto('/ui/saved', { replaceState: true });
	}
</script>

<div class="container">
	<div class="row">
		<div class="col-3">
			<p>
				Welcome to the Unity Management Console setup wizard. Here we set some values that are
				mandatory for the reliable operation of your Unity platform.
			</p>
			<p>
				During the setup process, you will have the opportunity to configure vital parameters such
				as server resources allocation, system performance thresholds, and logging options. These
				settings are essential for maintaining the stability and optimal performance of your Unity
				platform, allowing you to effectively monitor and manage the health of your Unity
				environment.
			</p>

			<p>
				Additionally, the setup wizard provides an intuitive interface that simplifies the
				configuration process, ensuring that even users with limited technical expertise can easily
				navigate and set up the mandatory values. With clear instructions and helpful tooltips, the
				wizard guides you through each step, minimizing the chances of errors and ensuring a smooth
				setup experience.
			</p>

			<p>
				By investing time and attention in properly configuring these mandatory values through the
				setup wizard, you can establish a robust foundation for your Unity platform, guaranteeing
				its reliability, scalability, and ability to handle the demands of your applications and
				users effectively.
			</p>
		</div>
		<div class="col">
			{#if $projectStore === ''}
				<form>
					<div class="form-group">
						<label for="project">Project Name</label>
						<input
							type="text"
							class="form-control"
							id="project"
							bind:value={project}
							style={isValid ? '' : 'border-color: red;'}
							required
						/>
						{#if !isValid}
							<div class="invalid-feedback" style="display: block;">
								Project name should be alphanumeric.
							</div>
						{/if}

						<div class="form-text">The project managing this Unity environment.</div>
					</div>

					<div class="form-group mt-4">
						<label for="venue">Venue Name</label>
						<input
							type="text"
							class="form-control"
							id="venue"
							bind:value={venue}
							style={isValidVenue ? '' : 'border-color: red;'}
							required
						/>
						{#if !isValidVenue}
							<div class="invalid-feedback" style="display: block;">
								Project name should be alphanumeric.
							</div>
						{/if}

						<div class="form-text">The venue this Unity environment is deployed into.</div>
					</div>

					<div class="form-group mt-4">
						<label for="privateSubnets">Private Subnets</label>
						<select multiple class="form-control" id="privateSubnets" bind:value={privateSubnets}>
							{#if $config}
								{#if $config.networkConfig}
									{#each $config?.networkConfig?.privatesubnets as option (option)}
										<option value={option}>{option}</option>
									{/each}
								{/if}
							{/if}
						</select>

						<div class="form-text">Select the private subnets you would like to use for Unity.</div>
					</div>

					<div class="form-group mt-4">
						<label for="publicSubnets">Public Subnets</label>
						<select multiple class="form-control" id="publicSubnets" bind:value={publicSubnets}>
							{#if $config}
								{#if $config.networkConfig}
									{#each $config?.networkConfig?.publicsubnets as option (option)}
										<option value={option}>{option}</option>
									{/each}
								{/if}
							{/if}
						</select>
						<div class="form-text">Select the public subnets you would like to use for Unity.</div>
					</div>

					<div class="accordion mt-4" id="optionalParametersAccordion">
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
									Optional Parameters
								</button>
							</h2>
							<div
								id="collapseOne"
								class="accordion-collapse collapse"
								aria-labelledby="headingOne"
								data-bs-parent="#optionalParametersAccordion"
							>
								<div class="accordion-body">
									<div class="form-group">
										<label for="key">Key</label>
										<input type="text" class="form-control" id="key" bind:value={key} />
									</div>

									<div class="form-group">
										<label for="value">Value</label>
										<input type="text" class="form-control" id="value" bind:value />
									</div>

									<button type="button" class="btn btn-primary mt-2" on:click={addToList}
										>Add</button
									>

									<div class="mt-3">
										<ul class="list-group">
											{#each list as item, index}
												<li class="list-group-item">
													{item.key}: {item.value}
													<button
														type="button"
														class="btn btn-danger btn-sm float-right"
														on:click={() => removeFromList(index)}>Remove</button
													>
												</li>
											{/each}
										</ul>
									</div>
								</div>
							</div>
						</div>
					</div>
					<button type="submit" on:click|preventDefault={handleSubmit} class="st-button large mt-5"
						>Save</button
					>
				</form>
			{:else}
				<ProgressFeedback />
			{/if}
		</div>
	</div>
</div>
