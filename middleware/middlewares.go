package middleware

import (
	"../helper"

	"net/http"
)

//expects base64 encoded user information in the header `X-App-User`
//example: eyJuYW1lIjoiQmVuIiwgImVtYWlsIjoiYmVuQGtyYXNoaWRidWlsdC5jb20iLCAiaWQiOiAxfQ==
func UserAuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, err := helper.GetUser(r)

	if err != nil || !user.IsValid() {
		http.Error(w, "Invalid or null user provided", http.StatusUnauthorized)
	} else {
		next(w, r)
	}
}
