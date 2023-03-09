<script>
	import { page } from '$app/stores';
</script>

<main class="page">
	<h1>{$page.status} {$page.error?.message}</h1>
	<a href="/">Aller à la page d'accueil</a>
</main>

<style>
	h1 {
		margin-bottom: 1rem;
	}
</style>
