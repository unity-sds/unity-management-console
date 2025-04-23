<script lang="ts">
  import "@nasa-jpl/stellar/css/button.css";
  import "@nasa-jpl/stellar/css/index.css";
  import "../custom.scss";
  import { HttpHandler } from "../data/httpHandler";
  import Navbar from "../components/Navbar.svelte";
  import { onMount } from "svelte";
  import { initialized, isLoading } from "../store/stores";
  import "../app.css";
  import Spinner from "../components/Spinner.svelte";

  onMount(async () => {
    let hasInitialized;

    // Subscribe to the store to get its current value
    initialized.subscribe(value => {
      hasInitialized = value;
    })();

    // If the initialization has not yet run, run it now
    if (!hasInitialized) {
      isLoading.set(true);
      const httpHandler = new HttpHandler();
      // Set up WebSocket for other event-based communications
      httpHandler.setupws();
      // Explicitly fetch config using our new HTTP method
      await httpHandler.fetchConfig();

      // Update the store to indicate that the initialization has run
      initialized.set(true);
      isLoading.set(false);
    }

  });
</script>

<svelte:head>
  <title>Unity Management Console</title>
</svelte:head>
{#if $isLoading}
  <Spinner />
{:else}
  <Navbar brand="Unity Management Console" logoutLink="/logout" />
  <slot />
{/if}

<!--<style lang="scss" global>-->
<!--  @import '../custom.scss';-->
<!--</style>-->
