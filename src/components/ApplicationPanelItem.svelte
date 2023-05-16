<script>
export let title = '';
export let description = '';
export let link = '';

</script>

<div class="col-lg-3 col-md-6 mb-4">
    <div class="card h-100">
        <div class="card-body">
            <h4 class="card-title">{title}</h4>
            <p class="card-text">{description}</p>
        </div>
        <div class="card-footer">
                <a href="{link}" class="btn btn-primary">Explore</a>
        </div>
    </div>
</div>
