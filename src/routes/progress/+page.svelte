<script>
import { onMount } from 'svelte';
import { fetchline } from './text.ts';
let progress = 0;
let installationComplete = false;
let text = '';
let lines = 0;
let maxLines = 100; 
onMount(() => {
        const interval2 = setInterval(() => {
            lines++;
            text += fetchline(lines)
            // Scroll to the bottom
            const textarea = document.getElementById('console');
            textarea.scrollTop = textarea.scrollHeight;

            if (lines >= maxLines) {
                clearInterval(interval2);
            }
        }, 100);
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
        clearInterval(interval2);
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
            <textarea class="form-control" id="console" rows="10" bind:value={text} readonly></textarea>
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
            <a href="/landing" class="btn btn-primary mt-3">Installation Complete</a>
        {/if}
    </div>
    </div>
</div>
