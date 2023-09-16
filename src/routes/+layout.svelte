<script lang="ts">
  import "@nasa-jpl/stellar/css/button.css";
  import "@nasa-jpl/stellar/css/index.css";
  import "../custom.scss";
  import { HttpHandler } from "../data/httpHandler";
  import Navbar from "../components/Navbar.svelte";
  import { onMount } from "svelte";
  import { initialized } from "../store/stores";
  import "../app.css";

  onMount(async () => {
    let hasInitialized;

    // Subscribe to the store to get its current value
    initialized.subscribe(value => {
      hasInitialized = value;
    })();

    // If the initialization has not yet run, run it now
    if (!hasInitialized) {
      const httpHandler = new HttpHandler();
      // if (typeof window !== 'undefined') {
      // 	createWebsocketStore('ws://' + window.location.host + '/ws');
      // }
      await httpHandler.setupws();
      // await httpHandler.fetchConfig();

      // Update the store to indicate that the initialization has run
      initialized.set(true);
    }

  });
</script>

<svelte:head>
  <title>Unity Management Console</title>
</svelte:head>
<Navbar brand="Unity Management Console" logoutLink="/logout" />
<slot />

<!--<style lang="scss" global>-->
<!--  @import '../custom.scss';-->
<!--</style>-->
