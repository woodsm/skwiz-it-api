package notification

import (
	"log"
	"github.com/benkauffman/skwiz-it-api/config"
	"gopkg.in/mailgun/mailgun-go.v1"
)

var conf = config.LoadConfig()

func SendEmail(emailAddr string, drawingId int64) {

	log.Printf("Mailgun %s for drawing %d", emailAddr, drawingId)
	mg := mailgun.NewMailgun(conf.MailGun.Domain, conf.MailGun.ApiKey, conf.MailGun.PublicApiKey)
	message := mg.NewMessage(
		"no-reply@skwiz.it",
		"Fancy subject!",
		"Your drawing is done"+string(drawingId),
		emailAddr)
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MailGun ID: %s Resp: %s\n", id, resp)

}
