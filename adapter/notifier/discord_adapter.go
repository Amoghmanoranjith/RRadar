package notifier

import (
	"encoding/json"
	"fmt"
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
	return json.Marshal(payload)
}

// this takes []byte returns error if any
