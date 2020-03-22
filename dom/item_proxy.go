package dom

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ItemProxy struct {
	LocalService  ItemService
	RemoteService ItemService
	CacheService  CacheItemService
}

func NewItemProxy(local ItemService, remote ItemService, cache CacheItemService) *ItemProxy {
	return &ItemProxy{
		LocalService:  local,
		RemoteService: remote,
		CacheService:  cache,
	}
}

func (proxy ItemProxy) RequestHandler(w http.ResponseWriter, r *http.Request) {
	var item *Item
	var err error

	itemID := mux.Vars(r)["id"]

	log.Println("GET item/" + itemID)

	item, err = proxy.LocalService.Item(itemID)

	if err != nil {
		log.Println("GET api.ml/item")

		item, err = proxy.RemoteService.Item(itemID)

		if err != nil {
			log.Fatal(err)
		}
		log.Println("WRITE api.ml/item To Local Cache")

		proxy.CacheService.Write(item)
	} else {
		log.Println("READ Local DB")
	}
	response, err := json.Marshal(item)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
