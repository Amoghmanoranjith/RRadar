package classifier

import (
	"encoding/json"
	"fmt"
	"rradar/model"
	modelLLM "rradar/model/llm"
)

type Gemini25FlashAdapter struct{}

// convert post to a http body suitable for gemini 2.5 flash*************************************

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}


func (adapter *Gemini25FlashAdapter) EncodeRequest(post model.Post) ([]byte, error) {
	prompt := modelLLM.BuildPrompt(post.Title, post.Content)
	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	return json.Marshal(req)
}

// convert response from gemini 2.5 flash to classifiedPost type*************************************
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type classification struct {
	Interesting bool    `json:"interesting"`
	Confidence  float32 `json:"confidence"`
	Reason      string  `json:"reason"`
}
func (adapter *Gemini25FlashAdapter) ParseResponse(data []byte, post model.Post) (model.ClassifiedPost, error) {
var resp geminiResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		return model.ClassifiedPost{}, err
	}

	if len(resp.Candidates) == 0 ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return model.ClassifiedPost{}, fmt.Errorf("empty Gemini response")
	}

	var c classification

	if err := json.Unmarshal(
		[]byte(resp.Candidates[0].Content.Parts[0].Text),
		&c,
	); err != nil {
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
