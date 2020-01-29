package service

import (
	"time"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"fmt"
	"bytes"
	"strings"
	"net/url"
	"github.com/tenahubclientdocker/entity"
)

type cookie struct {
	Key        string
	Expiration time.Time
}

type response struct {
	Status  string
	Content interface{}
}

var loggedIn = make([]cookie, 10)

//const baseURL string = "https://tenahubapi.herokuapp.com/v1"
const baseURL string = "http://localhost:8181/v1"
var BaseURL = baseURL
func getResponse(request *http.Request) []byte {
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	if 200 != resp.StatusCode {
		panic(errors.New("status not correct"))
	}

	return body
}

// PostUser posts user to api
func PostUser(user *entity.User) error {
	requestBody, err := json.MarshalIndent(user, "", "\n")
	URL := fmt.Sprintf("%s/%s", baseURL, "users")

	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("reading: %s", err)
	if err != nil {
		fmt.Println(err)
		return err
	}

	cmp := strings.Compare(string(body), "Not Found")
	fmt.Println(string(body))
	fmt.Printf("comparison: %d",cmp)
	if  string(body) == "Not Found" {
		return errors.New("duplicate")
	}
	fmt.Println(string(body))
	return nil
}

// Authenticate authenticates user
func Authenticate(user *entity.User) (*entity.User, error) {
	URL := fmt.Sprintf("%s/%s", baseURL, "user")

	formval := url.Values{}
	formval.Add("email", user.Email)
	formval.Add("password", user.Password)

	resp, err := http.PostForm(URL, formval)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	respjson := struct {
		Status  string
		Content entity.User
	}{}

	err = json.Unmarshal(body, &respjson)

	fmt.Println(respjson)

	if respjson.Status == "error" {
		return nil, errors.New("error")
	}

	return &respjson.Content, nil
}

// GetHealthcenters gets healthcenters
func GetHealthcenters(name string, column string) ([]entity.Hcrating, error) {
	URL := fmt.Sprintf("%s/%s?search-key=%s&column=%s", baseURL, "healthcenters/search", name, column)

	fmt.Println(URL)
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	hcs := []entity.Hcrating{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &hcs)

	if err != nil {
		return nil, err
	}

	fmt.Println(hcs)

	return hcs, nil
}

// GetHealthcenter gets health center by id
func GetHealthcenter(id uint) (*entity.HealthCenter, error) {
	URL := fmt.Sprintf("%s/%s/%d", baseURL, "healthcenter", id)

	fmt.Println(URL)
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	hcs := entity.HealthCenter{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &hcs)

	if err != nil {
		return nil, err
	}

	fmt.Println(hcs)

	return &hcs, nil
}

// GetServices gets healthcenters services
func GetServices(id uint) ([]entity.Service, error) {
	URL := fmt.Sprintf("%s/%s/%d", baseURL, "services", id)

	fmt.Println(URL)
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	services := []entity.Service{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &services)

	if err != nil {
		return nil, err
	}

	fmt.Println(services)

	return services, nil
}

// GetRating gets healthcenters rating
func GetRating(id uint) (float64, error) {
	URL := fmt.Sprintf("%s/%s/%d", baseURL, "rating", id)

	fmt.Println(URL)
	resp, err := http.Get(URL)

	if err != nil {
		return 0.0, err
	}

	rating := struct {
		Rating float64
	}{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0.0, err
	}

	err = json.Unmarshal(body, &rating)

	if err != nil {
		return 0.0, err
	}

	fmt.Println(rating)

	return rating.Rating, nil
}

// PostFeedback posts feedback
func PostFeedback(comment *entity.Comment) error {
	URL := fmt.Sprintf("%s/%s", baseURL, "rating")

	data, err := json.MarshalIndent(comment, "", "\t")

	if err != nil {
		return err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}

// CheckValidity checks if user is valid to give feedback
func CheckValidity(uid uint, hid uint) (string, error) {
	URL := fmt.Sprintf("%s/comments/check", baseURL)

	comment := entity.Comment{HealthCenterID: hid, UserID: uid}

	output, err := json.MarshalIndent(&comment, "", "\t")

	if err != nil {
		return "", err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(output))

	if err != nil {
		return "", err
	}

	l := resp.ContentLength
	data := make([]byte, l)
	resp.Body.Read(data)

	result := struct {
		Status string
	}{}

	err = json.Unmarshal(data, &result)

	if err != nil {
		return "", err
	}

	if strings.Compare(result.Status, "valid") == 0 {
		return result.Status, nil
	}

	return "", nil

}

// GetTop returns top rated healthcenters
func GetTop(amount uint)([]entity.Hcrating, error) {
	URL := fmt.Sprintf("%s/%s/%s/%d", baseURL, "healthcenters", "top", amount)
	client := http.Client{}

	request, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	l := resp.ContentLength
	data := make([]byte, l)
	resp.Body.Read(data)

	result := []entity.Hcrating{}

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return result, nil
}

// GetFeedback gets feedback
func GetFeedback(id uint)([]entity.UserComment, error) {
	URL := fmt.Sprintf("%s/%s/%d",baseURL, "comments", id)
	fmt.Println(URL)
	comments := []entity.UserComment{}

	client := http.Client{}

	request, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)

	l := resp.ContentLength
	data := make([]byte, l)
	resp.Body.Read(data)

	err = json.Unmarshal(data, &comments)

	if err != nil {
		return nil, err
	}
	fmt.Println(comments)
	return comments, nil

}

func FetchUsers() ([]entity.User, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%s/users/user/type", baseURL)
	req, _ := http.NewRequest("GET", URL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var users []entity.User
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

