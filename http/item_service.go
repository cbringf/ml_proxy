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

func (s ItemService) Item(id string) (*dom.Item, *dom.Error) {
	var mlResponse dom.Item

	itemRes, err := request(fmt.Sprintf("https://api.mercadolibre.com/items/%s", id))

	if err != nil {
		return nil, err
	}

	childrenRes, err := request(fmt.Sprintf("https://api.mercadolibre.com/items/%s/children", id))

	if err != nil {
		return nil, err
	}

	uError := json.Unmarshal(itemRes, &mlResponse)

	if uError != nil {
		return nil, dom.UnknownError()
	}

	uError = json.Unmarshal(childrenRes, &mlResponse.Children)

	if uError != nil {
		return nil, dom.UnknownError()
	}

	return &mlResponse, nil
}

func request(url string) ([]byte, *dom.Error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, dom.UnknownError()
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, dom.UnknownError()
	} else if response.StatusCode != 200 {
		var apiError = dom.UnknownError()

		json.Unmarshal(responseBody, apiError)

		return nil, apiError
	} else {
		return responseBody, nil
	}
}
