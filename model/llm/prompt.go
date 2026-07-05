package llm

import "strings"

var Prompt = `
You classify Reddit posts.

Interesting = people hiring backend dev.

Input

Title:
{{TITLE}}

Content:
{{CONTENT}}

Good:
- backend
- API
- Go Java Rust Python C# Node
- microservice
- distributed
- SQL NoSQL
- Redis
- Kafka RabbitMQ
- Docker
- Kubernetes
- AWS GCP Azure
- auth
- payment
- performance
- scaling
- bug fix
- refactor
- migration
- integration

Bad:
- frontend
- React Next Angular Vue
- HTML CSS
- WordPress Shopify Wix Bubble Webflow
- mobile only
- UI UX
- SEO
- marketing
- writing
- video
- crypto

Bad:
- for hire
- open to work
- available
- resume
- portfolio
- agency ad
- looking for client
- looking for job
- selling own service

Rule:
- hiring backend -> interesting
- backend > frontend -> interesting
- frontend >= backend -> not interesting
- author want job -> not interesting
- author hire people -> interesting

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