<script lang="ts">
	import AdvancedVar from './advanced_var.svelte';

	type JsonValue = string | { [key: string]: JsonValue };
	export let json: { [key: string]: JsonValue };
	export let indent = 0;
	export let path: string[] = [];
	export let editMode: Boolean = false;
</script>

<div>
	{#each Object.entries(json) as [key, value]}
		<div style="padding-left:{indent}px; margin-bottom: 5px;">
			{#if typeof value === 'object'}
				<div class="st-typography-label" style="margin-bottom: 5px;">{key}:</div>
				<AdvancedVar bind:json={value} indent={indent + 10} path={path.concat([key])} {editMode} />
			{:else}
				<div style=" display: flex; align-items: center;">
					<div class="st-typography-label">{key}:&nbsp;</div>
					{#if editMode}
						<input class="st-input" style="margin-left: 5px;" bind:value={json[key]} />
					{:else}
						<div class="st-typography-bold">{value}</div>
					{/if}
				</div>
			{/if}
		</div>
	{/each}
</div>
