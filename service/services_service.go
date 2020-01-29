package service

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/tenahubclientdocker/entity"
)

func FetchService(id uint) ([]entity.Service, error) {
	client := &http.Client{}
	fmt.Println(id)
	URL := fmt.Sprintf("%s/services/%d", baseURL, id)
	fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var service []entity.Service
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &service)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return service, nil
}

func FetchPendingServices(id uint) ([]entity.Service, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/pending/services/%d", baseURL, id)
	fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var services []entity.Service
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
