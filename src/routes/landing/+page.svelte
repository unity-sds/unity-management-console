<script lang="ts">
  import { config, parametersStore } from '../../store/stores';
  import type { Config } from '../../data/unity-cs-manager/protobuf/extensions';
  import ControlPanelItem from '../../components/ControlPanelItem.svelte';
  import { onMount } from 'svelte';
  import ExtendedSemver from '../../lib/ExtendedSemver';

  // Initialize ExtendedSemver for version comparisons
  const semver = new ExtendedSemver();

  let conf: Config | null;
  let latestVersion: string | null = null;
  let updateAvailable = false;
  let checkingForUpdates = false;
  let updateError = false;
  let releaseData: any = null; // Store the raw GitHub release data

  // Update process states
  let updating = false;
  let updateSuccess = false;
  let updateErrorMessage = '';

  config.subscribe((value) => {
    conf = value;
    // When config changes and we have the latest version, check if an update is available
    if (value && value.version && latestVersion) {
      try {
        // Use ExtendedSemver for comparison to handle non-standard version formats
        updateAvailable = semver.gt(latestVersion, value.version);
      } catch (e) {
        console.error('Error comparing versions:', e);
        updateAvailable = false;
      }
    }
  });

  onMount(async () => {
    await checkForUpdates();
  });

  async function checkForUpdates() {
    checkingForUpdates = true;
    updateError = false;

    try {
      // Fetch the latest release from GitHub API
      const response = await fetch(
        'https://api.github.com/repos/unity-sds/unity-management-console/releases/latest'
      );

      if (!response.ok) {
        throw new Error(`GitHub API returned ${response.status}`);
      }

      const data = await response.json();
      releaseData = data; // Store the raw release data for UI display

      // Extract version from tag name (e.g., "v1.0.0" -> "1.0.0") using ExtendedSemver's coerce
      const tagName = data.tag_name;
      const coercedVersion = semver.coerce(tagName);

      if (!coercedVersion) {
        throw new Error(`Invalid version format: ${tagName}`);
      }

      latestVersion = coercedVersion;

      // Check if current version is older than latest version
      if (conf && conf.version) {
        try {
          // Use ExtendedSemver for comparison
          updateAvailable = semver.gt(latestVersion, conf.version);
        } catch (e) {
          console.error('Error comparing versions:', e);
          updateAvailable = false;
        }
      }
    } catch (error) {
      console.error('Error checking for updates:', error);
      updateError = true;
    } finally {
      checkingForUpdates = false;
    }
  }

  // Function to trigger the update process
  async function startUpdate() {
    updating = true;
    updateSuccess = false;

    try {
      const response = await fetch('../api/update-management-console', { method: 'POST' });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error ${response.status}: ${errorText}`);
      }

      // Parse the response
      const result = await response.json();

      if (result.success) {
        updateSuccess = true;
        // Refresh the page after 45 seconds when update is successful
        setTimeout(() => {
          window.location.reload();
        }, 1000 * 45);
      } else {
        throw new Error(result.error || 'Unknown error during update');
      }
    } catch (error) {
      console.error('Error during update:', error);
      updateErrorMessage = error instanceof Error ? error.message : String(error);
    } finally {
      updating = false;
    }
  }

  let setuprun: boolean;
  let bootstrapfailed: boolean;
  let bootstrapped: boolean;
  $: {
    if (conf && conf.bootstrap == 'complete') {
      bootstrapped = true;
      bootstrapfailed = false;
    } else if (conf && conf.bootstrap == 'failed') {
      bootstrapped = false;
      bootstrapfailed = true;
    } else if (conf && conf.bootstrap == '') {
      bootstrapped = false;
      bootstrapfailed = false;
    }
    setuprun = !!(conf && conf.updatedby !== '');
  }
  $: cardData = [
    {
      title: 'Unity Marketplace',
      description: 'Explore the Unity Marketplace.',
      link: '/management/ui/marketplace',
      disabled: !setuprun
    },
    {
      title: 'Application Management',
      description: 'Manage your applications.',
      link: '/management/ui/applications',
      disabled: !setuprun
    }
  ];

  $: project = $config ? $config.applicationConfig?.Project : '(loading)';
  $: venue = $config ? $config.applicationConfig?.Venue : '(loading)';
  $: currentVersion = $config ? $config.version : '(loading)';
</script>

<header class="bg-blue-600 text-white text-center py-12 mb-12">
  <h1 class="text-4xl">Welcome to the Unity Management Console</h1>
  <p class="text-lg leading-6">Control Panel</p>
  <div class="mt-4 text-sm">
    <span>Version: {currentVersion}</span>
    {#if checkingForUpdates}
      <span class="ml-2 inline-block animate-pulse">Checking for updates...</span>
    {:else if updateError}
      <span class="ml-2 inline-block text-yellow-300">
        Failed to check for updates.
        <button on:click={checkForUpdates} class="underline ml-1">Retry</button>
      </span>
    {:else if updateAvailable}
      <span class="ml-2 inline-block bg-yellow-500 text-black px-2 py-1 rounded relative group">
        Update available! ({releaseData?.tag_name || `v${latestVersion}`})

        {#if updating}
          <span class="ml-1 inline-flex items-center">
            <svg
              class="animate-spin -ml-1 mr-2 h-4 w-4 text-black"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              />
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
            Updating...
          </span>
        {:else if updateSuccess}
          <span class="ml-1 text-green-900 font-bold"> Update downloaded! Reloading soon... </span>
        {:else if updateError}
          <span class="ml-1 text-red-900 font-bold"> Update failed </span>
          <button
            on:click={startUpdate}
            class="ml-1 bg-blue-700 hover:bg-blue-800 text-white px-2 py-0.5 rounded text-sm"
          >
            Retry
          </button>
        {:else}
          <button
            on:click={startUpdate}
            class="ml-1 bg-blue-700 hover:bg-blue-800 text-white px-2 py-0.5 rounded text-sm"
          >
            Update now
          </button>
          <a
            href="https://github.com/unity-sds/unity-management-console/releases/latest"
            target="_blank"
            rel="noopener noreferrer"
            class="underline ml-1 text-xs"
          >
            View details
          </a>
        {/if}

        {#if releaseData}
          <div
            class="absolute left-0 bottom-full mb-2 hidden group-hover:block bg-gray-900 text-white p-2 rounded shadow-lg z-10 w-64 text-left"
          >
            <div class="font-bold">{releaseData.name || releaseData.tag_name}</div>
            <div class="text-xs max-h-36 overflow-y-auto">
              {#if releaseData.body}
                <p>
                  {releaseData.body.length > 200
                    ? releaseData.body.substring(0, 200) + '...'
                    : releaseData.body}
                </p>
              {/if}
              <p class="mt-1 text-gray-300">
                Published: {new Date(releaseData.published_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        {/if}
      </span>
    {:else if latestVersion}
      <span class="ml-2 inline-block text-green-300">Up to date</span>
    {/if}
  </div>
</header>

<div class="container mx-auto">
  <div class="flex justify-center">
    <div class="flex-initial">
      {#if bootstrapfailed}
        <div class="managementfeedback">
          <h5 class="text-xl">The Bootstrap Process Failed Please Check The Logs</h5>
        </div>
      {:else if !bootstrapped}
        <div class="managementfeedback">
          <h5 class="text-xl">Bootstrap is either in progress or has not been run</h5>
        </div>
      {:else}
        <div class="managementfeedback">
          <ul class="list-decimal pl-5">
            <li class="bg-gray-200 p-4 border-b border-gray-300">
              Project: {project}
            </li>
            <li class="bg-gray-200 p-4">Venue: {venue}</li>
          </ul>
        </div>
      {/if}
    </div>
  </div>

  <div class="flex justify-center mt-12 text-center">
    {#each cardData as card (card.title)}
      <ControlPanelItem
        title={card.title}
        description={card.description}
        link={card.link}
        disabled={card.disabled}
      />
    {/each}
  </div>
</div>
