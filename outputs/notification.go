package outputs

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// Outputs structure contains variables required to send notifications
type Outputs struct {
	SlackWebhook   string
	SendgridAPIKey string
	Email          string
	Client         http.Client
	Sendgrid       *sendgrid.Client
}

// New creates a new output structure
func New(SlackWebhook string, SendgridAPIKey string, Email string) *Outputs {
	o := &Outputs{
		SlackWebhook:   SlackWebhook,
		SendgridAPIKey: SendgridAPIKey,
		Email:          Email,

		Client: http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}},
	}

	if o.SendgridAPIKey != "" {
		o.Sendgrid = sendgrid.NewSendClient(o.SendgridAPIKey)
	}

	return o
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

	resp, err := o.Client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("Could not send notification")
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
}

// SendNotificationToSlackDirect does what it says
func SendNotificationToSlackDirect(message string, webhook string) {
	message = strings.Replace(message, "\"", "\\\"", -1)
	jsonData := `{"text": "` + message + `"}`
	req, err := http.NewRequest("POST", webhook, bytes.NewBufferString(jsonData))
	if err != nil {
		logrus.WithError(err).Warn("Could not send notification")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("Could not send notification")
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
}

// SendNotificationToEmailViaSendgrid does what it says
func (o *Outputs) SendNotificationToEmailViaSendgrid(message string, title string) {
	from := mail.NewEmail("TheWidow", "alerts@thewidow.com")
	to := mail.NewEmail(o.Email, o.Email)

	emailMessage := mail.NewSingleEmail(from, title, to, "", message)

	_, err := o.Sendgrid.Send(emailMessage)
	if err != nil {
		logrus.WithError(err).Warning("Could not send sendgrid notification")
	}
}
