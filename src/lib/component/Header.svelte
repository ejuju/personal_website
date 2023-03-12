<script>
	import Logo from './Logo.svelte';
</script>

<header>
	<a href="/" id="brand"><Logo /> Julien Sellier</a>

	<nav>
		<a href="/">Home</a>
		<a href="/cv">CV</a>
		<a href="/blog">Blog</a>
		<a href="/contact">Contact</a>
	</nav>
</header>

<style>
	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem var(--space-page-side);
		border-bottom: 1px solid #222;
	}

	#brand {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: bold;
	}

	#brand :global(svg) {
		height: 1.66rem;
		color: var(--clr-green);
	}

	nav {
		display: flex;
		gap: 1rem;
	}

	nav > a {
		display: flex;
		align-items: center;
	}
</style>
