# Personal website
# My personal website

Online CV and contact form.

### Stack

Stack was kept as simple as possible.
The code simply consists of Go HTTP server that renders HTML templates (using Go's std library html/template).
Zero client-side JS is loaded and no CSS library / pre-processor has been used.

The website pages are rendered on startup (equivalent to SSG).
Global CSS is inlined on top of every page.

HTML template files and other static assets are embedded in the Go binary (using go:embed).

### Resume PDF generation

`resume.go` contains the code used to generate the resume as a PDF.