package app

type branding struct {
	Name string
	Font string
}

var defaultBranding = branding{
	Name: "Julien Sellier",
	Font: "JetBrainsMono",
}
