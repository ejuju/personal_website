<script>
	import { page } from '$app/stores';
	import PageMeta from '$lib/component/PageMeta.svelte';
</script>

<PageMeta title={$page.status + ' ' + $page.error?.message} description="Oops, an error occured." />
<main class="page">
	<section class="intro">
		<h1>
			<span>😓</span>
			{$page.status}
			{$page.error?.message}
		</h1>
		<nav>
			<a href="/">Aller à la page d'accueil</a>
		</nav>
	</section>
</main>

<style>
	h1 {
		margin-bottom: 1rem;
	}
</style>
