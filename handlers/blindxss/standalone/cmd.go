// Package blindxss is a tool used for detecting blindxss
package blindxss

import (
	"github.com/sirupsen/logrus"
)

const (
	version = "1.0.0"
	authors = "cosmoscrew"
	website = "https://github.com/cosmoscrew"
)

var (
	SlackWebhook *string
)

func main() {
	logrus.WithFields(logrus.Fields{
		"version": version,
		"authors": authors,
		"website": website,
	}).Info("Starting thewidow-xss agent")

}
