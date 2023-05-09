package app

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

// Resume data

type resume struct {
	Name     string
	TagLine  map[lang]string
	PDFTitle string
	// Experiences
	ExperiencesTitle          map[lang]string
	Experiences               []experience
	ExperienceDurationKey     map[lang]string
	ExperienceCompanyKey      map[lang]string
	ExperienceLocationKey     map[lang]string
	ExperienceTechnologiesKey map[lang]string
	ExperienceDescriptionKey  map[lang]string
	ExperienceNow             map[lang]string
	ExperienceMonths          map[lang]string
	// Skills
	SkillsTitle map[lang]string
	Skills      []skill
	// Languages
	LanguagesTitle map[lang]string
	Languages      []language
	// External links
	ExternalLinksTitle map[lang]string
	ExternalLinks      []externalLink
	// Contact
	ContactLinksTitle map[lang]string
	ContactLinks      []contactLink
	// Source code
	SourceCodeText map[lang]string
	SourceCodeURL  string
	GeneratedAt    map[lang]string
}

type experience struct {
	Title          map[lang]string
	Company        string
	From           time.Time
	To             time.Time
	Duration       string
	Description    map[lang]string
	Location       string
	SkillsAndTools []string
}

// Returns the number of months from the start to the end of the work experience.
func (exp *experience) Months() int { return int(exp.To.Sub(exp.From).Hours() / (24 * 30)) }

type skill struct {
	Title map[lang]string
	Tools []string
}

type language struct {
	Flag  string
	Name  map[lang]string
	Level map[lang]string
}

type externalLink struct {
	Label map[lang]string
	URL   string
}

type contactLink struct {
	Label map[lang]string
	URL   string
}

