<script lang="ts">
	import ScaleOut from "./common/ScaleOut.svelte";
	import Modal from "./common/Modal.svelte";
	import { HttpHandler, reapplyApplication } from "../data/httpHandler";
	import { goto } from "$app/navigation";

	import checkIcon from "../icons/check.svg";

	export let title = "";
	export let description = "";
	export let link = "";
	export let status = "";
	export let appPackage = "";
	export let appName = "";
	export let deployment = "";

	let latestStatus = "";
	$: combinedStatus = latestStatus || status;

	console.log({ appPackage, appName, deployment });

	export let objectnumber = 0;

	let isUninstalling = false;
	const uninstallApp = () => {
		console.log({ appName, appPackage, deployment });
		return;
		isUninstalling = true;
		const httphandler = new HttpHandler();
		console.log("Uninstalling " + appName);
		httphandler.uninstallSoftware(appName, appPackage, deployment);
	};

	let uninstallComplete = false;
	let uninstallInProgress = false;
	let uninstallError = false;
	async function handleUninstall() {
		const url = `../api/uninstall_application/${appName}/version/${deployment}`;
		const res = await fetch(url);
		if (!res.ok) {
			console.warn("Error uninstalling!");
			return;
		}
		console.log("Starting uninstall!");
		uninstallInProgress = true;
	}

	let statusInterval: any = null;
	$: {
		if (uninstallInProgress && !uninstallComplete && !statusInterval) {
			statusInterval = setInterval(async () => {
				const res = await fetch(
					`../api/install_application/status/${appName}/version/${deployment}`,
				);
				if (res.status === 404) {
					uninstallInProgress = false;
					uninstallComplete = true;
					clearInterval(statusInterval);
				}
				if (!res.ok) {
					console.warn("Error getting status!");
					clearInterval(statusInterval);
				}
				const json = await res.json();
				latestStatus = json.Status;
				if (json.Status === "UNINSTALL FAILED") {
					clearInterval(statusInterval);
					uninstallInProgress = false;
					uninstallError = true;
				}
			}, 5000);
		} else {
			clearInterval(statusInterval);
		}
	}

	$: console.log({ combinedStatus });

	const handleKeydown = (event: KeyboardEvent) => {
		if (event.ctrlKey && event.key === objectnumber.toString()) {
			handleUninstall();
		} else if (event.key === objectnumber.toString()) {
			goto(link);
		}
	};

	const reapplyApp = () => {
		console.log("Reapplying: " + title);
		reapplyApplication(title, appName, deployment);
	};

	let showLogs = false;
	let logInterval: any = null;
	let logs = "";
	let selectedLogOption = "";

	async function fetchLogs() {
		const url =
			selectedLogOption === "uninstall"
				? `../api/uninstall_application/logs/${appName}/${deployment}`
				: `../api/install_application/logs/${appName}/${deployment}`;
		const res = await fetch(url);
		if (!res.ok) {
			console.warn("Can't get logs!");
			if (logInterval) clearInterval(logInterval);
			return;
		}

		logs = await res.text();
		if (uninstallComplete && logs) {
			clearInterval(logInterval);
		}
	}

	async function getLogs() {
		if (!selectedLogOption) {
			clearInterval(logInterval);
			showLogs = false;
			return;
		}

		await fetchLogs();
		showLogs = true;

		if (selectedLogOption === "uninstall") {
			logInterval = setInterval((_) => {
				fetchLogs();
			}, 5000);
		}
	}

	$: if (!showLogs && logInterval) {
		clearInterval(logInterval);
	}
</script>

<div class="lg:w-1/3 md:w-1/2 mb-4" style="flex: 0 0 auto">
	<div class="bg-white border rounded shadow-md h-full">
		<div
			style="
				display: flex;
				flex-direction: column;
				align-items: center;
				padding: 5px;
			"
		>
			<span class="st-typography-header">{title}</span>
			<span class="st-typography-bold">Application: {appName}</span>
			<div
				style="display: flex; gap: 10px; margin: 10px; justify-content: center"
			>
				<span class="st-typography-bold">Installation Status:</span>
				{#if combinedStatus === 'COMPLETE'}
				<span class="st-typography-small-caps" style="color: green">Done</span>
				{:else}
				<span class="st-typography-small-caps" style="color: red"
					>{combinedStatus}</span
				>
				{/if}
			</div>
		</div>
		{#if combinedStatus !== 'UNINSTALLED'}
		<div class="p-4 border-t" style="text-align: center">
			<!--       {#if isUninstalling}
        <div style="display: flex; gap: 5px; align-items: center; justify-content: center;">
          <span class="st-typography-medium">Uninstalling....</span>
          <ScaleOut size={20} />
        </div>
      {:else} -->
			<a href="{link}" on:keydown="{handleKeydown}" class="st-button"
				>Explore</a
			>
			<!-- <button on:click={reapplyApp} on:keydown={handleKeydown} class="st-button"
          >Reapply Installation
        </button> -->
			{#if uninstallInProgress}
			<button
				class="st-button tertiary"
				disabled
				style="color: red; margin-top: 5px"
			>
				Uninstalling...
			</button>
			{:else if uninstallError}
			<button
				class="st-button tertiary"
				disabled
				style="color: red; margin-top: 5px"
			>
				Uninstall Error
			</button>
			{:else if !uninstallComplete}
			<button
				on:click="{handleUninstall}"
				on:keydown="{handleKeydown}"
				class="st-button tertiary"
				style="color: red; margin-top: 5px"
			>
				Uninstall
			</button>
			{/if}
			<!--         {#if uninstallInProgress || uninstallComplete || uninstallError}
          <button
            class="st-button secondary"
            style="margin-top: 5px;"
            on:click={() => getLogs(true, true)}
            >Show Uninstall Logs
          </button>
        {/if} -->
			<select
				class="st-select"
				bind:value="{selectedLogOption}"
				style="height: 34px"
				on:change="{getLogs}"
			>
				<option value="">Show Logs:</option>
				<option value="install">Install</option>
				<option value="uninstall">Uninstall</option>
			</select>
			{#if uninstallComplete}
			<div class="st-typography-small-caps">Uninstall Complete!</div>
			{/if}
			<!-- {/if} -->
			<!-- <button class="st-button secondary" style="margin-top: 5px;" on:click={(_) => getLogs()}
          >Show Install Logs
        </button> -->
		</div>
		{/if}
	</div>
</div>

<Modal bind:showModal="{showLogs}">
	<h2 slot="header">
		<span class="st-typography-bold"> Install Logs for {title} </span>
	</h2>

	{#if logs}
	<pre>
  {logs}
</pre
	>
	{/if} </Modal
>>
