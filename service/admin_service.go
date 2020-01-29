package service

import (
	"encoding/json"
	"fmt"
	"github.com/tenahubclientdocker/entity"
	"io/ioutil"
	"net/http"
)
// var baseURL = "http://localhost:8181/v1"

func FetchAdmin(id int) (*entity.User, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/users/%d", baseURL, id)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	adminData := entity.User{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &adminData)
	fmt.Println("error is ",err)
	if err != nil {
		return nil, err
	}
	return &adminData, nil
}
