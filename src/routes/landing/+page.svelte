<script lang="ts">
import { writable } from 'svelte/store';
import { projectStore, venueStore } from '../../store/stores';
import ControlPanelItem from '../../components/ControlPanelItem.svelte';
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
    {title: 'Core Management', description: 'Manage your core settings and features.', link: '/setup', disabled: false},
    {title: 'Unity Marketplace', description: 'Explore the Unity Marketplace.', link: '/marketplace', disabled: setuprun},
    {title: 'Application Management', description: 'Manage your applications.', link: '/applications', disabled: setuprun},
    {title: 'Extension Management', description: 'Manage your hosted extensions.', link: '#', disabled: setuprun}
]
</script>
<header class="bg-primary text-white text-center py-5 mb-5">
    <h1>Welcome to the Unity Management Console</h1>
    <p class="lead">Control Panel</p>
</header>
<div class="container">
    <div class="row justify-content-md-center">
        <div class="col col-md-auto">
            {#if $projectStore}
                <ul class="list-group">
                    <li class="list-group-item">Project: {project}</li>
                    <li class="list-group-item">Venue: {venue}</li>
                </ul>
            {:else}
                <div>
            <h5>Setup has not been run, please go to Core Management</h5>
                </div>
            {/if}
        </div>
    </div>
    <div class="row text-center mt-5">
        {#each cardData as card (card.title)}
            <ControlPanelItem title={card.title} description={card.description} link={card.link} disabled={card.disabled} />
        {/each}
    </div>
</div>
