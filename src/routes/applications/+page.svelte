<script lang="ts">
import { writable } from 'svelte/store';
import { projectStore, venueStore } from '../../store/stores';
import ApplicationPanelItem from '../../components/ApplicationPanelItem.svelte';
const page = writable('');
let project: string;
let venue: string
projectStore.subscribe(value => {
    project = value;
});
venueStore.subscribe(value => {
    venue = value;
})

let setuprun = true;
$: {
    // If projectStore is not null or an empty string, set setuprun to false
    if ($projectStore && $projectStore.trim() !== '') {
        setuprun = false;
        console.log("store set")
    } else {
        setuprun = true;
        console.log("store not set")
    }
}
$: cardData = [
    {title: 'test_sps', description: 'Unity SPS', link: '/ui/applications/unity-sps/test_sps/explore', disabled:false},
    {title: 'test_deployment', description: 'Unity EKS.', link: '/ui/applications/unity-eks/test_deployment/explore', disabled:false},

]
</script>
<header class="bg-primary text-white text-center py-5 mb-5">
    <h1>Installed Applications</h1>
</header>
<div class="container">
    <div class="row text-center mt-5">
        {#each cardData as card (card.title)}
            <ApplicationPanelItem title={card.title} description={card.description} link={card.link} />
        {/each}
    </div>
</div>
