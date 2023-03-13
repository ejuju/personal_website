<script>
	import PageMeta from '$lib/component/PageMeta.svelte';
	import { content } from '$lib/content/content';
</script>

<PageMeta title="CV" description="Julien Sellier's CV" />
<main class="page">
	<section class="intro">
		<h1>
			<span>📜</span>
			Curriculum vitae
		</h1>
		<p>
			I'm a passionate self-taught software engineer<br />
			experienced in backend and frontend web development.
		</p>
		<nav>
			<a href="/contact">Get in touch</a>
		</nav>
	</section>

	<section id="experiences" class="tile">
		<h2>Experiences</h2>

		{#each content.experiences as experience}
			<section>
				<h3>{experience.title} at {experience.company}</h3>
				<p class="Description"><span>ℹ️</span>{experience.description}</p>
				<p class="Duration"><span>🗓️</span>{experience.when}</p>
				<p class="Location"><span>🌍</span>{experience.location}</p>
				<div class="Tools">
					<span>⚒️</span>
					<ul class="inlinelist">
						{#each experience.stack as skill}
							<li>{skill}</li>
						{/each}
					</ul>
				</div>
			</section>
		{/each}
	</section>

	<section id="skills" class="tile">
		<h2>Main skills</h2>

		{#each content.skills as skill}
			<section>
				<h3>{skill.title}</h3>
				<ul class="inlinelist">
					{#each skill.tools as tool}
						<li>{tool}</li>
					{/each}
				</ul>
			</section>
		{/each}
	</section>

	<section id="languages" class="tile">
		<h2>Spoken languages</h2>

		{#each content.languages as language}
			<section>
				<span>{language.flag}</span>
				<div>
					<h3>{language.name}</h3>
					<p>{language.level}</p>
				</div>
			</section>
		{/each}
	</section>
</main>

<style lang="scss">
	.tile {
		margin-top: 5rem;
	}

	#skills,
	#experiences,
	#languages {
		h2 {
			text-align: center;
			padding-bottom: 2rem;
			border-bottom: 1px solid #f0f0f030;
		}

		& > section {
			margin-top: 2rem;
		}
	}

	#skills {
		& > section {
			padding: 2rem 1rem;
			background-color: var(--clr-bg-1);
			border-left: 1px solid var(--clr-green);

			ul {
				margin-top: 1rem;
			}
		}
	}

	#experiences {
		& > section {
			padding: 2rem 1rem;
			background-color: var(--clr-bg-1);
			border-left: 1px solid var(--clr-green);

			.Location,
			.Description,
			.Duration,
			.Tools {
				display: flex;
				align-items: baseline;
				gap: 0.75rem;
				margin-top: 1rem;
			}
		}
	}

	#languages {
		section {
			display: flex;
			align-items: baseline;
			gap: 1rem;
			span {
				font-size: 1.25rem;
			}
		}
	}
</style>
