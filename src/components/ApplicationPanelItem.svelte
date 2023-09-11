<script>
  import { uninstallApplication } from "../data/httpHandler";

  export let title = "";
  export let description = "";
  export let link = "";
  export let status = "";
  export let appPackage = "";
  export let deployment = "";
  const uninstallApp = () => {
    uninstallApplication(title, appPackage, deployment);
  };
</script>

<div class="col-lg-3 col-md-6 mb-4">
  <div class="card h-100">
    <div class="card-body">
      <h4 class="card-title">{title}</h4>
      <p class="card-text">{status} - {description}</p>
    </div>
    <div class="card-footer">
      <a href="{link}" class="btn btn-primary">Explore</a>
      <a on:click={uninstallApp} class="btn btn-primary">Uninstall</a>
    </div>
  </div>
</div>
