package main

import (
	"log"
	"net/http"

	"github.com/benkauffman/skwiz-it-api/handler"
	"github.com/benkauffman/skwiz-it-api/middleware"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	listen := "0.0.0.0:3000"

	log.Printf("Starting server and listening on %s", listen)

	router := mux.NewRouter().StrictSlash(true)

	privateBase := mux.NewRouter()
	router.PathPrefix("/api/v1/private").Handler(negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.HandlerFunc(middleware.UserAuthMiddleware),
		negroni.Wrap(privateBase),
	))
	private := privateBase.PathPrefix("/api/v1/private").Subrouter()
	private.Methods("GET").Path("/section/type").HandlerFunc(handler.GetSectionType)
	private.Methods("POST").Path("/section/{type}").HandlerFunc(handler.SaveSection)

	publicBase := mux.NewRouter()
	router.PathPrefix("/api/v1/public").Handler(negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.Wrap(publicBase),
	))
	public := publicBase.PathPrefix("/api/v1/public").Subrouter()
	public.Methods("POST").Path("/register").HandlerFunc(handler.RegisterUser)
	public.Methods("GET").Path("/drawing/{id}").HandlerFunc(handler.GetDrawing)
	public.Methods("GET").Path("/drawings").HandlerFunc(handler.GetDrawings)

	log.Fatal(http.ListenAndServe(listen, router))
}