var resumeData = resume{
	Name: "Julien Sellier",
	TagLine: map[lang]string{
		english: "Passionate self-taught software engineer,\nspecialised in backend and frontend web development.",
		french:  "Développeur auto-ditacte passionné,\nspecialisé en développement web (backend et frontend).",
	},
	PDFTitle:                  "My resume",
	ExperiencesTitle:          map[lang]string{english: "Work experience", french: "Expériences"},
	ExperienceDurationKey:     map[lang]string{english: "Duration", french: "Durée"},
	ExperienceCompanyKey:      map[lang]string{english: "Organisation", french: "Organisation"},
	ExperienceLocationKey:     map[lang]string{english: "Location", french: "Lieu"},
	ExperienceTechnologiesKey: map[lang]string{english: "Technologies", french: "Technologies"},
	ExperienceDescriptionKey:  map[lang]string{english: "Description", french: "Description"},
	ExperienceNow:             map[lang]string{english: "now", french: "maintenant"},
	ExperienceMonths:          map[lang]string{english: "months", french: "mois"},
	Experiences: []experience{
		{
			Title: map[lang]string{
				english: "Web development tutor",
				french:  "Formateur en développement web",
			},
			Company:  "Orange, Prison de Melun, Mission Locale, Code Phenix, L'Ilot",
			Location: "Paris, France",
			From:     mustParseTime("January 2023"),
			To:       time.Time{},
			Description: map[lang]string{
				english: "Taught web development fundamentals with various social programs for (ex-) prisoners and youth at risk.",
				french:  "J'ai pu initié et formé des (ex-) détenus et des jeunes en difficulté au fondamentaux du développement web.",
			},
			SkillsAndTools: []string{"HTML", "CSS", "JavaScript", "HTTP"},
		},
		{
			Title: map[lang]string{
				english: "Backend software engineer",
				french:  "Développeur backend",
			},
			Company:  "Canal+",
			Location: "Paris, France",
			From:     mustParseTime("January 2022"),
			To:       mustParseTime("October 2022"),
			Description: map[lang]string{
				english: "Contributed to the development of a new live video streaming solution based on DASH and HLS.",
				french:  "J'ai contribué au développement d'une nouvelle solution de live streaming de vidéo basé sur DASH et HLS.",
			},
			SkillsAndTools: []string{"Golang", "Docker", "Kubernetes", "PostgreSQL", "Bash", "Gitlab CI", "AWS"},
		},
		{
			Title: map[lang]string{
				english: "Freelance web developer",
				french:  "Web développeur freelance",
			},
			Company:  "Record Eye, Cyclic Studio, etc.",
			Location: "Paris, France",
			From:     mustParseTime("September 2020"),
			To:       mustParseTime("January 2022"),
			Description: map[lang]string{
				english: "Handled frontend and backend web development projects.",
				french:  "J'ai géré les projets de développement front et back de plusieurs PMEs",
			},
			SkillsAndTools: []string{"Golang", "TypeScript", "Svelte / Vue / React", "HTML", "CSS", "HTTP", "GCP"},
		},
		{
			Title: map[lang]string{
				english: "Chief Operations Officer",
				french:  "Directeur des opérations",
			},
			Company:  "Green Online",
			Location: "Amsterdam, Netherlands",
			From:     mustParseTime("September 2018"),
			To:       mustParseTime("April 2020"),
			Description: map[lang]string{
				english: "Managed the expansion and operation of our web application in 5 new European countries.",
				french:  "Je me suis occupé de l'expansion et la gestion de notre application web dans 5 nouveaux pays européens.",
			},
			SkillsAndTools: []string{"Ruby on Rails", "GCP"},
		},
	},
	SkillsTitle: map[lang]string{english: "Skills", french: "Compétences"},
	Skills: []skill{
		{
			Title: map[lang]string{english: "Programming languages", french: "Langages de programmation"},
			Tools: []string{"Golang", "JavaScript / Typescript"},
		},
		{
			Title: map[lang]string{english: "Website development", french: "Développement de site web"},
			Tools: []string{"HTTP", "HTML", "CSS", "JS", "Svelte / Vue /React", "A11y"},
		},
		{
			Title: map[lang]string{english: "DevOps & CI/CD", french: "DevOps & CI/CD"},
			Tools: []string{"Linux", "Bash", "Ansible", "Gitlab CI / Github Actions", "Docker / Podman", "Kubernetes"},
		},
		{
			Title: map[lang]string{english: "Database", french: "Bases de données"},
			Tools: []string{"PostgreSQL", "MongoDB", "SQLite", "BoltDB"},
		},
		{
			Title: map[lang]string{english: "SE Practices", french: "Pratiques de développement logiciel"},
			Tools: []string{"TDD / BDD", "Clean architecture", "Pair / mob programming"},
		},
	},
	LanguagesTitle: map[lang]string{english: "Languages", french: "Langues"},
	Languages: []language{
		{
			Flag:  "🇫🇷",
			Name:  map[lang]string{english: "French", french: "Français"},
			Level: map[lang]string{english: "Native", french: "Langue maternelle"},
		},
		{
			Flag:  "🇬🇧",
			Name:  map[lang]string{english: "English", french: "Anglais"},
			Level: map[lang]string{english: "Bilingual", french: "Bilingue"},
		},
		{
			Flag:  "🇪🇸",
			Name:  map[lang]string{english: "Spanish", french: "Espagnol"},
			Level: map[lang]string{english: "Working proficiency", french: "Niveau professionnel"},
		},
		{
			Flag:  "🇳🇱",
			Name:  map[lang]string{english: "Dutch", french: "Néerlandais"},
			Level: map[lang]string{english: "Basic understanding", french: "Compréhension basique"},
		},
	},
	ExternalLinksTitle: map[lang]string{english: "External links", french: "Liens externes"},
	ExternalLinks: []externalLink{
		{
			Label: map[lang]string{english: "GitHub", french: "GitHub"},
			URL:   "https://github.com/ejuju",
		},
		{
			Label: map[lang]string{english: "Website", french: "Site web"},
			URL:   "https://juliensellier.com",
		},
		{
			Label: map[lang]string{english: "Algorithmic art", french: "Art algorithmique"},
			URL:   "https://instagram.com/algo.croissant",
		},
	},
	ContactLinksTitle: map[lang]string{english: "Contact", french: "Contact"},
	ContactLinks: []contactLink{
		{
			Label: map[lang]string{english: "Email address", french: "Adresse email"},
			URL:   "mailto:admin@juliensellier.com",
		},
		{
			Label: map[lang]string{english: "Online contact form", french: "Formulaire de contact en ligne"},
			URL:   "https://juliensellier.com/contact#form",
		},
	},
	SourceCodeText: map[lang]string{
		english: "The code I wrote handle my website and to generate this resume as a PDF is available on my GitHub: ",
		french:  "Le code que j'utilise sur mon site web et pour génerer ce PDF est disponible sur mon GitHub: ",
	},
	SourceCodeURL: "https://github.com/ejuju/personal_website",
	GeneratedAt:   map[lang]string{english: "PDF generated on ", french: "PDF généré le "},
}

func generateAndServeResumeFile(content resume, l lang) http.HandlerFunc {
	buf := &bytes.Buffer{}
	err := generateResumePDF(buf, content, l)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf.Bytes())
	}
}

// PDF generation

const a4WidthPt, a4HeightPt = 595.28, 842.89

var (
	titleFontSize  = 24.0
	bigFontSize    = 16.0
	normalFontSize = 9.5
	marginTopSize  = 80.0
	marginSideSize = 50.0 // left and right margins
	textColor      = [3]int{0, 0, 0}
	textDimColor   = [3]int{50, 50, 50}
	midColor       = [3]int{127, 127, 127}
)

