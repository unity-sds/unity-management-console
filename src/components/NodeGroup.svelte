<script lang="ts">
  import type { NodeGroupType } from "../data/entities";

  export let nodeGroup : NodeGroupType;
  export let index: number;
  export let onRemove: (index: number) => void;

  // consider defining your constants like '2' and '10' if you use them multiple times in your code
  const LABEL_COL_CLASS = 'col-sm-2 col-form-label';
  const INPUT_COL_CLASS = 'col-sm-10';
</script>

<div>
  <h3>{nodeGroup.name}</h3>
  <div class="form-group row">
    <label for={`minNodes${index}`} class={LABEL_COL_CLASS}>Min Nodes</label>
    <div class={INPUT_COL_CLASS}>
      <input type="number" class="form-control" id={`minNodes${index}`} bind:value={nodeGroup.settings.MinNodes} />
    </div>
  </div>
  <div class="form-group row">
    <label for={`maxNodes${index}`} class="col-sm-2 col-form-label"
    >Max Nodes</label
    >
    <div class="col-sm-10">
      <input
        type="number"
        class="form-control"
        id={`maxNodes${index}`}
        bind:value={nodeGroup.settings.MaxNodes}
      />
    </div>
  </div>
  <div class="form-group row">
    <label for={`desiredNodes${index}`} class="col-sm-2 col-form-label"
    >Desired Nodes</label
    >
    <div class="col-sm-10">
      <input
        type="number"
        class="form-control"
        id={`desiredNodes${index}`}
        bind:value={nodeGroup.settings.DesiredNodes}
      />
    </div>
  </div>
  <div class="form-group row">
    <label for={`instanceType${index}`} class="col-sm-2 col-form-label"
    >Instance Type</label
    >
    <div class="col-sm-10">
      <input
        type="text"
        class="form-control"
        id={`instanceType${index}`}
        bind:value={nodeGroup.settings.InstanceType}
      />
    </div>
  </div>
  <button type="button" on:click={() => onRemove(index)} class="btn btn-danger">Remove Node Group</button>
</div>