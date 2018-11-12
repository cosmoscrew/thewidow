// Package https implements functions for https server
package https

import (
	"net/http"
	"time"
)

// NewHTTPServer starts a new HTTP Server
func NewHTTPServer(mux *http.ServeMux, IoTimeout int, IdleTimeout int) *http.Server {
	return &http.Server{
		ReadTimeout:  time.Duration(IoTimeout) * time.Second,
		WriteTimeout: time.Duration(IoTimeout) * time.Second,
		IdleTimeout:  time.Duration(IdleTimeout) * time.Second,
		Handler:      mux,
	}
}

// NewRedirectToHTTPSServer redirects the user to https connection
func NewRedirectToHTTPSServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return NewHTTPServer(mux, 5, 120)
}

// NewHTTPSServer creates a new https server
func NewHTTPSServer(mux *http.ServeMux, CertPath string, KeyPath string) (*http.Server, error) {
	err := CheckCertificate(CertPath, KeyPath)
	if err != nil {
		return nil, err
	}

	server := NewHTTPServer(mux, 5, 120)
	server.Addr = ":443"

	go server.ListenAndServeTLS(CertPath, KeyPath)
	return server, nil
}
