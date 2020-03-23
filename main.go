package main

import (
	"log"

	s "database/sql"
	h "net/http"

	"github.com/cbringf/proxy/dom"
	"github.com/cbringf/proxy/http"
	"github.com/cbringf/proxy/mysql"

	"github.com/gorilla/mux"
)

func handleRequest(w h.ResponseWriter, r *h.Request) {
	db, _ := s.Open("mysql", "root:sniPer$3@/ml_proxy")
	dbItemService := mysql.ItemService{DB: db}
	proxy := dom.NewItemProxy(
		dbItemService,
		http.ItemService{HTTP: *h.DefaultClient},
		dbItemService,
		db,
	)
	reqInfo := proxy.RequestHandler(w, r)
	proxy.LogRequest(&reqInfo)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items/{id}", handleRequest).Methods("GET")

	log.Fatal(h.ListenAndServe(":8080", router))
}
