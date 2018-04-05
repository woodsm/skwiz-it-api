package middleware

import (
	"net/http"
	"encoding/base64"
	"github.com/benkauffman/skwiz-it-api/model"
	"encoding/json"
)

//expects base64 encoded user information in the header `X-App-User`
//example: eyJuYW1lIjoiQmVuIiwgImVtYWlsIjoiYmVuQGtyYXNoaWRidWlsdC5jb20iLCAiaWQiOiAxfQ==
func UserAuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	bytes, _ := base64.StdEncoding.DecodeString(r.Header.Get("X-App-User"))

	user := new(model.User)
	err := json.Unmarshal(bytes, user)

	if err != nil || !user.IsValid() {
		http.Error(w, "Invalid or null user provided", http.StatusUnauthorized)
	} else {
		next(w, r)
	}
}
