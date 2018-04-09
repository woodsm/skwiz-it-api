package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"strconv"

	"log"

	"github.com/benkauffman/skwiz-it-api/database"
	"github.com/benkauffman/skwiz-it-api/helper"
	"github.com/benkauffman/skwiz-it-api/model"
	"github.com/benkauffman/skwiz-it-api/notification"
	"github.com/benkauffman/skwiz-it-api/storage"
	"github.com/gorilla/mux"
)

func GetSectionType(w http.ResponseWriter, r *http.Request) {

	bytes, err := json.Marshal(database.GetNeededSection())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.WriteJsonResponse(w, bytes)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db := database.CheckHealth()
	s3 := storage.CheckHealth()
	mg := notification.CheckHealth()

	if db && s3 && mg {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Fatalf("Internal services are not responding as expected: DB = %t, S3 = %t, MG = %t", db, s3, mg)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	user := new(model.User)
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err = database.UpsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.WriteJsonResponse(w, bytes)
}

// data:image/png;base64,iVkhdfjdAjdfirtn=
func SaveSection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeOf := vars["type"]

	if !helper.Contains(helper.GetSections(), typeOf) {
		http.Error(w, "Section type \""+typeOf+"\" is not valid", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		println(err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	base64Str := helper.TrimQuotes(string(body))

	fileId, err := storage.SaveToS3(base64Str[strings.IndexByte(base64Str, ',')+1:])

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	user, _ := helper.GetUser(r)

	d, err := database.SaveSection(user.Id, typeOf, helper.GetUrl(fileId))

	bytes, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.WriteJsonResponse(w, bytes)
}

func GetDrawings(w http.ResponseWriter, r *http.Request) {

	list := database.GetDrawings()

	bytes, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	helper.WriteJsonResponse(w, bytes)
}

func GetDrawing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.ParseInt(strId, 10, 64)

	d, err := database.GetDrawing(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "unable to encode drawing", http.StatusInternalServerError)
		return
	}

	helper.WriteJsonResponse(w, bytes)
}
