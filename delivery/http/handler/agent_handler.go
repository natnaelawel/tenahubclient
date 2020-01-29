package handler

import (
	"net/http"
	"html/template"
	"fmt"
	"encoding/json"
	"bytes"
	"strconv"
	"github.com/tenahubclientdocker/service"
	"github.com/tenahubclientdocker/entity"
	"github.com/tenahubclientdocker/form"
	"net/url"
	"github.com/tenahubclientdocker/rtoken"
)


type AgentHandler struct {
	temp *template.Template
	CsrfSignKey []byte
	loggedInAgent *entity.User
}
func NewAgentHandler(T *template.Template, csk []byte) *AgentHandler {
	return &AgentHandler{temp: T, CsrfSignKey: csk}
}

func (adh *AgentHandler) AddAgent(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	//username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phonenum")
	password := []byte(r.FormValue("password"))

	hashedPassword,err := HashPassword(password)
	//data := entity.Agent{FirstName:firstName, LastName:lastName, UserName:username, Email:email,PhoneNumber:phone,Password:hashedPassword}
	data := entity.User{FirstName:firstName, LastName:lastName, Email:email, PhoneNumber:phone, Password:hashedPassword, Role:"agent"}
	jsonValue, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/%s", service.BaseURL, "users")
	_, err = http.Post(url,"application/json",bytes.NewBuffer(jsonValue))
	var status addStatus
	if err != nil {
		status.Success = false
	}else {
		status.Success = true
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}


func (adh *AgentHandler) EditAgent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Editing agent")
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	//username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phonenumber")
	password := r.FormValue("password")
	id,_ := strconv.Atoi(r.FormValue("hidden_id"))
	f, _, _ := r.FormFile("upload_image")
	if f != nil{
		_, err := FileUpload(r,"agent_uploads")
		if err != nil{
			fmt.Println(err)
		}
	}


	data := entity.User{ID:uint(id), FirstName:firstName, LastName:lastName, Email:email,PhoneNumber:phone,Password:password}
	jsonValue, _ := json.Marshal(data)
	URL := fmt.Sprintf("%s/%s/%d",service.BaseURL,"users", id)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(jsonValue))
	resp, err := client.Do(req)
	fmt.Printf("status code : %d\n", resp.StatusCode)
	var status addStatus
	if err != nil {
		panic(err)
		status.Success = false
	}else {
		status.Success = true
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

func (adh *AgentHandler) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	id,_ := strconv.Atoi(r.FormValue("hidden_id"))
	URL := fmt.Sprintf("%s/%s/%d",service.BaseURL,"users",id)

	req, err := http.NewRequest(http.MethodDelete,URL,nil)
	var status addStatus

	if err != nil {
		status.Success = false
	}else {
		status.Success = true
	}
	_, err = client.Do(req)

	if err != nil {
		status.Success = false
	}else {
		status.Success = true
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}
type addStatus struct {
	Success bool
}
type agentDatas struct {
	Agent entity.User
	HealthCenters []entity.HealthCenter
	PendingServices []entity.Service
	Form form.Input
}
func (ah *AgentHandler) AgentPage(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(ah.CsrfSignKey)
	agentForm := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		CSRF:    token,
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	//fmt.Println(r.URL.RawQuery)
	//fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	agentData, err := service.FetchAgent(id)
	if err != nil {
		http.Redirect(w, r, "http://localhost:8282/login", http.StatusSeeOther)
		return
	}
	ah.loggedInAgent = agentData
	healthcentersByAgent, err := service.FetchHealthCenterByAgentId(uint(id))
	pendingServices, err := service.FetchPendingServices(uint(id))

	fmt.Println("agent are ", agentData)
	fmt.Println("healtcenters are ", healthcentersByAgent)
	fmt.Println("pending services are ", pendingServices)
	data := agentDatas{Agent: *agentData, HealthCenters:healthcentersByAgent, PendingServices:pendingServices, Form:agentForm}
	//fmt.Println("the data is ", data)
	//if err != nil {
	//	w.WriteHeader(http.StatusNoContent)
	//}
	ah.temp.ExecuteTemplate(w, "agent_home.layout",data)
	//adh.temp.ExecuteTemplate(w, "agent_home.layout", data{admin,agents, healthCenters, users})
}

func (ah *AgentHandler) AddHealthCenter(w http.ResponseWriter, r *http.Request) {
	//c, err := r.Cookie("agent")
	//id, _ := strconv.Atoi(c.Value)
	id := ah.loggedInAgent.ID
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phonenum")
	city := r.FormValue("city")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm{
		fmt.Println("password is not same")
		return
	}

	data := entity.HealthCenter{Name:name,Email:email,PhoneNumber:phone,City:city,Password:password, AgentID:uint(id)}
	fmt.Println("the data is ", data)
	jsonValue, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/%s", service.BaseURL, "healthcenter/addhealthcenter")
	res, err := http.Post(url,"application/json",bytes.NewBuffer(jsonValue))
	var status addStatus
	fmt.Println(res.StatusCode)
	if err != nil {
		status.Success = false
	}else {
		status.Success = true
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}
