# Personal website

## Tasks

CICD:
- ❗ Build & deploy pipeline
- [ ] E2E testing

Reliability and observability:
- ❗ Rate limiting middleware
- [ ] Record request and response body sizes for report
- [ ] Record response status for report
- [ ] Show 404s in report
- [ ] Improve health report email (HTML template)

DX improvements:
- Try github.com/signintech/gopdf instead of current PDF library.

UX improvements:
- [ ] Improve contact form submission confirmation email (HTML template)
- [ ] Gzip middleware
- [ ] Style customizer (dark/light mode, accent colors, border-radius)
- [ ] Service worker for offline access
- [ ] Custom 404 page
- [ ] Add french translation for website and CV

Legal:
- [ ] Complete website info page

Marketing & SEO:
- [ ] Add JSON+LD and meta OG properties to pages
- [ ] Add image banner for page meta properties (OG / JSON+LD)

## Fun web experiments ideas

- [] Add algorithmic art (p5.js)
- [] Morse code to/from text

```bash
KEY="$(cat ~/.ssh/github_actions_personal_website)"
USERNAME="github"
HOST="juliensellier.com"
./cicd/deploy.sh
```