package core

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// XSSReport structure contains a single XSS report
type XSSReport struct {
	InnerHTML    string
	URL          string
	Cookie       string
	OpenerCookie string
	OpenerURL    string
	OpenerHTML   string
	IP           string
}

// NewReport creates a new xss report
func NewReport(Values map[string]string) XSSReport {
	return XSSReport{
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
func (x XSSReport) String() string {
	return fmt.Sprintf("Blind XSS Attempt Recieved\nURL: `%s`\tIP: `%s`\nCookies: `%s`\nInnerHTML: \n```\n%s\n```\nOpenerCookie: `%s`\nOpenerURL: `%s`\nOpenerHTML: `%s`\n",
		x.URL, x.IP, x.Cookie, x.InnerHTML, x.OpenerCookie, x.OpenerURL, x.OpenerHTML)
}

// Print prints a xss report
func (x XSSReport) Print() {
	logrus.WithFields(logrus.Fields{
		"innerHTML":    x.InnerHTML,
		"url":          x.URL,
		"cookie":       x.Cookie,
		"openerCookie": x.OpenerCookie,
		"openerUrl":    x.OpenerURL,
		"openerHtml":   x.OpenerHTML,
		"IP":           x.IP,
	}).Info("Blind XSS attempt Recieved")
}
