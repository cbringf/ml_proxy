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
	ID                   int64     `db:"id"`
	ItemID               string    `db:"item_id"`
	Remote               bool      `db:"remote"`
	RemoteResponseTime   int       `db:"remote_response_time"`
	RemoteResponseStatus int       `db:"remote_response_status"`
	ResponseStatus       int       `db:"response_status"`
	ResponseTime         int       `db:"response_time"`
	RequestDate          time.Time `db:"request_date"`
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
			err := proxy.CacheService.Write(item)

			if err != nil {
				log.Printf("ERROR WRITING api.ml/item To Local Cache")
			}
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
	} else {
		_, err = stmt.Exec(&reqInfo.ItemID, &reqInfo.Remote, &reqInfo.ResponseStatus, &reqInfo.ResponseTime, &reqInfo.RemoteResponseStatus, &reqInfo.RemoteResponseTime, &reqInfo.RequestDate)

		if err != nil {
			log.Println("FAILED WRITE Request Info")
		}
	}
}

func (proxy ItemProxy) ReadRequests() ([]RequestInfo, *Error) {
	var result = make([]RequestInfo, 0)

	rows, _ := proxy.DB.Query("SELECT id, item_id, (remote = b'1'), response_status, response_time, request_date, remote_response_time, remote_response_status FROM request")

	for rows.Next() {
		var aux RequestInfo

		_ = rows.Scan(&aux.ID, &aux.ItemID, &aux.Remote, &aux.ResponseStatus, &aux.ResponseTime, &aux.RequestDate, &aux.RemoteResponseTime, &aux.RemoteResponseStatus)

		result = append(result, aux)
	}
	return result, nil
}
