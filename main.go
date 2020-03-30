package main

import (
	"fmt"
	"log"
	"os"

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
	snapshot := dom.BuildSnapshot(requestList)

	snapshot.HandleRequest(w, r)
}

func main() {
	dbConnStr := os.Getenv("DB_CONN_STR")

	if dbConnStr == "" {
		fmt.Println("LOADING DEV ENVIRONMENT")

		dbConnStr = "root:testDb@tcp(127.0.0.1:3306)/ml_proxy"
	}

	db, _ = s.Open("mysql", fmt.Sprintf("%s?parseTime=true", dbConnStr))
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items/{id}", handleItemRequest).Methods("GET")
	router.HandleFunc("/health", handleHealthRequest).Methods("GET")

	fmt.Println("LISTENING ON PORT 8080")
	log.Fatal(h.ListenAndServe(":8080", router))
}
