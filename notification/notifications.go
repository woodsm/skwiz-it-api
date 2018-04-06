package notification

import (
	"log"
	"github.com/benkauffman/skwiz-it-api/database"
	"github.com/benkauffman/skwiz-it-api/config"
	"gopkg.in/mailgun/mailgun-go.v1"
)

var conf = config.LoadConfig()

func SendEmails(drawingId int64) {
	log.Printf("Sending email to users for drawing %d completion", drawingId)
	for _, emailAddr := range database.GetEmailAddresses(drawingId) {
		log.Printf("Sending email in background thread to %s for drawing %d", emailAddr, drawingId)
		go sendEmail(emailAddr, drawingId)
	}
}

func sendEmail(emailAddr string, drawingId int64) {

	log.Printf("Mailgun %s for drawing %d", emailAddr, drawingId)
	mg := mailgun.NewMailgun(conf.MailGun.Domain, config.MailGun.ApiKey, config.MailGun.PublicApiKey)
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
