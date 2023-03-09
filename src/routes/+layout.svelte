<script lang="ts">
	import '$lib/style/reset.css';
	import '$lib/style/global.scss';

	import Header from '../lib/component/Header.svelte';
	import Footer from '../lib/component/Footer.svelte';
</script>

<Header />
<slot />
<Footer />
