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

var config *dom.Config
var dbConnStr string

func openDbConn() (*sql.DB, error) {
	db, err := s.Open("mysql", fmt.Sprintf("%s?parseTime=true", dbConnStr))

	if err != nil {
		fmt.Println("UNABLE to load DB connection")
	}

	return db, err
}

func handleItemRequest(w h.ResponseWriter, r *h.Request) {
	db, _ := openDbConn()

	defer db.Close()

	dbItemService := mysql.ItemService{DB: db}
	proxy := dom.NewItemProxy(
		dbItemService,
		http.ItemService{
			HTTP:   *h.DefaultClient,
			Config: config,
		},
		dbItemService,
		db,
	)
	reqInfo := proxy.RequestHandler(w, r)
	proxy.LogRequest(&reqInfo)
}

func handleHealthRequest(w h.ResponseWriter, r *h.Request) {
	var snapshot dom.SysRequestSnapshot

	db, _ := openDbConn()

	defer db.Close()

	proxy := dom.ItemProxy{
		DB: db,
	}
	requestList, err := proxy.ReadRequests()

	if err != nil {
		snapshot = dom.SysRequestSnapshot{
			SnapshotError: err,
		}
	} else {
		snapshot = dom.BuildSnapshot(requestList)
	}

	snapshot.HandleRequest(w, r)
}

func main() {
	var err error

	dbConnStr = os.Getenv("DB_CONN_STR")
	config, err = dom.Load()

	if err != nil {
		return
	}

	if dbConnStr == "" {
		fmt.Println("LOADING DEV ENVIRONMENT")

		dbConnStr = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.DB.User,
			config.DB.Pass,
			config.DB.Host,
			config.DB.Port,
			config.DB.Name,
		)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items/{id}", handleItemRequest).Methods("GET")
	router.HandleFunc("/health", handleHealthRequest).Methods("GET")

	fmt.Println("SERVER RUNNING")
	log.Fatal(h.ListenAndServe(":8080", router))
}
