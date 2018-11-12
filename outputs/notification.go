package outputs

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Outputs structure contains variables required to send notifications
type Outputs struct {
	SlackWebhook string

	client http.Client
}

// New creates a new output structure
func New() *Outputs {
	o := &Outputs{}

	o.client = http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

	return o
}

// Notify sends a message using whatever methods available
func (o *Outputs) Notify(message string) {
	if o.SlackWebhook != "" {
		o.SendNotificationToSlack(message)
	}
}

// SendNotificationToSlack sends a notification to slack
func (o *Outputs) SendNotificationToSlack(message string) {
	message = strings.Replace(message, "\"", "\\\"", -1)
	jsonData := `{"text": "` + message + `"}`
	req, err := http.NewRequest("POST", o.SlackWebhook, bytes.NewBufferString(jsonData))
	if err != nil {
		logrus.WithError(err).Warn("Could not send notification")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("Could not send notification")
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
}
