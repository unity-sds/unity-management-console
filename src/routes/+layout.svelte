<script lang="ts">
	//import 'bootstrap/dist/css/bootstrap.min.css';
        import '@nasa-jpl/stellar/css/button.css';
import '@nasa-jpl/stellar/css/index.css';
import '../custom.scss';
import Navbar from '../components/Navbar.svelte';
</script>

<svelte:head>
	<title>Unity Management Console</title>
</svelte:head>
<Navbar brand="Unity Management Console" logoutLink="/logout" />
<slot />

<style lang="scss" global>
  @import "../custom.scss";
</style>
