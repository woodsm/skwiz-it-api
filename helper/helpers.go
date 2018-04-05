package helper

import (
	"github.com/benkauffman/skwiz-it-api/config"
	"net/http"
)

var conf = config.LoadConfig()

func GetUrl(file string) string {
	//https://s3-us-west-2.amazonaws.com/skwiz-it-media/840f12a9-0a59-4ee7-a756-23bb74b8d3db.png
	return "https://s3-" + conf.S3.Region + ".amazonaws.com/" + conf.S3.Bucket + "/" + file
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1: len(s)-1]
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
