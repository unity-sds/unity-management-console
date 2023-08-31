<script>
export let title = '';
export let description = '';
export let link = '';
export let disabled = false;
</script>

<div class="col-lg-4 col-md-6 mb-4">
    <div class="card h-100">
        <div class="card-body">
            <h4 class="card-title">{title}</h4>
            <p class="card-text">{description}</p>
        </div>
        <div class="card-footer">
            {#if !disabled}
                <a href="{link}" class="btn btn-primary">Go!</a>
            {:else}
                <a href="#" class="btn disabled">Go!</a>
            {/if}
        </div>
    </div>
</div>
