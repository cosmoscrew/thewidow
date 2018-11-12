package blindxss

import (
	"log"
	"net/http"
	"testing"

	"github.com/cosmoscrew/thewidow/outputs"
	"github.com/cosmoscrew/thewidow/servers/https"
)

func TestServer(t *testing.T) {
	m := &http.ServeMux{}
	s := https.NewHTTPServer(m, 5, 120)
	o := outputs.New()
	o.SlackWebhook = ""
	s.Addr = ":8082"
	_ = New(m, "http://localhost:8082/m", o)
	log.Fatal(s.ListenAndServe())
}
