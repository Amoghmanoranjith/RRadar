package model

import "strings"

const Prompt = `
You classify Reddit posts.

Interesting = someone looking to hire a freelance software developer for paid project work.

The author should be looking for someone to build, fix, extend, or maintain software on a paid basis.

Not interesting = traditional employment, company recruiting, staffing agencies, internships, or people looking for jobs.

Input

Title:
{{TITLE}}

Content:
{{CONTENT}}

Good signals:
- looking for a developer
- looking for a programmer
- need a developer
- need a programmer
- hiring a freelancer
- hiring a contractor
- freelance project
- contract work
- paid project
- build my website
- build my web app
- build my mobile app
- build SaaS
- MVP
- custom software
- API integration
- payment integration
- authentication
- bug fixes
- debugging
- maintenance
- feature development
- migrations
- refactoring
- automation
- web scraping
- AI integration
- chatbot
- dashboard
- admin panel
- backend
- frontend
- full stack
- fullstack
- web developer
- software developer
- software engineer
- Go
- Java
- Rust
- Python
- C#
- Node.js
- TypeScript
- JavaScript
- React
- Next.js
- Angular
- Vue
- Express
- NestJS
- FastAPI
- Spring Boot
- Django
- ASP.NET
- SQL
- PostgreSQL
- MySQL
- MongoDB
- Redis
- Docker
- Kubernetes
- AWS
- Azure
- GCP
- REST
- GraphQL
- gRPC
- CI/CD

Not interesting:
- company hiring employees
- full-time job
- part-time employee
- permanent position
- internship
- graduate program
- recruiter posting jobs
- staffing agency recruiting
- hiring for an internal team
- corporate careers post
- open to work
- for hire
- available for work
- looking for a job
- job search
- resume
- CV
- portfolio
- agency advertisement
- looking for clients
- selling own services
- WordPress
- Shopify
- Wix
- Bubble
- Webflow
- SEO
- marketing
- copywriting
- content writing
- video editing
- graphic design
- crypto shilling
- NFT promotion

Rules:
- A client, founder, startup, or business wants to hire a freelancer to build or improve software -> interesting.
- Paid project work, freelance work, consulting, or contract work -> interesting.
- One-off tasks, feature requests, bug fixes, integrations, or maintenance -> interesting.
- Company hiring employees for permanent or internal roles -> not interesting.
- Recruiters or staffing agencies hiring on behalf of companies -> not interesting.
- Author is looking for work -> not interesting.
- Author is advertising themselves or their agency -> not interesting.
- Unpaid collaboration, equity-only offers, hackathons, or volunteer work -> not interesting.
- If it is unclear whether the author wants a freelancer or an employee, classify as not interesting.

Return ONLY:
{"interesting":true,"confidence":0.95,"reason":"<40 words>"}

No markdown.
No extra text.
Only JSON.
`

func BuildPrompt(title, content string) string {
	return strings.NewReplacer(
		"{{TITLE}}", title,
		"{{CONTENT}}", content,
	).Replace(Prompt)
}
