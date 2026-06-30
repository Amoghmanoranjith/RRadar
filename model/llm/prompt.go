package llm

import "strings"

var Prompt = `
You classify Reddit freelance posts.

Input:

Title:
{{TITLE}}

Content:
{{CONTENT}}

Interesting if backend is primary.

Accept:
- backend
- APIs
- Go Java Rust Python C# Node
- microservices
- distributed systems
- databases
- SQL NoSQL
- caching
- Redis
- Kafka RabbitMQ
- Docker
- Kubernetes
- AWS GCP Azure
- CI/CD
- auth
- payments
- performance
- scaling
- bug fixes
- refactoring
- migrations
- integrations

Reject:
- frontend-heavy
- React
- Next.js
- Angular
- Vue
- HTML CSS
- WordPress
- Shopify
- Wix
- Bubble
- Webflow
- mobile-only
- UI/UX
- SEO
- marketing
- writing
- video editing
- crypto promotion

If backend > frontend => interesting.
If frontend >= backend => not interesting.

Return only:

{"interesting":true,"confidence":0.93,"reason":"short reason"}

Rules:
- confidence between 0 and 1.
- reason under 40 words.
- NO markdown.
- NO explanations outside the JSON.
`


func BuildPrompt(title, content string) string {
	return strings.NewReplacer(
		"{{TITLE}}", title,
		"{{CONTENT}}", content,
	).Replace(Prompt)
}