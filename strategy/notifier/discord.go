package notifier

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	adapter "rradar/adapter/notifier"
	clientHTTP "rradar/http"
	"rradar/model"
)

type Discord struct {
	webhookURL string
}

func NewDiscord(webhookURL string) *Discord {
	return &Discord{
		webhookURL: webhookURL,
	}
}

// discord has a rate limiting of 5 per 2 seconds

// this should get classified post
// get the bytes for this built using an adapter
// make the request and return error if any

func (d Discord) Notify(classifiedPost model.ClassifiedPost) error {
	a := &adapter.DiscordAdapter{}

	data, err := a.EncodeRequest(classifiedPost)
	if err != nil {
		return fmt.Errorf("encode discord request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		d.webhookURL,
		bytes.NewReader(data),
	)
	if err != nil {
		slog.Error(
			"failed to create discord webhook request",
			"component", "Discord",
			"operation", "Notify",
			"cause", "http.NewRequest",
			"error", err,
		)
		return fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := clientHTTP.Client.Do(req)
	if err != nil {
		slog.Error(
			"failed to send discord webhook request",
			"component", "Discord",
			"operation", "Notify",
			"cause", "clientHTTP.Client.Do",
			"error", err,
		)
		return fmt.Errorf("send webhook request: %w", err)
	}
	defer resp.Body.Close()

	slog.Info(
		"Discord webhook response",
		"status", resp.Status,
		"status_code", resp.StatusCode,
	)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord returned status %d", resp.StatusCode)
	}

	return nil
}