package notification

import (
	"log"
	"strconv"

	"gopkg.in/mailgun/mailgun-go.v1"
	"bytes"
	"os"
	"io"
	"../helper"
	"strings"
	"../config"
)

var conf = config.LoadConfig()

func CheckHealth() (bool) {
	mg := mailgun.NewMailgun(conf.MailGun.Domain, conf.MailGun.ApiKey, conf.MailGun.PublicApiKey)
	return mg.ApiKey() == conf.MailGun.ApiKey
}

func SendEmail(emailAddr string, drawingId int64) {
	url := conf.App.Domain + "/#/finished?groupId=" + strconv.FormatInt(drawingId, 10)

	log.Printf("Mailgun %s for drawing %d", emailAddr, drawingId)
	mg := mailgun.NewMailgun(conf.MailGun.Domain, conf.MailGun.ApiKey, conf.MailGun.PublicApiKey)
	message := mg.NewMessage(
		"notifications@skwiz.it",
		"Your exquisite corpse drawing!",
		LoadMessage("txt", url),
		emailAddr)

	message.SetHtml(LoadMessage("html", url), )
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MailGun ID: %s Resp: %s\n", id, resp)

}

func LoadMessage(ext string, url string) (string) {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open("./template.email." + ext)
	helper.CheckError(err)
	_, err = io.Copy(buf, f)
	helper.CheckError(err)
	f.Close()
	return strings.Replace(string(buf.Bytes()), "{URL_LINK}", url, 1)
}
