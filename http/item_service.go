package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cbringf/proxy/dom"
)

type ItemService struct {
	HTTP http.Client
}

func (s ItemService) Item(id string) (*dom.Item, error) {
	var mlResponse dom.Item
	// var mlChildrenResponse []dom.ItemChild

	itemResponse, err := s.HTTP.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s", id))

	if err != nil {
		return nil, err
	}

	childrenResponse, err := s.HTTP.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s/children", id))

	if err != nil {
		return nil, err
	}

	item, _ := ioutil.ReadAll(itemResponse.Body)
	children, _ := ioutil.ReadAll(childrenResponse.Body)

	json.Unmarshal(item, &mlResponse)
	json.Unmarshal(children, &mlResponse.Children)

	return &mlResponse, nil
}
