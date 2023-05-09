package app

type branding struct {
	Name             string
	ContactEmailAddr string
	Font             string
}

var defaultBranding = branding{
	Name:             "Julien Sellier",
	ContactEmailAddr: "admin@juliensellier.com",
	Font:             "JetBrainsMono",
}
