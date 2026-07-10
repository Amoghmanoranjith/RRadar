package notifier

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"rradar/model"
)

type DiscordAdapter struct{}

type Payload struct {
	Content string `json:"content"`
}

// this takes classifiedPost returns []byte
func (adapter *DiscordAdapter) EncodeRequest(classifiedPost model.ClassifiedPost) ([]byte, error) {
	message := fmt.Sprintf(
		"## 📰 %s\n\n"+
			"**Reason:** %s\n"+
			"**Published:** <t:%d:F>\n"+
			"**Link:** %s",
		classifiedPost.Title,
		classifiedPost.Reason,
		classifiedPost.Published.Unix(),
		classifiedPost.Link,
	)
	payload := Payload{
		Content: message,
	}
	bytes, err := json.Marshal(payload)
	if err != nil{
		slog.Error(
			"failed to marshal discord request",
			"component", "DiscordAdapter",
			"operation", "EncodeRequest",
			"cause", "json.Marshal",
			"error", err,
		)
		return nil, err
	}
	return bytes, nil
}

// this takes []byte returns error if any
