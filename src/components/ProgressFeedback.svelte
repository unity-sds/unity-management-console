<script>
    import {onDestroy, onMount} from "svelte";
    import {HttpHandler} from "../data/httpHandler";
    import {messageStore} from "../store/stores";
    let progress = 0;
    let installationComplete=false;
    let socket = new HttpHandler();
    onDestroy(() => {
        socket.closeSocket();
    });
    onMount(() => {
        //socket.installSoftwareSocket();
        const interval = setInterval(() => {
            if (progress < 100) {
                progress += 10;
            }
            if (progress >= 100) {
                clearInterval(interval);
                installationComplete = true;
            }
        }, 1000);
        return () => {
            clearInterval(interval);
        }
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
            <div class="progress col-md-6" style="height: 50px;">
                <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar"
                     style="width: {progress}%;"
                     aria-valuenow={+progress}
                     aria-valuemin={Number(0)}
                     aria-valuemax={Number(100)}
                >
                    {progress}%
                </div>
            </div>
            <div class="col-md-3"></div>
            {#if installationComplete}
                <a href="/ui/landing" class="btn btn-primary mt-3">Installation Complete</a>
            {/if}
        </div>
    </div>
</div>
