<script lang="ts">
  import { onMount } from "svelte";
  import NodeGroup from './NodeGroup.svelte';
  import type {Product, NodeGroupType} from "../data/entities";

  export let product: Product | null;
  export let nodeGroups: NodeGroupType[] ;

  const defaultNodeGroup = () => ({
    name: '',
    settings: {
      MinNodes: 0,
      MaxNodes: 0,
      DesiredNodes: 0,
      InstanceType: ''
    }
  });

  let newNodeGroup = defaultNodeGroup();

  function addNodeGroup() {
    nodeGroups = [...nodeGroups, newNodeGroup];
    newNodeGroup = defaultNodeGroup();
  }

  function removeNodeGroup(index: number) {
    nodeGroups = nodeGroups.filter((_, i) => i !== index);
  }

  onMount(async () => {
    if (product != null) {
      nodeGroups = product.DefaultDeployment.EksSpec.NodeGroups.map((ng) => {
        let name = Object.keys(ng)[0];
        return {
          name: name,
          settings: {
            MinNodes: Number(ng[name].MinNodes),
            MaxNodes: Number(ng[name].MaxNodes),
            DesiredNodes: Number(ng[name].DesiredNodes),
            InstanceType: ng[name].InstanceType
          }
        };
      });
    }
  });
</script>

<div class="accordion mt-4" id="accordionExample">
  <!-- omitted for brevity -->
  <div class="accordion-body">
    <h2>Node Groups</h2>
    {#each nodeGroups as nodeGroup, index (nodeGroup.name)}
      <NodeGroup {nodeGroup} {index} onRemove={removeNodeGroup} />
    {/each}
    <hr/>
    <h3>Add Node Group</h3>
    <div class="form-group row">
      <label for="newNodeGroupName" class="col-sm-2 col-form-label">Name</label>
      <div class="col-sm-10">
        <input
          type="text"
          class="form-control"
          id="newNodeGroupName"
          bind:value={newNodeGroup.name}
        />
      </div>
    </div>
    <div class="form-group row">
      <label for="newNodeGroupMinNodes" class="col-sm-2 col-form-label"
      >Min Nodes</label
      >
      <div class="col-sm-10">
        <input
          type="number"
          class="form-control"
          id="newNodeGroupMinNodes"
          bind:value={newNodeGroup.settings.MinNodes}
        />
      </div>
    </div>
    <div class="form-group row">
      <label for="newNodeGroupMaxNodes" class="col-sm-2 col-form-label"
      >Max Nodes</label
      >
      <div class="col-sm-10">
        <input
          type="number"
          class="form-control"
          id="newNodeGroupMaxNodes"
          bind:value={newNodeGroup.settings.MaxNodes}
        />
      </div>
    </div>
    <div class="form-group row">
      <label for="newNodeGroupDesiredNodes" class="col-sm-2 col-form-label"
      >Desired Nodes</label
      >
      <div class="col-sm-10">
        <input
          type="number"
          class="form-control"
          id="newNodeGroupDesiredNodes"
          bind:value={newNodeGroup.settings.DesiredNodes}
        />
      </div>
    </div>
    <div class="form-group row">
      <label for="newNodeGroupInstanceType" class="col-sm-2 col-form-label"
      >Instance Type</label
      >
      <div class="col-sm-10">
        <input
          type="text"
          class="form-control"
          id="newNodeGroupInstanceType"
          bind:value={newNodeGroup.settings.InstanceType}
        />
      </div>
    </div>
    <button type="button" on:click={addNodeGroup} class="btn btn-primary"
    >Add Node Group</button
    >
  </div>
</div>