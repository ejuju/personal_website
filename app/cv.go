package app

var cvPageData = map[string]any{
	"Experiences": []struct {
		Title          string
		Company        string
		From           string
		To             string
		Duration       string
		Description    string
		Location       string
		SkillsAndTools []string
	}{
		{
			Title:          "Web development teacher",
			Company:        "Orange, Prison de Melun, Mission Locale, Code Phenix and L'Ilot",
			Location:       "Paris, France",
			From:           "2023-01-14",
			To:             "now",
			Description:    "Did frontend and backend development for SMBs.",
			SkillsAndTools: []string{"HTML", "CSS", "JavaScript", "P5.js"},
		},
		{
			Title:          "Freelance software engineer",
			Company:        "Self-employed",
			Location:       "Paris, France",
			From:           "2020-06-24",
			To:             "now",
			Description:    "Did frontend and backend development for SMBs.",
			SkillsAndTools: []string{"Golang", "TypeScript"},
		},
		{
			Title:          "Backend software engineer",
			Company:        "Canal+",
			Location:       "Paris, France",
			From:           "2022-03-14",
			To:             "2022-09-21",
			Description:    "Built video streaming solutions (over DASH and HLS).",
			SkillsAndTools: []string{"Golang", "Docker", "Kubernetes", "PostgreSQL", "MongoDB", "Bash", "CI/CD", "Gitlab CI"},
		},
		{
			Title:          "Chief Operations Officer",
			Company:        "Green Online",
			Location:       "Amsterdam, Netherlands",
			From:           "2018-09-14",
			To:             "2020-04-24",
			Description:    "Managed the launch and operation of our website services in 5 European countries.",
			SkillsAndTools: []string{"Ruby", "GCP"},
		},
	},
	"Skills": []struct {
		Title string
		Tools []string
	}{
		{
			Title: "Programming languages",
			Tools: []string{"Golang", "JavaScript / Typescript"},
		},
		{
			Title: "Website development",
			Tools: []string{
				"HTML and A11y",
				"CSS",
				"Svelte (SvelteKit)",
				"Vue (Nuxt)",
				"React (Next)",
				"Technical SEO",
			},
		},
		{
			Title: "DevOps",
			Tools: []string{
				"CI/CD",
				"Bash",
				"Ansible",
				"Gitlab CI / Github Actions",
				"Docker / Podman",
				"Kubernetes",
			},
		},
		{
			Title: "Database",
			Tools: []string{"PostgreSQL", "MongoDB", "SQLite"},
		},
		{
			Title: "CMS",
			Tools: []string{"Wordpress", "Strapi", "Pocketbase"},
		},
		{
			Title: "Hosting",
			Tools: []string{"GCP", "AWS", "Vercel", "Scaleway"},
		},
		{
			Title: "OS",
			Tools: []string{"Linux", "OpenBSD"},
		},
		{
			Title: "Creative coding",
			Tools: []string{"P5.js", "Three.js", "Sonic Pi"},
		},
	},
	"Languages": []struct {
		Flag  string
		Name  string
		Level string
	}{
		{
			Name:  "French",
			Flag:  "🇫🇷",
			Level: "Native",
		},
		{
			Name:  "English",
			Flag:  "🇬🇧",
			Level: "Bilingual",
		},
		{
			Name:  "Spanish",
			Flag:  "🇪🇸",
			Level: "Working proficiency",
		},
		{
			Name:  "Dutch",
			Flag:  "🇳🇱",
			Level: "Basic understanding",
		},
	},
}
