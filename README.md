# Personal website

## Tasks

CICD:
- [ ] E2E testing
- [ ] Downtime monitoring & alerting

Reliability and observability:
- [ ] ❗ Rate limiting middleware (or delegate to caddy if possible)
- [ ] Record request and response body sizes for report
- [ ] Record response status for report
- [ ] Show 404s in report
- [ ] Improve health report email (HTML template)

DX improvements:
- [ ] Try github.com/signintech/gopdf instead of current PDF library.
- [ ] Dynamically generate sitemap
- [ ] Put all user-facing text content in a Go struct

UX improvements:
- [ ] Improve contact form submission confirmation email (use a HTML template?)
- [ ] Style customizer (dark/light mode, accent colors, border-radius)
- [ ] Service worker for offline access
- [ ] Add french translation for website pages (CV is already done)

Legal:
- [ ] Complete website info page

Marketing & SEO:
- [ ] Add JSON+LD and meta OG properties to pages (include image banner)

## Fun web experiments ideas

- [] Add algorithmic art (p5.js)
- [] Morse code to/from text

```bash
KEY="$(cat ~/.ssh/github_actions_personal_website)"
USERNAME="github"
HOST="juliensellier.com"
./cicd/deploy.sh
```