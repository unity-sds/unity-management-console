<script>
    import {onDestroy} from "svelte";
    import {HttpHandler} from "../data/httpHandler";
    import { installComplete, messageStore } from "../store/stores";
    let socket = new HttpHandler();
    onDestroy(() => {
        socket.closeSocket();
    });
</script>

<div class="container d-flex align-items-center justify-content-center vh-100">
    <div class="text-center w-100">
        <div class="row text-center">
            <h2>Installing Application Stack</h2>
        </div>
        <div class="row">
            <div class="form-group col-md-12">
                <textarea class="form-control" id="console" rows="10" bind:value={$messageStore} readonly></textarea>

            </div>
        </div>
        <div class="row mt-3">
            <div class="col-md-3"></div>
            <div class="col-md-3"></div>
            {#if installComplete}
                <a href="/ui/landing" class="btn btn-primary mt-3">Installation Complete</a>
            {/if}
        </div>
    </div>
</div>
