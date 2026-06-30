package classifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	clientHTTP "rradar/http"
	modelLLM "rradar/model/llm"
	modelGemini "rradar/model/llm/gemini"
	modelXML "rradar/model/xml"
)

type Gemini25Flash struct {
	apiKey string
}

func NewGemini25Flash(apiKey string) *Gemini25Flash {
    return &Gemini25Flash{
        apiKey: apiKey,
    }
}

func (g Gemini25Flash) Classify(entry modelXML.Entry) (modelLLM.Entry, error) {

	prompt := modelLLM.BuildPrompt(entry.Title, entry.Content)

	body := modelGemini.GenerateContentRequest{
		Contents: []modelGemini.Content{
			{
				Parts: []modelGemini.Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return modelLLM.Entry{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent",
		bytes.NewReader(data),
	)
	if err != nil {
		return modelLLM.Entry{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", g.apiKey)

	resp, err := clientHTTP.Client.Do(req)
	if err != nil {
		return modelLLM.Entry{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return modelLLM.Entry{}, fmt.Errorf("gemini returned %d: %s", resp.StatusCode, b)
	}

	var result modelGemini.GenerateContentResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return modelLLM.Entry{}, err
	}

	if len(result.Candidates) == 0 ||
		len(result.Candidates[0].Content.Parts) == 0 {

		return modelLLM.Entry{}, errors.New("gemini returned no candidates")
	}

	text := result.Candidates[0].Content.Parts[0].Text

	// Parse the LLM's response into your Entry.
	// Assuming BuildPrompt asks Gemini to return JSON.
	var classified modelLLM.Entry

	if err := json.Unmarshal([]byte(text), &classified); err != nil {
		return modelLLM.Entry{}, err
	}

	return classified, nil
}
