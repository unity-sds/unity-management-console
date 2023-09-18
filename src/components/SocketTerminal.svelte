<script lang="ts">
  import { installError, installRunning, messageStore } from "../store/stores";
  import { addLine, handleResize, ourxterm } from "../data/ourxterm";
  import ResizeObserver from "svelte-resize-observer";
  import "xterm/css/xterm.css";


  let installRunningValue: boolean;
  const unsubscribeInstallRunning = installRunning.subscribe(value => {
    installRunningValue = value;
  });

  let installErrorValue: boolean;
  const unsubscribeErrorRunning = installError.subscribe(value => {
    installErrorValue = value;
  });
  let textarea: HTMLTextAreaElement;

  const unsubscribe = messageStore.subscribe((value) => {
    // This code runs whenever messageStore changes
    const allLines = value.split("\n");

    const noEmptyStrings = allLines.filter((str) => str !== "");
    let line = noEmptyStrings[noEmptyStrings.length - 1];

    console.log("sending line: " + line);
    addLine(line + "\n");

  });

  // function handleResize() {
  //   if (termFit) {
  //     termFit.fit();
  //   }
  // }
  // Clean up the subscription when the component is destroyed
  import { onDestroy } from "svelte";

  onDestroy(() => {
    unsubscribe();
    unsubscribeInstallRunning();
    unsubscribeErrorRunning();
  });
</script>
<div>
  <div use:ourxterm={""} />
  <div class="observer">
    <ResizeObserver on:resize={handleResize} />
  </div>
</div>

