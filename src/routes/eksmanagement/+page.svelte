<script>
    let clusterName = "";
    let nodeGroups = [{
        instanceType: "",
        nodeCount: "",
        nodeGroupName: ""
    }];

    function addNodeGroup() {
        nodeGroups = [...nodeGroups, {instanceType: "", nodeCount: "", nodeGroupName: ""}];
    }

    function deployCluster() {
        // Here, you would have code to deploy the EKS cluster.
        // I'm just logging the inputs for demonstration.
        console.log({
            clusterName,
            nodeGroups
        });
    }
</script>

<div class="container">
    <div class="row">
        <div class="col-md-12">
            <h1 class="my-4">Deploy EKS Cluster</h1>
            <form on:submit|preventDefault={deployCluster}>
                <div class="form-group">
                    <label for="cluster-name">Cluster Name</label>
                    <input id="cluster-name" class="form-control" bind:value={clusterName} required>
                </div>
                <h5 class="mt-4">Node Groups</h5>
                {#each nodeGroups as nodeGroup, i (i)}
                    <div class="form-group">
                        <h6>Node Group {i + 1}</h6>
                        <label for={`nodegroup-name-${i}`}>Nodegroup Name</label>
                        <input id={`nodegroup-name-${i}`} class="form-control" bind:value={nodeGroup.nodeGroupName} required>
                        <label for={`instance-type-${i}`}>Instance Type</label>
                        <input id={`instance-type-${i}`} class="form-control" bind:value={nodeGroup.instanceType} required>
                        <label for={`node-count-${i}`}>Node Count</label>
                        <input id={`node-count-${i}`} class="form-control" type="number" bind:value={nodeGroup.nodeCount} required>
                    </div>
                {/each}
                <button class="btn btn-secondary" type="button" on:click={addNodeGroup}>Add Node Group</button>
                <br/>
                <button class="btn btn-primary mt-3" type="submit">Deploy</button>
            </form>
        </div>
    </div>
</div>
