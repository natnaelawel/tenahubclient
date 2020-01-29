package service

import (
	"github.com/tenahubclientdocker/entity"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
)



func FetchSession(uuid string)(*entity.Session, error){
	client := http.Client{}
	URL := fmt.Sprintf("%s%s?uuid=%s", baseURL, "/session", uuid)
	fmt.Println(URL)
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("here1")
		return nil, err
	}
	var session entity.Session
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("here2")
		return nil, err
	}

	err = json.Unmarshal(body, &session)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println("here3")
		return nil, err
	}
	return &session, nil
}

func StoreSession(session *entity.Session)(*entity.Session, error){
	client := http.Client{}
	URL := fmt.Sprintf("%s%s", baseURL, "/session")
	fmt.Println(URL)
	data, err := json.MarshalIndent(session, "", "\t\t")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("service: %s\n", session.SigningKey)
	fmt.Printf("session uuid: %s\n", session.UUID)
	sess, err := FetchSession(session.UUID)

	if err != nil {
		fmt.Printf("error: %s\n",err)
		return nil, err
	}

	return sess, nil
}

func DeleteSession(uuid string)error{
	client := http.Client{}
	URL := fmt.Sprintf("%s%s%s", baseURL, "/session/", uuid)
	req, _ := http.NewRequest(http.MethodDelete, URL, nil)

	_, err := client.Do(req)

	if err!=nil{
		return err
	}

	return nil
}