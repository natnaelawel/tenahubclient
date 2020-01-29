package handler

import (
	"html/template"
	"net/http"
	"github.com/tenahubclientdocker/entity"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"
)

type ServiceHandler struct {
	temp *template.Template
}
func NewServiceHandler(T *template.Template) *ServiceHandler {
	return &ServiceHandler{temp: T}
}


func (adh *ServiceHandler) AddService(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("healthcenter_id"))
	name := r.FormValue("name")
	description := r.FormValue("description")

	// healthcenter id is get from the cookie
	data := entity.Service{Name:name, Description:description,HealthCenterID: uint(id)}
	jsonValue, _ := json.Marshal(data)
	response, err := http.Post("http://localhost:8181/v1/service","application/json",bytes.NewBuffer(jsonValue))
	var status addStatus
	if err != nil {
		status.Success = false
		fmt.Println(err)
	}else {
		status.Success = true
		fmt.Println(response.StatusCode)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}


func (adh *ServiceHandler) EditService(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut{
		name := r.FormValue("name")
		description := r.FormValue("description")
		id,_ := strconv.Atoi(r.FormValue("hidden_service_id"))
		// healthcenter

		data := entity.Service{ID :uint(id),Name:name, Description:description}
		jsonValue, _ := json.Marshal(data)
		URL := fmt.Sprintf("http://localhost:8181/v1/service/%d", id)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(jsonValue))
		_, err = client.Do(req)
		var status addStatus
		if err != nil {
			status.Success = false
			fmt.Println(err)
		}else {
			status.Success = true
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	}

}

func (adh *ServiceHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	id,_ := strconv.Atoi(r.FormValue("hidden_service_id"))
	URL := fmt.Sprintf("http://localhost:8181/v1/service/%d",id)
	fmt.Println("We are here")
	req, err := http.NewRequest(http.MethodDelete,URL,nil)
	var status addStatus

	if err != nil {
		status.Success = false
		fmt.Println(err)
	}else {
		status.Success = true
		fmt.Println(req.URL.Path)
	}
	res, err := client.Do(req)
	if err != nil {
		status.Success = false
		fmt.Println(err)
	}else {
		status.Success = true
		fmt.Println(res.StatusCode)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

