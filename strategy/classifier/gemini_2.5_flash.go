package classifier

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	clientHTTP "rradar/http"
	adapter "rradar/adapter/classifier"
	"rradar/model"
)

type Gemini25Flash struct {
	apiKey string
}

func NewGemini25Flash(apiKey string) *Gemini25Flash {
    return &Gemini25Flash{
        apiKey: apiKey,
    }
}

func (g Gemini25Flash) Classify(post model.Post) (model.ClassifiedPost, error) {
	a := &adapter.Gemini25FlashAdapter{}
	data, err := a.EncodeRequest(post)
	if err != nil {
		return model.ClassifiedPost{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent",
		bytes.NewReader(data),
	)
	if err != nil {
		return model.ClassifiedPost{}, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", g.apiKey)
	
	resp, err := clientHTTP.Client.Do(req)
	if err != nil {
		return model.ClassifiedPost{}, err
	}
	defer resp.Body.Close()
	// read the data in response body, convert it to []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return model.ClassifiedPost{}, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return model.ClassifiedPost{}, fmt.Errorf("gemini returned %d: %s", resp.StatusCode, b)
	}
	return a.DecodeResponse(data, post)
}
