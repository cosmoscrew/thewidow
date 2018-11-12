package blindxss

import (
	"log"
	"net/http"
	"testing"

	"github.com/cosmoscrew/thewidow/servers/https"
)

func TestServer(t *testing.T) {
	m := &http.ServeMux{}
	s := https.NewHTTPServer(m, 5, 120)
	s.Addr = ":8082"
	_ = New(m, "http://localhost:8082/m")
	log.Fatal(s.ListenAndServe())
}
