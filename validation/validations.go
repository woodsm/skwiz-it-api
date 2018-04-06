package validation

import (
	"github.com/mailgun/mailgun-go"
	"github.com/benkauffman/skwiz-it-api/config"
	"log"
)

var conf = config.LoadConfig()

func IsValidEmail(emailAddr string) (bool) {

	mg := mailgun.NewMailgun(conf.MailGun.Domain, config.MailGun.ApiKey, config.MailGun.PublicApiKey)
	ev, err := mg.ValidateEmail(emailAddr)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return ev.IsValid
}
