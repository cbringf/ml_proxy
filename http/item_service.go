package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cbringf/ml_proxy/main"
)

type ItemService struct {
	HTTP http.Client
}

func (s *main.ItemService) Item(id string) *main.Item {
	var mlResponse main.Item
	var mlChildrenResponse []main.ItemChild

	response1, err1 := http.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s", id))
	response2, err2 := http.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s/children", id))

	if err1 != nil {
		fmt.Printf("Http request fails with error %s", err1)
	} else if err2 != nil {
		fmt.Printf("Http request fails with error %s", err2)
	} else {
		data1, _ := ioutil.ReadAll(response1.Body)
		data2, _ := ioutil.ReadAll(response2.Body)

		json.Unmarshal(data2, &mlChildrenResponse)
		json.Unmarshal(data1, &mlResponse)
	}
}
