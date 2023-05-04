package app

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-pdf/fpdf"
)

// Resume data

type resume struct {
	TagLine       string
	Experiences   []experience
	Skills        []skill
	Languages     []language
	ExternalLinks map[string]string
	ContactLinks  map[string]string
}

type experience struct {
	Title          string
	Company        string
	From           string
	To             string
	Duration       string
	Description    string
	Location       string
	SkillsAndTools []string
}

type skill struct {
	Title string
	Tools []string
}

type language struct {
	Flag  string
	Name  string
	Level string
}

var resumeData = resume{
	TagLine: "Passionate self-taught software engineer\nexperienced in backend and frontend web development.",
	Experiences: []experience{
		{
			Title:          "Web development tutor",
			Company:        "Orange, Prison de Melun, Mission Locale, Code Phenix, L'Ilot",
			Location:       "Paris, France",
			From:           "January 2023",
			To:             "now",
			Description:    "Taught web development fundamentals with various social programs for (ex-) prisoners and youth at risk.",
			SkillsAndTools: []string{"HTML", "CSS", "JavaScript", "HTTP"},
		},
		{
			Title:          "Backend software engineer",
			Company:        "Canal+",
			Location:       "Paris, France",
			From:           "January 2022",
			To:             "October 2022",
			Description:    "Built video streaming solutions (over DASH and HLS).",
			SkillsAndTools: []string{"Golang", "Docker", "Kubernetes", "PostgreSQL", "Bash", "Gitlab CI", "AWS"},
		},
		{
			Title:          "Freelance software engineer",
			Company:        "Record Eye, Cyclic Studio and other SMBs",
			Location:       "Paris, France",
			From:           "September 2020",
			To:             "January 2022",
			Description:    "Handled frontend and backend web development projects.",
			SkillsAndTools: []string{"Golang", "TypeScript", "Svelte / Vue / React", "HTML", "CSS", "HTTP", "GCP"},
		},
		{
			Title:          "Chief Operations Officer",
			Company:        "Green Online",
			Location:       "Amsterdam, Netherlands",
			From:           "September 2018",
			To:             "April 2020",
			Description:    "Managed the launch and operation of our website services in 5 European countries.",
			SkillsAndTools: []string{"Ruby on Rails", "GCP"},
		},
	},
	Skills: []skill{
		{Title: "Programming languages", Tools: []string{"Golang", "JavaScript / Typescript"}},
		{Title: "Website development", Tools: []string{"HTTP", "HTML", "CSS", "JS", "Svelte / Vue /React", "A11y"}},
		{Title: "DevOps & CI/CD", Tools: []string{"Linux", "Bash", "Ansible", "Gitlab CI / Github Actions", "Docker / Podman", "Kubernetes"}},
		{Title: "Database", Tools: []string{"PostgreSQL", "MongoDB", "SQLite", "BoltDB"}},
		{Title: "SE Practices", Tools: []string{"TDD / BDD", "Clean architecture", "Pair / mob programming"}},
	},
	Languages: []language{
		{Flag: "🇫🇷", Name: "French", Level: "Native"},
		{Flag: "🇬🇧", Name: "English", Level: "Bilingual"},
		{Flag: "🇪🇸", Name: "Spanish", Level: "Working proficiency"},
		{Flag: "🇳🇱", Name: "Dutch", Level: "Basic understanding"},
	},
	ExternalLinks: map[string]string{
		"GitHub":           "https://github.com/ejuju",
		"Personal website": "https://juliensellier.com",
		"Algorithmic art":  "https://instagram.com/algo.croissant",
	},
	ContactLinks: map[string]string{
		"Email address":       "mailto:admin@juliensellier.com",
		"Online contact form": "https://juliensellier.com/contact#form",
	},
}

