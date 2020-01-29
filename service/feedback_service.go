package service

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/tenahubclientdocker/entity"
)

func FetchFeedbacks(id uint) ([]entity.Comment, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/feedback/%d", baseURL, id)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var feedbacks []entity.Comment
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &feedbacks)
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}


