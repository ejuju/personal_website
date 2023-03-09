<script lang="ts">
	export let title: string;
	export let description: string;
</script>

<svelte:head>
	<title>{title} - Julien Sellier</title>
	<meta name="robots" content="index, follow" />
	<meta name="description" content={description} />
	<meta name="author" content="Julien Sellier" />
</svelte:head>
