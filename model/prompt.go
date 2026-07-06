package model

import "strings"


const Prompt = `
You classify Reddit posts.

Interesting = people hiring software engineers for backend, frontend, or full-stack work.

Input

Title:
{{TITLE}}

Content:
{{CONTENT}}

Good:
- backend
- frontend
- full stack
- fullstack
- software engineer
- web developer
- API
- Go
- Java
- Rust
- Python
- C#
- Node
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
- microservices
- distributed systems
- SQL
- NoSQL
- PostgreSQL
- MySQL
- MongoDB
- Redis
- Kafka
- RabbitMQ
- Docker
- Kubernetes
- AWS
- GCP
- Azure
- authentication
- authorization
- payments
- performance
- scaling
- debugging
- bug fixing
- refactoring
- migrations
- integrations
- CI/CD
- GraphQL
- REST
- gRPC

Bad:
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

Bad:
- for hire
- open to work
- available for work
- resume
- CV review
- portfolio
- agency advertisement
- looking for client
- looking for job
- selling own service

Rules:
- author is hiring developers or engineers -> interesting
- backend, frontend, or full-stack hiring -> interesting
- contract, freelance, internship, or full-time hiring -> interesting
- author is seeking a job -> not interesting
- author is advertising themselves or their agency -> not interesting
- if unclear whether the post is hiring, prefer not interesting

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
