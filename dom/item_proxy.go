package dom

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ItemProxy struct {
	DB            *sql.DB
	LocalService  ItemService
	RemoteService ItemService
	CacheService  CacheItemService
}

type RequestInfo struct {
	ID                   int64
	ItemID               string
	Remote               bool
	RemoteResponseTime   int
	RemoteResponseStatus int
	ResponseStatus       int
	ResponseTime         int
	RequestDate          time.Time
}

func NewItemProxy(local ItemService, remote ItemService, cache CacheItemService, db *sql.DB) *ItemProxy {
	return &ItemProxy{
		LocalService:  local,
		RemoteService: remote,
		CacheService:  cache,
		DB:            db,
	}
}

func (proxy ItemProxy) RequestHandler(w http.ResponseWriter, r *http.Request) RequestInfo {
	var item *Item
	var reqInfo = RequestInfo{
		ResponseStatus: 200,
	}
	var response []byte

	handleStart := time.Now()
	itemID := mux.Vars(r)["id"]
	reqInfo.ItemID = itemID

	log.Println("GET item/" + itemID)

	item, localErr := proxy.LocalService.Item(itemID)

	if localErr != nil {
		remoteStart := time.Now()
		reqInfo.ResponseStatus = localErr.Status
		reqInfo.Remote = true

		log.Println("GET api.ml/item")

		reqInfo.RemoteResponseStatus = 200
		item, err := proxy.RemoteService.Item(itemID)
		reqInfo.RemoteResponseTime = int(time.Since(remoteStart).Milliseconds())

		if err != nil {
			reqInfo.RemoteResponseStatus = err.Status
			response, _ = json.Marshal(err)

			log.Println("ERROR Remote Api Call")
		} else {
			log.Println("WRITE api.ml/item To Local Cache")

			response, _ = json.Marshal(item)
			proxy.CacheService.Write(item)
		}
	} else {
		response, _ = json.Marshal(item)
		log.Println("READ Local DB")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	reqInfo.ResponseTime = int(time.Since(handleStart).Milliseconds())
	reqInfo.RequestDate = handleStart

	return reqInfo
}

func (proxy ItemProxy) LogRequest(reqInfo *RequestInfo) {
	var query = `
		INSERT INTO request (item_id, remote, response_status, response_time, remote_response_status, remote_response_time, request_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := proxy.DB.Prepare(query)

	if err != nil {
		log.Println("FAILED WRITE Request Info")
	}

	_, err = stmt.Exec(&reqInfo.ItemID, &reqInfo.Remote, &reqInfo.ResponseStatus, &reqInfo.ResponseTime, &reqInfo.RemoteResponseStatus, &reqInfo.RemoteResponseTime, &reqInfo.RequestDate)

	if err != nil {
		log.Println("FAILED WRITE Request Info")
	}
}
