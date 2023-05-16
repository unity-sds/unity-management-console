<script>
import { writable } from 'svelte/store';

const gatewayUrl = writable('Not deployed');

function deployGateway() {
    // Here, you would have code to deploy the API Gateway.
    // I'm just setting a dummy URL for demonstration.
    gatewayUrl.set('https://example.com/api-gateway');
}
</script>

<div class="container">
    <div class="row">
        <div class="col-md-12">
            <h1 class="my-4">API Gateway Configuration</h1>
            <button class="btn btn-success mb-4" on:click={deployGateway}>Deploy</button>
            <div class="alert alert-info">
                <strong>Gateway URL:</strong> {$gatewayUrl}
            </div>
        </div>
    </div>
</div>
