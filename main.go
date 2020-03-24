package main

import (
	"log"

	"database/sql"
	s "database/sql"
	h "net/http"

	"github.com/cbringf/proxy/dom"
	"github.com/cbringf/proxy/http"
	"github.com/cbringf/proxy/mysql"

	"github.com/gorilla/mux"
)

var db *sql.DB

func handleItemRequest(w h.ResponseWriter, r *h.Request) {
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

func handleHealthRequest(w h.ResponseWriter, r *h.Request) {
	proxy := dom.ItemProxy{
		DB: db,
	}
	requestList, _ := proxy.ReadRequests()
	snapshot := dom.SysRequestSnapshot{
		SysRequestList: dom.BuildSnapshot(requestList),
	}

	snapshot.HandleRequest(w, r)
}

func main() {
	db, _ = s.Open("mysql", "root:sniPer$3@/ml_proxy?parseTime=true")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items/{id}", handleItemRequest).Methods("GET")
	router.HandleFunc("/health", handleHealthRequest).Methods("GET")

	log.Fatal(h.ListenAndServe(":8080", router))
}
