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

type resume struct {
	TagLine          string
	Experiences      []experience
	Skills           []skill
	Languages        []language
	ContactEmailAddr string
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
	ContactEmailAddr: "admin@juliensellier.com",
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
	pdf.AddUTF8FontFromBytes("JetBrainsMono", "", mustReadEmbeddedFile(staticFilesFS, "static/JetBrainsMono-Regular.ttf"))
	pdf.AddUTF8FontFromBytes("JetBrainsMono", "B", mustReadEmbeddedFile(staticFilesFS, "static/JetBrainsMono-Bold.ttf"))
	pdf.SetFont("JetBrainsMono", "", normalFontSize)

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
		pdf.SetTextColor(rgb(midColor))
		defer pdf.SetTextColor(rgb(textColor))
		pdf.Text(marginSideSize+3, a4HeightPt-4*normalFontSize, txt)
		if !isLastPage {
			return
		}
		pdf.Ln(6 * normalFontSize)
		srcCodeURL := "https://github.com/ejuju/personal_website"
		pdf.MultiCell(0, normalFontSize+4, "The code used to generate this PDF is available on my GitHub:", "", "", false)
		pdf.SetFontStyle("U")
		defer pdf.SetFontStyle("")
		pdf.CellFormat(0, normalFontSize+4, srcCodeURL, "", 1, "", false, 0, srcCodeURL)
	})

	// Create page 1
	pdf.AddPage()

	// Add title
	pdf.Bookmark("Julien Sellier", 0, -1)
	pdf.SetFontSize(titleFontSize)
	pdf.SetFontStyle("B")
	pdf.MultiCell(0, titleFontSize, "Julien Sellier", "", "C", false)

	// Add sub-title
	pdf.Ln(2 * normalFontSize)
	pdf.SetFontSize(normalFontSize)
	pdf.SetFontStyle("")
	pdf.SetTextColor(rgb(textDimColor))
	pdf.MultiCell(0, normalFontSize+4, resumeData.TagLine, "", "C", false)
	pdf.SetTextColor(rgb(textColor))

	// Add horizontal line below sub-title
	pdf.Ln(3 * normalFontSize)
	left, _, right, _ := pdf.GetMargins()
	pdf.Rect(left, pdf.GetY(), a4WidthPt-2*right, 0.5, "F")

	// Add experiences
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Experiences", func() {
		for _, exp := range resumeData.Experiences {
			pdf.Bookmark(fmt.Sprintf("%s (%s to %s)", exp.Title, exp.From, exp.To), 2, -1)
			pdf.Ln(2 * normalFontSize)
			pdf.SetFontStyle("B")
			pdf.MultiCell(0, normalFontSize+4, exp.Title, "", "", false)
			pdf.SetFontStyle("")
			pdf.Ln(0.25 * normalFontSize)
			addKV(pdf, 88, "From", exp.From+" to "+exp.To)
			addKV(pdf, 88, "Company", exp.Company)
			addKV(pdf, 88, "Location", exp.Location)
			addKV(pdf, 88, "Technologies", strings.Join(exp.SkillsAndTools, ", "))
			addKV(pdf, 88, "Description", exp.Description)
		}
		pdf.AddPage() // move on to page 2 for other sections
	})

	// Add skills
	addSection(pdf, "Skills", func() {
		for _, skill := range resumeData.Skills {
			pdf.Bookmark(skill.Title, 2, -1)
			pdf.Ln(1 * normalFontSize)
			pdf.SetFontStyle("B")
			pdf.MultiCell(0, normalFontSize+4, skill.Title, "", "", false)
			pdf.SetFontStyle("")
			pdf.Ln(0.25 * normalFontSize)
			pdf.SetTextColor(rgb(textDimColor))
			pdf.MultiCell(0, normalFontSize+4, strings.Join(skill.Tools, ", "), "", "", false)
			pdf.SetTextColor(rgb(textColor))
		}
	})

	// Add languages
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Languages", func() {
		pdf.Ln(0.75 * normalFontSize)
		for _, lang := range resumeData.Languages {
			pdf.Bookmark(lang.Name, 2, -1)
			pdf.Ln(0.25 * normalFontSize)
			pdf.SetFontStyle("B")
			pdf.Write(normalFontSize+4, fmt.Sprintf("%-10s", lang.Name+" "))
			pdf.SetFontStyle("")
			pdf.SetTextColor(rgb(textDimColor))
			pdf.CellFormat(0, normalFontSize+4, lang.Level, "", 1, "", false, 0, "")
			pdf.SetTextColor(rgb(textColor))
		}
	})

	// Add links
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Links", func() {
		for text, url := range map[string]string{
			"GitHub":           "https://github.com/ejuju",
			"Personal website": "https://www.juliensellier.com",
		} {
			pdf.Bookmark(text, 2, -1)
			pdf.Ln(1 * normalFontSize)
			pdf.SetFontStyle("B")
			pdf.CellFormat(0, normalFontSize+4, text, "", 1, "", false, 0, "")
			pdf.SetFontStyle("U")
			pdf.SetTextColor(rgb(textDimColor))
			pdf.CellFormat(0, normalFontSize+4, url, "", 2, "", false, 0, url)
			pdf.SetTextColor(rgb(textColor))
		}
		defer pdf.SetFontStyle("")
	})

	// Add contact section
	pdf.Ln(3 * normalFontSize)
	addSection(pdf, "Contact", func() {
		addLink := func(label, linkText, linkURL string) {
			pdf.Bookmark(label, 2, -1)
			pdf.Ln(1 * normalFontSize)
			pdf.SetFontStyle("B")
			pdf.CellFormat(0, normalFontSize+4, label, "", 1, "", false, 0, "")
			pdf.SetFontStyle("U")
			defer pdf.SetFontStyle("")
			pdf.SetTextColor(rgb(textDimColor))
			defer pdf.SetTextColor(rgb(textColor))
			pdf.CellFormat(0, normalFontSize+4, linkText, "", 2, "", false, 0, linkURL)
		}

		addLink("Email address", resumeData.ContactEmailAddr, "mailto:"+resumeData.ContactEmailAddr)
		addLink("Online contact form", "https://www.juliensellier.com/contact", "https://www.juliensellier.com/contact")
	})

	return pdf.Output(w)
}

func rgb(clr [3]int) (r, g, b int) { return clr[0], clr[1], clr[2] }

func addSection(pdf *fpdf.Fpdf, heading string, cb func()) {
	pdf.Bookmark(heading, 1, -1)
	pdf.SetFontSize(bigFontSize)
	pdf.MultiCell(0, bigFontSize+4, heading, "", "", false)
	pdf.SetFontSize(normalFontSize)
	cb()
}

func addKV(pdf *fpdf.Fpdf, keyCellWidth float64, k, v string) {
	pdf.SetTextColor(rgb(textDimColor))
	pdf.CellFormat(keyCellWidth, normalFontSize+4, k, "", 0, "", false, 0, "")
	pdf.SetTextColor(rgb(textColor))
	pdf.MultiCell(0, normalFontSize+4, v, "", "", false)
}

func mustReadEmbeddedFile(fs embed.FS, fname string) []byte {
	raw, err := fs.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return raw
}
