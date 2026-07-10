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

type Mistral struct {
	apiKey string
}

func NewMistral(apiKey string) *Mistral {
	return &Mistral{
		apiKey: apiKey,
	}
}

func (mistral Mistral) Classify(post model.Post) (model.ClassifiedPost, error) {
	a := &adapter.MistralAdapter{}
	data, err := a.EncodeRequest(post)
	if err != nil {
		return model.ClassifiedPost{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.mistral.ai/v1/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		slog.Error(
			"failed to create post request",
			"component", "Mistral",
			"operation", "Classify",
			"cause", "http.NewRequest",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", mistral.apiKey)

	resp, err := clientHTTP.Client.Do(req)
	if err != nil {
		slog.Error(
			"failed to make post request to mistral api",
			"component", "Mistral",
			"operation", "Classify",
			"cause", "clientHTTP.Client.Do",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}

	defer resp.Body.Close()

	slog.Info(
		"Mistral response",
		"status", resp.Status,
		"status_code", resp.StatusCode,
	)

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(
			"failed to read from raw response from mistral api",
			"component", "Mistral",
			"operation", "Classify",
			"cause", "io.ReadAll",
			"error", err,
		)
		return model.ClassifiedPost{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.ClassifiedPost{}, fmt.Errorf("mistral returned %d: %s", resp.StatusCode, string(data))
	}

	return a.DecodeResponse(data, post)
}