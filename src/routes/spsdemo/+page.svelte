<script lang="ts">
import ProgressFeedback from "../../components/ProgressFeedback.svelte";
import { HttpHandler } from "../../data/httpHandler";
let meta = `{
	"metadata": {
		"metadataversion": "unity-cs-0.1",
		"exectarget": "act",
		"deploymentname": "managementdashboard",
		"services": [
			{"name":"ryantestdeploy","source":"unity-sds/unity-sps-prototype","version":"xxx","branch":"main"}
		],
		"extensions":{
			"kubernetes":{
				"clustername":"unity-sps-managementdashboard",
				"owner":"ryan",
				"projectname":"testproject",
				"nodegroups":{
					"group1": {
						"instancetype": "m5.xlarge",
						"nodecount":"1"
					}
				}
			}
		}
	}
}`
function installSPS() {
  let socket = new HttpHandler();
  socket.installSoftwareSocket(meta);
}

</script>

<div class="container">

  <button type="button" class="btn btn-primary mt-2" on:click={installSPS}>Install SPS</button>
  <ProgressFeedback/>

</div>
