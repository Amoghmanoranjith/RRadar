package classifier

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	adapter "rradar/adapter/classifier"
	clientHTTP "rradar/http"
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
		slog.Error(
			"failed to create post request",
			"component", "Gemini25Flash",
			"operation", "Classify",
			"cause", "http.NewRequest",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", g.apiKey)
	
	resp, err := clientHTTP.Client.Do(req)
	if err != nil {
		slog.Error(
			"failed to make post request to gemini 2.5 flash api",
			"component", "Gemini25Flash",
			"operation", "Classify",
			"cause", "clientHTTP.Client.Do",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}
	defer resp.Body.Close()
	// read the data in response body, convert it to []byte
	slog.Info(
    "Gemini Flash response",
    "status", resp.Status,
    "status_code", resp.StatusCode,
)	// read the data in response body, convert it to []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(
			"failed to read from raw response from gemini 2.5 flash api",
			"component", "Gemini25Flash",
			"operation", "Classify",
			"cause", "io.ReadAll",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return model.ClassifiedPost{}, fmt.Errorf("gemini returned %d: %s", resp.StatusCode, b)
	}
	return a.DecodeResponse(data, post)
}
