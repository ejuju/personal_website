package app

var resumeTmplData = map[string]any{
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
			Title:          "Web development tutor",
			Company:        "Orange, Prison de Melun, Mission Locale, Code Phenix and L'Ilot",
			Location:       "Paris, France",
			From:           "January 2023",
			To:             "now",
			Description:    "Taught web development fundamentals to (ex-) prisoners.",
			SkillsAndTools: []string{"HTML", "CSS", "JavaScript", "HTTP"},
		},
		{
			Title:          "Backend software engineer",
			Company:        "Canal+",
			Location:       "Paris, France",
			From:           "January 2022",
			To:             "October 2022",
			Description:    "Built video streaming solutions (over DASH and HLS).",
			SkillsAndTools: []string{"Golang", "Docker", "Kubernetes", "PostgreSQL", "Bash", "Gitlab CI"},
		},
		{
			Title:          "Freelance software engineer",
			Company:        "Self-employed",
			Location:       "Paris, France",
			From:           "September 2020",
			To:             "now",
			Description:    "Handled frontend and backend web development for SMBs.",
			SkillsAndTools: []string{"Golang", "TypeScript"},
		},
		{
			Title:          "Chief Operations Officer",
			Company:        "Green Online",
			Location:       "Amsterdam, Netherlands",
			From:           "September 2018",
			To:             "April 2020",
			Description:    "Managed the launch and operation of our website services in 5 European countries.",
			SkillsAndTools: []string{"Ruby", "GCP"},
		},
	},
	"Skills": []struct {
		Title string
		Tools []string
	}{
		{Title: "Programming languages", Tools: []string{"Golang", "JavaScript / Typescript"}},
		{Title: "Website development", Tools: []string{"HTML", "CSS", "JS / TypeScript", "Svelte / Vue / React"}},
		{Title: "DevOps & CI/CD", Tools: []string{"Linux", "Bash", "Ansible", "Gitlab CI / Github Actions", "Docker / Podman", "Kubernetes"}},
		{Title: "Database", Tools: []string{"PostgreSQL", "MongoDB", "SQLite", "BoltDB"}},
		{Title: "Creative coding", Tools: []string{"P5.js", "Three.js", "Sonic Pi"}},
		// {Title: "CMS", Tools: []string{"Wordpress", "Strapi", "Pocketbase"}},
		// {Title: "Hosting", Tools: []string{"GCP", "AWS", "Vercel", "Scaleway"}},
		// {Title: "OS", Tools: []string{"Linux", "OpenBSD"}},
	},
	"Languages": []struct {
		Flag  string
		Name  string
		Level string
	}{
		{Flag: "🇫🇷", Name: "French", Level: "Native"},
		{Flag: "🇬🇧", Name: "English", Level: "Bilingual"},
		{Flag: "🇪🇸", Name: "Spanish", Level: "Working proficiency"},
		{Flag: "🇳🇱", Name: "Dutch", Level: "Basic understanding"},
	},
}