func generateResumePDF(w io.Writer, content resume, l lang) error {
	pdf := fpdf.New("P", "pt", "A4", "")

	// Set metadata
	pdf.SetCreationDate(time.Now())
	pdf.SetAuthor(content.Name, true)
	pdf.SetLang(string(l))
	pdf.SetTitle(content.PDFTitle, true)

	// Setup font
	font := "Roboto"
	pdf.AddUTF8FontFromBytes(font, "", mustReadEmbeddedFile(staticFilesFS, "static/"+font+"-Regular.ttf"))
	pdf.AddUTF8FontFromBytes(font, "B", mustReadEmbeddedFile(staticFilesFS, "static/"+font+"-Bold.ttf"))
	pdf.SetFont(font, "", normalFontSize)

	// Setup default styles
	pdf.SetTopMargin(marginTopSize)
	pdf.SetLeftMargin(marginSideSize)
	pdf.SetRightMargin(marginSideSize)
	pdf.SetTextColor(rgb(textColor))
	pdf.SetFillColor(rgb(textDimColor))

	// Setup footer callback
	pdf.AliasNbPages("{max_page}")
	pdf.SetFooterFuncLpi(func(isLastPage bool) {
		txt := fmt.Sprintf("Page %d/{max_page}", pdf.PageCount())
		setTempTextColor(pdf, midColor, func() {
			pdf.Text(marginSideSize+3, a4HeightPt-4*normalFontSize, txt)

			if !isLastPage {
				return
			}
			pdf.Ln(8 * normalFontSize)
			pdf.Write(normalFontSize+4, content.SourceCodeText[l]+"\n")
			setTempFontStyle(pdf, "U", func() { addClickableURL(pdf, content.SourceCodeURL) })
			pdf.Write(normalFontSize+4, "\n"+content.GeneratedAt[l]+time.Now().Format("02/01/2006 (15:04:05)"))
		})
	})

	// Create page 1
	pdf.AddPage()

	// Add title
	pdf.Bookmark(defaultBranding.Name, 0, -1)
	setTempFontSize(pdf, titleFontSize, func() {
		setTempFontStyle(pdf, "B", func() {
			pdf.MultiCell(0, titleFontSize, defaultBranding.Name, "", "C", false)
		})
	})

	// Add sub-title
	pdf.Ln(2 * normalFontSize)
	setTempTextColor(pdf, textDimColor, func() {
		pdf.MultiCell(0, normalFontSize+4, content.TagLine[l], "", "C", false)
	})

	// Add horizontal line below sub-title
	pdf.Ln(3 * normalFontSize)
	left, _, right, _ := pdf.GetMargins()
	pdf.Rect(left, pdf.GetY(), a4WidthPt-2*right, 0.5, "F")

	// Add experiences
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, content.ExperiencesTitle[l], func() {
		for _, exp := range content.Experiences {
			pdf.Bookmark(fmt.Sprintf("%s (%s)", exp.Title[l], exp.Company), 2, -1)
			pdf.Ln(2.5 * normalFontSize)

			setTempFontStyle(pdf, "B", func() {
				pdf.MultiCell(0, normalFontSize+4, exp.Title[l], "", "", false)
			})

			pdf.Ln(0.5 * normalFontSize)
			fromStr := exp.From.Format("01/2006")
			toStr := content.ExperienceNow[l]
			if !exp.To.IsZero() {
				toStr = exp.To.Format("01/2006")
			}
			dur := fromStr + " - " + toStr
			if !exp.To.IsZero() {
				dur += fmt.Sprintf(" (%d %s)", exp.Months(), content.ExperienceMonths[l])
			}
			addKV(pdf, 88, content.ExperienceDurationKey[l], fromStr+" - "+toStr, midColor, textDimColor, "", "")
			pdf.Ln(0.125 * normalFontSize)
			addKV(pdf, 88, content.ExperienceCompanyKey[l], exp.Company, midColor, textDimColor, "", "")
			pdf.Ln(0.125 * normalFontSize)
			addKV(pdf, 88, content.ExperienceLocationKey[l], exp.Location, midColor, textDimColor, "", "")
			pdf.Ln(0.125 * normalFontSize)
			addKV(pdf, 88, content.ExperienceTechnologiesKey[l], strings.Join(exp.SkillsAndTools, ", "), midColor, textDimColor, "", "")
			pdf.Ln(0.125 * normalFontSize)
			addKV(pdf, 88, content.ExperienceDescriptionKey[l], exp.Description[l], midColor, textDimColor, "", "")
		}

		pdf.AddPage() // move on to page 2 for other sections
	})

	// Add skills
	addSection(pdf, content.SkillsTitle[l], func() {
		for _, skill := range content.Skills {
			pdf.Bookmark(skill.Title[l], 2, -1)

			setTempFontStyle(pdf, "B", func() {
				pdf.Ln(1 * normalFontSize)
				pdf.MultiCell(0, normalFontSize+4, skill.Title[l], "", "", false)
			})

			setTempTextColor(pdf, textDimColor, func() {
				pdf.Ln(0.25 * normalFontSize)
				pdf.MultiCell(0, normalFontSize+4, strings.Join(skill.Tools, ", "), "", "", false)
			})
		}
	})

	// Add languages
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, content.LanguagesTitle[l], func() {
		pdf.Ln(0.75 * normalFontSize)
		for _, lang := range content.Languages {
			pdf.Bookmark(lang.Name[l], 2, -1)

			pdf.Ln(0.25 * normalFontSize)
			addKV(pdf, 66, lang.Name[l], lang.Level[l], textDimColor, midColor, "B", "")
		}
	})

	// Add links
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, content.ExternalLinksTitle[l], func() {
		pdf.Ln(0.25 * normalFontSize)
		for _, link := range content.ExternalLinks {
			pdf.Bookmark(link.Label[l], 2, -1)

			pdf.Ln(0.75 * normalFontSize)
			setTempTextColor(pdf, textDimColor, func() {
				setTempFontStyle(pdf, "B", func() { pdf.CellFormat(106, normalFontSize+4, link.Label[l]+" ", "", 0, "", false, 0, "") })
				setTempFontStyle(pdf, "U", func() { addClickableURL(pdf, link.URL) })
			})
		}
	})

	// Add contact section
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, content.ContactLinksTitle[l], func() {
		for _, link := range content.ContactLinks {
			pdf.Bookmark(link.Label[l], 2, -1)

			pdf.Ln(1 * normalFontSize)
			setTempFontStyle(pdf, "B", func() {
				pdf.CellFormat(0, normalFontSize+4, link.Label[l], "", 1, "", false, 0, "")
			})
			setTempFontStyle(pdf, "U", func() {
				setTempTextColor(pdf, textDimColor, func() {
					addClickableURL(pdf, link.URL)
				})
			})
		}
	})

	return pdf.Output(w)
}