func generateAndServeResumeFile(resumeData resume) http.HandlerFunc {
	buf := &bytes.Buffer{}
	err := generateResumePDF(buf, resumeData)
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

func generateResumePDF(w io.Writer, resumeData resume) error {
	pdf := fpdf.New("P", "pt", "A4", "")

	// Setup font
	pdf.AddUTF8FontFromBytes("IBMPlexSans", "", mustReadEmbeddedFile(staticFilesFS, "static/IBMPlexSans-Regular.ttf"))
	pdf.AddUTF8FontFromBytes("IBMPlexSans", "B", mustReadEmbeddedFile(staticFilesFS, "static/IBMPlexSans-Bold.ttf"))
	pdf.SetFont("IBMPlexSans", "", normalFontSize)

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
			srcCodeURL := "https://github.com/ejuju/personal_website"
			pdf.Write(normalFontSize+4, "The code used to generate this PDF is available on my GitHub: ")
			setTempFontStyle(pdf, "U", func() { addClickableURL(pdf, srcCodeURL) })
		})
	})

	// Create page 1
	pdf.AddPage()

	// Add title
	pdf.Bookmark("Julien Sellier", 0, -1)
	setTempFontSize(pdf, titleFontSize, func() {
		setTempFontStyle(pdf, "B", func() {
			pdf.MultiCell(0, titleFontSize, "Julien Sellier", "", "C", false)
		})
	})

	// Add sub-title
	pdf.Ln(2 * normalFontSize)
	setTempTextColor(pdf, textDimColor, func() {
		pdf.MultiCell(0, normalFontSize+4, resumeData.TagLine, "", "C", false)
	})

	// Add horizontal line below sub-title
	pdf.Ln(3 * normalFontSize)
	left, _, right, _ := pdf.GetMargins()
	pdf.Rect(left, pdf.GetY(), a4WidthPt-2*right, 0.5, "F")

	// Add experiences
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Experiences", func() {
		for _, exp := range resumeData.Experiences {
			pdf.Bookmark(fmt.Sprintf("%s (%s to %s)", exp.Title, exp.From, exp.To), 2, -1)
			pdf.Ln(2.5 * normalFontSize)

			setTempFontStyle(pdf, "B", func() {
				pdf.MultiCell(0, normalFontSize+4, exp.Title, "", "", false)
			})

			pdf.Ln(0.5 * normalFontSize)

			addKV(pdf, 88, "From", exp.From+" to "+exp.To, midColor, textDimColor, "", "")
			addKV(pdf, 88, "Company", exp.Company, midColor, textDimColor, "", "")
			addKV(pdf, 88, "Location", exp.Location, midColor, textDimColor, "", "")
			addKV(pdf, 88, "Technologies", strings.Join(exp.SkillsAndTools, ", "), midColor, textDimColor, "", "")
			addKV(pdf, 88, "Description", exp.Description, midColor, textDimColor, "", "")
		}

		pdf.AddPage() // move on to page 2 for other sections
	})

	// Add skills
	addSection(pdf, "Skills", func() {
		for _, skill := range resumeData.Skills {
			pdf.Bookmark(skill.Title, 2, -1)

			setTempFontStyle(pdf, "B", func() {
				pdf.Ln(1 * normalFontSize)
				pdf.MultiCell(0, normalFontSize+4, skill.Title, "", "", false)
			})

			setTempTextColor(pdf, textDimColor, func() {
				pdf.Ln(0.25 * normalFontSize)
				pdf.MultiCell(0, normalFontSize+4, strings.Join(skill.Tools, ", "), "", "", false)
			})
		}
	})

	// Add languages
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Languages", func() {
		pdf.Ln(0.75 * normalFontSize)
		for _, lang := range resumeData.Languages {
			pdf.Bookmark(lang.Name, 2, -1)

			pdf.Ln(0.25 * normalFontSize)
			addKV(pdf, 66, lang.Name, lang.Level, textDimColor, midColor, "B", "")
		}
	})

	// Add links
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Links", func() {
		pdf.Ln(0.25 * normalFontSize)
		for text, url := range resumeData.ExternalLinks {
			pdf.Bookmark(text, 2, -1)

			pdf.Ln(0.75 * normalFontSize)
			setTempTextColor(pdf, textDimColor, func() {
				setTempFontStyle(pdf, "B", func() { pdf.CellFormat(106, normalFontSize+4, text+" ", "", 0, "", false, 0, "") })
				setTempFontStyle(pdf, "U", func() { addClickableURL(pdf, url) })
			})
		}
	})

	// Add contact section
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Contact", func() {
		for label, url := range resumeData.ContactLinks {
			pdf.Bookmark(label, 2, -1)

			pdf.Ln(1 * normalFontSize)
			setTempFontStyle(pdf, "B", func() {
				pdf.CellFormat(0, normalFontSize+4, label, "", 1, "", false, 0, "")
			})
			setTempFontStyle(pdf, "U", func() {
				setTempTextColor(pdf, textDimColor, func() {
					addClickableURL(pdf, url)
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
