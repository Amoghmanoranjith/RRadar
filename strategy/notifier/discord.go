package notifier

import (
	"bytes"
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

func (d Discord) Notify(classifiedPost model.ClassifiedPost) (error error) {
	// this gon notify the webhook
	// structure the message
	a := &adapter.DiscordAdapter{}
	data, err := a.EncodeRequest(classifiedPost)
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
