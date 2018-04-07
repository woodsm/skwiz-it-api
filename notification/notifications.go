package notification

import (
	"log"
	"github.com/benkauffman/skwiz-it-api/config"
	"gopkg.in/mailgun/mailgun-go.v1"
	"strconv"
)

var conf = config.LoadConfig()

func CheckHealth() (bool) {
	mg := mailgun.NewMailgun(conf.MailGun.Domain, conf.MailGun.ApiKey, conf.MailGun.PublicApiKey)
	return mg.ApiKey() == conf.MailGun.ApiKey
}

func SendEmail(emailAddr string, drawingId int64) {

	log.Printf("Mailgun %s for drawing %d", emailAddr, drawingId)
	mg := mailgun.NewMailgun(conf.MailGun.Domain, conf.MailGun.ApiKey, conf.MailGun.PublicApiKey)
	message := mg.NewMessage(
		"notifications@skwiz.it",
		"Your exquisite corpse drawing!",
		`Your drawing has been completed, check it out here:
`+ conf.App.Domain+ "/#/finished?groupId="+ strconv.FormatInt(drawingId, 10),
		emailAddr)
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MailGun ID: %s Resp: %s\n", id, resp)

}
