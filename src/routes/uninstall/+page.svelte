<script lang="ts">

  // This function will run once the component is mounted to the DOM
  import { HttpHandler } from "../../data/httpHandler";
  import ProgressFeedback from "../../components/ProgressFeedback.svelte";

  function handleButtonClick() {
    const httpHandler = new HttpHandler();

    httpHandler.uninstallAllSoftware();
  }


</script>

<div>
  <h1>Uninstalling full stack</h1>
  <button on:click={handleButtonClick}>Go!</button>
</div>

<ProgressFeedback />
