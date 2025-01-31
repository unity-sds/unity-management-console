<script lang="ts">
  import { config, installError, installRunning, parametersStore } from '../../store/stores';
  import ProgressFeedback from '../../components/ProgressFeedback.svelte';
  import { derived, get } from 'svelte/store';
  import { HttpHandler } from '../../data/httpHandler';
  import type { Parameters_Parameter } from '../../data/unity-cs-manager/protobuf/extensions';
  import InputField from '../../components/InputField.svelte';
  import SelectField from '../../components/SelectField.svelte';
  import OptionalParametersAccordion from '../../components/OptionalParametersAccordion.svelte';

  let httpHandler = new HttpHandler();

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
      return await res.json();
    }
    return <ParameterResponse>{};
  }
</script>

<div class="container mx-auto px-4">
  {#await getSSMParams()}
    <strong>Loading...</strong>
  {:then res}
    {#each Object.entries(res.parameterlist) as [key, param]}
      <p><strong>{key}:</strong>&nbsp;{param.value}</p>
    {/each}
  {/await}
</div>
