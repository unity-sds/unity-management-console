<script lang="ts">
  import {
    config,
    installError,
    installRunning,
    parametersStore,
    projectStore,
    venueStore
  } from '../../store/stores';
  import ProgressFeedback from '../../components/ProgressFeedback.svelte';
  import { derived, get } from 'svelte/store';
  import { HttpHandler } from '../../data/httpHandler';
  import type { Parameters_Parameter } from '../../data/unity-cs-manager/protobuf/extensions';
  import InputField from '../../components/InputField.svelte';
  import SelectField from '../../components/SelectField.svelte';
  import OptionalParametersAccordion from '../../components/OptionalParametersAccordion.svelte';

  let httpHandler = new HttpHandler();

  let parameters = get(parametersStore); // Get the current value of the store

  let key = '';
  let value = '';
  let list: { key: string; value: string }[] = [];

  $: console.log(config);
  $: console.log($config);

  const venueAndProjectStore = derived(
    [venueStore, projectStore],
    ([$venueStore, $projectStore]) => {
      // Do validations and derivations here
      return {
        venueIsValid: /^[a-z0-9]+$/i.test($venueStore),
        projectIsValid: /^[a-z0-9]+$/i.test($projectStore)
      };
    }
  );

  function handleInputChange(e: CustomEvent) {
    const target = e.detail.target as HTMLInputElement;
    const { id, value } = target;
    if (id === 'project') projectStore.set(value);
    if (id === 'venue') venueStore.set(value);
    // Handle other changes
  }

  function handleSubmit() {
    console.log('project: ' + get(projectStore));
    console.log('venue: ' + get(venueStore));
    const unsubscribe = parametersStore.subscribe((items) => {
      let l = items.parameterlist;
      l['project'] = createBaseParameters_Parameter({
        name: 'project',
        type: 'String',
        insync: true,
        value: get(projectStore),
        tracked: true
      });
      l['venue'] = createBaseParameters_Parameter({
        name: 'venue',
        type: 'String',
        insync: true,
        value: get(venueStore),
        tracked: true
      });
      installRunning.set(true);
      installError.set(false);
      console.log(items.parameterlist);
      httpHandler.updateParameters(items.parameterlist);
    });
    unsubscribe();
  }

  parametersStore.subscribe((value) => {
    console.log('Value: ' + value);
    parameters = value; // Update parameters whenever the store changes

    const venuekey = '/unity/core/venue';
    const projectkey = '/unity/core/project';

    // Check if the value is not null and has the key
    if (
      parameters &&
      parameters.parameterlist &&
      Object.prototype.hasOwnProperty.call(parameters?.parameterlist, projectkey)
    ) {
      projectStore.set(parameters.parameterlist[projectkey].value);
    } else {
      console.log('Key does not exist or parameters is null/undefined.');
    }

    // Check if the value is not null and has the key
    if (
      parameters &&
      parameters.parameterlist &&
      Object.prototype.hasOwnProperty.call(parameters?.parameterlist, venuekey)
    ) {
      venueStore.set(parameters.parameterlist[venuekey].value);
    } else {
      console.log('Key does not exist or parameters is null/undefined.');
    }

    for (const key in parameters.parameterlist) {
      if (key !== '/unity/core/venue' && key !== '/unity/core/project') {
        list = [...list, { key: key, value: parameters.parameterlist[key].value }];
      }
    }
  });

  let privateSubnets: string[] = [];
  let publicSubnets: string[] = [];
  let privateOptions = $config?.networkConfig?.privatesubnets || [];
  let publicOptions = $config?.networkConfig?.publicsubnets || [];
  const transformedPrivateOptions = privateOptions.map((option) => ({
    value: option,
    label: option
  }));
  const transformedPublicOptions = publicOptions.map((option) => ({
    value: option,
    label: option
  }));

  function createBaseParameters_Parameter(p: {
    insync: boolean;
    name: string;
    tracked: boolean;
    type: string;
    value: string;
  }): Parameters_Parameter {
    return { name: p.name, value: p.value, type: p.type, tracked: p.tracked, insync: p.insync };
  }
</script>

<div class="container mx-auto px-4">
  <div class="flex flex-wrap -mx-4">
    <div class="w-full lg:w-1/4 px-4">
      <p class="mb-4">
        Welcome to the Unity Management Console setup wizard. Here we set some values that are
        mandatory for the reliable operation of your Unity platform.
      </p>
      <p class="mb-4">
        During the setup process, you will have the opportunity to configure vital parameters such
        as server resources allocation, system performance thresholds, and logging options. These
        settings are essential for maintaining the stability and optimal performance of your Unity
        platform, allowing you to effectively monitor and manage the health of your Unity
        environment.
      </p>

      <p class="mb-4">
        Additionally, the setup wizard provides an intuitive interface that simplifies the
        configuration process, ensuring that even users with limited technical expertise can easily
        navigate and set up the mandatory values. With clear instructions and helpful tooltips, the
        wizard guides you through each step, minimizing the chances of errors and ensuring a smooth
        setup experience.
      </p>

      <p class="mb-4">
        By investing time and attention in properly configuring these mandatory values through the
        setup wizard, you can establish a robust foundation for your Unity platform, guaranteeing
        its reliability, scalability, and ability to handle the demands of your applications and
        users effectively.
      </p>
    </div>
    <div class="w-full lg:w-3/4 px-4">
      {#if $installRunning === false}
        <form class="space-y-4">
          <InputField
            label="Project Name"
            id="project"
            isValid={$venueAndProjectStore.projectIsValid}
            on:input={handleInputChange}
            subtext="The project managing this Unity environment."
            value={$projectStore}
            disabled={true}
          />
          <InputField
            label="Venue Name"
            id="venue"
            isValid={$venueAndProjectStore.venueIsValid}
            on:input={handleInputChange}
            subtext="The venue this Unity environment is deployed into."
            value={$venueStore}
            disabled={true}
          />

          <SelectField
            label="Private Subnets"
            id="privateSubnets"
            multiple={true}
            options={transformedPrivateOptions}
            bind:value={privateSubnets}
            subtext="Select the private subnets you would like to use for Unity."
          />
          <SelectField
            label="Public Subnets"
            id="publicSubnets"
            multiple={true}
            options={transformedPublicOptions}
            bind:value={publicSubnets}
            subtext="Select the public subnets you would like to use for Unity."
          />

          <OptionalParametersAccordion {list} {key} {value} />
          <button
            type="submit"
            on:click|preventDefault={handleSubmit}
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded mt-5"
            >Save
          </button>
        </form>
      {:else}
        <ProgressFeedback />
      {/if}
    </div>
  </div>
</div>
