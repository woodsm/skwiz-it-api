package helper

import (
	"net/http"
	"encoding/json"
	"encoding/base64"
	"log"
	"../config"
	"../model"
)

var conf = config.LoadConfig()

func GetUrl(file string) string {
	//https://s3-us-west-2.amazonaws.com/skwiz-it-media/840f12a9-0a59-4ee7-a756-23bb74b8d3db.png
	return "https://s3-" + conf.S3.Region + ".amazonaws.com/" + conf.S3.Bucket + "/" + file
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func Contains(arr [3]string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func WriteJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}

func GetUser(r *http.Request) (*model.User, error) {
	bytes, _ := base64.StdEncoding.DecodeString(r.Header.Get("X-App-User"))

	user := new(model.User)
	err := json.Unmarshal(bytes, user)

	return user, err
}

func GetSections() [3]string {
	return [3]string{"top", "middle", "bottom"}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func GetStringFromNullable(s *string) string {
	if s == nil {
		temp := "" // *string cannot be initialized
		s = &temp  // in one statement
	}
	value := *s // safe to dereference the *string

	return value
}
