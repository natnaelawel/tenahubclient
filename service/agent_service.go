package service

import (
	"encoding/json"
	"fmt"
	"github.com/tenahubclientdocker/entity"
	"io/ioutil"
	"net/http"
)

func FetchAgent(id int) (*entity.User, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/users/%d", baseURL, id)
	fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	userdata := entity.User{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(body)
	err = json.Unmarshal(body, &userdata)
	if err != nil {
		return nil, err
	}
	fmt.Println(userdata)
	return &userdata, nil
}

func FetchAgents() ([]entity.User, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/users/agent/type", baseURL)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var agents []entity.User
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &agents)
	if err != nil {
		return nil, err
	}
	return agents, nil
}
