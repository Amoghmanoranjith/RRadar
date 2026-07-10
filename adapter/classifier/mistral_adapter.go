package classifier

import (
	"encoding/json"
	"log/slog"
	"rradar/model"
)

type MistralAdapter struct{}

// =========================
// Request
// =========================

type mistralRequest struct {
	Model    string            `json:"model"`
	Messages []mistralMessage  `json:"messages"`
}

type mistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (a *MistralAdapter) EncodeRequest(post model.Post) ([]byte, error) {
	prompt := model.BuildPrompt(post.Title, post.Content)

	req := mistralRequest{
		Model: "mistral-medium-latest",
		Messages: []mistralMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	bytes, err := json.Marshal(req)
	if err != nil {
		slog.Error(
			"failed to marshal mistral request",
			"component", "MistralAdapter",
			"operation", "EncodeRequest",
			"cause", "json.Marshal",
			"error", err,
		)
		return nil, err
	}
	return bytes, err
}

// =========================
// Response
// =========================

type mistralResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}


func (a *MistralAdapter) DecodeResponse(data []byte, post model.Post) (model.ClassifiedPost, error) {
	var resp mistralResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		slog.Error(
			"failed to unmarshal mistral response",
			"component", "MistralAdapter",
			"operation", "DecodeResponse",
			"cause", "json.Unmarshal",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}

	if len(resp.Choices) == 0 {
		return model.ClassifiedPost{}, nil
	}

	var c classification

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &c); err != nil {
		slog.Error(
			"failed to unmarshal mistral response",
			"component", "MistralAdapter",
			"operation", "DecodeResponse",
			"cause", "json.Unmarshal",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}

	return model.ClassifiedPost{
		Title:       post.Title,
		Content:     post.Content,
		Link:        post.Link,
		Author:      post.Author,
		Published:   post.Published,
		Subreddit:   post.Subreddit,
		Interesting: c.Interesting,
		Confidence:  c.Confidence,
		Reason:      c.Reason,
	}, nil
}