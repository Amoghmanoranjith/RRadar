package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	clientHTTP "rradar/http"
	modelLLM "rradar/model/llm"
)

type Discord struct {
	webhookURL string
}

type Payload struct {
	Content string `json:"content"`
}

func NewDiscord(webhookURL string) *Discord {
	return &Discord{
		webhookURL: webhookURL,
	}
}

// discord has a rate limiting of 5 per 2 seconds

func (d Discord) Notify(entry modelLLM.Entry) (error error) {
	// this gon notify the webhook
	// structure the message
	message := fmt.Sprintf(
		"## 📰 %s\n\n"+
			"**Reason:** %s\n"+
			"**Published:** <t:%d:F>\n"+
			"**Link:** %s",
		entry.Title,
		entry.Reason,
		entry.Published.Unix(),
		entry.Link,
	)

	payload := Payload{
		Content: message,
	}
	data, err := json.Marshal(payload)
	fmt.Println(string(data))
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		d.webhookURL,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = clientHTTP.Client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
