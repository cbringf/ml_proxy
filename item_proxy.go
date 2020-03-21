package main

type ItemProxy struct {
	LocalService  ItemService
	RemoteService ItemService
}

func NewItemProxy(local ItemService, remote ItemService) *ItemProxy {
	return &ItemProxy{
		LocalService:  local,
		RemoteService: remote,
	}
}

func (proxy *ItemProxy) Item(id string) (*Item, error) {
	res, err := proxy.LocalService.Item(id)

	if err != nil {
		res, err := proxy.RemoteService.Item(id)

		return res, err
	}
	return res, err
}