func rgb(clr [3]int) (r, g, b int) { return clr[0], clr[1], clr[2] }

func addSection(pdf *fpdf.Fpdf, heading string, cb func()) {
	pdf.Bookmark(heading, 1, -1)
	pdf.SetFontSize(bigFontSize)
	pdf.SetFontStyle("B")
	pdf.MultiCell(0, bigFontSize+4, heading, "", "", false)
	pdf.SetFontStyle("")
	pdf.SetFontSize(normalFontSize)
	cb()
}

func addKV(pdf *fpdf.Fpdf, keyCellWidth float64, k, v string, kClr, vClr [3]int, kStyle, vStyle string) {
	setTempTextColor(pdf, kClr, func() {
		setTempFontStyle(pdf, kStyle, func() {
			pdf.CellFormat(keyCellWidth, normalFontSize+4, k, "", 0, "", false, 0, "")
		})
	})
	setTempTextColor(pdf, vClr, func() {
		setTempFontStyle(pdf, vStyle, func() {
			pdf.MultiCell(0, normalFontSize+4, v, "", "", false)
		})
	})
}

func addClickableURL(pdf *fpdf.Fpdf, url string) {
	urlText := url
	switch {
	case strings.HasPrefix(url, "mailto:"):
		urlText = strings.TrimPrefix(url, "mailto:")
	case strings.HasPrefix(url, "https://"):
		urlText = strings.TrimPrefix(url, "https://")
	}
	pdf.CellFormat(0, normalFontSize+4, urlText, "", 2, "", false, 0, url)
}

func mustReadEmbeddedFile(fs embed.FS, fname string) []byte {
	raw, err := fs.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return raw
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("January 2006", s)
	if err != nil {
		panic(err)
	}
	return t
}

func setTempFontStyle(pdf *fpdf.Fpdf, style string, cb func()) {
	pdf.SetFontStyle(style)
	defer pdf.SetFontStyle("")
	cb()
}

func setTempFontSize(pdf *fpdf.Fpdf, size float64, cb func()) {
	pdf.SetFontSize(size)
	defer pdf.SetFontSize(normalFontSize)
	cb()
}

func setTempTextColor(pdf *fpdf.Fpdf, color [3]int, cb func()) {
	pdf.SetTextColor(rgb(color))
	defer pdf.SetTextColor(rgb(textColor))
	cb()
}
