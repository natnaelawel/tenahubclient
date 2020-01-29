package handler

import (
	"net/http"
	"html/template"
	"github.com/tenahubclientdocker/service"
	"fmt"
	"encoding/json"
	"bytes"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"github.com/tenahubclientdocker/entity"
	"net/url"
	"github.com/tenahubclientdocker/rtoken"
	"github.com/tenahubclientdocker/form"
)


type AdminHandler struct {
	temp *template.Template
	userHandl *UserHandler
	CsrfSignKey  []byte
}
func NewAdminHandler(T *template.Template, uh *UserHandler, csk []byte) *AdminHandler {
	return &AdminHandler{temp: T, userHandl:uh, CsrfSignKey:csk}
}


func (adh *AdminHandler) AllAgents(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	users, err := service.FetchAgent(6)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		adh.temp.ExecuteTemplate(w, "admin_home.layout", nil)
	}
	adh.temp.ExecuteTemplate(w, "check.html", users)
}
type data struct {
	Admin *entity.User
	Agent []entity.User
	HealthCenter []entity.HealthCenter
	User []entity.User
	Form form.Input
}
func (adh *AdminHandler) AdminPage(w http.ResponseWriter, r *http.Request) {

		token, err := rtoken.CSRFToken(adh.CsrfSignKey)
		agentForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		//id, _ := strconv.Atoi(r.Form.Get("id"))
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		//fmt.Println(r.URL.RawQuery)
		admin, err := service.FetchAdmin(id)
		//fmt.Println(admin)
		agents, err := service.FetchAgents()
		healthCenters, err := service.FetchHealthCenters()
		users, err := service.FetchUsers()

		if err != nil {
			fmt.Println("here")
			fmt.Println(err)
			w.WriteHeader(http.StatusNoContent)
			return
			//http.Redirect(w, r, "http://localhost:8282/admin/login", http.StatusSeeOther)
		}
		//fmt.Println(admin)
		adh.temp.ExecuteTemplate(w, "admin_home.layout", data{admin,agents, healthCenters, users, agentForm})
		return
}

func (adh *AdminHandler) EditAdmin(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("admin")
	id, _ := strconv.Atoi(c.Value)

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	//username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm {
		return
	}
	fileName, err := FileUpload(r,"admin_uploads")
	if err != nil{
		fmt.Println(err)
	}

	data := entity.User{FirstName:firstName, LastName:lastName, Email:email,PhoneNumber:phone,Password:password,ProfilePic:fileName}
	jsonValue, _ := json.Marshal(data)
	client := &http.Client{}

	URL := fmt.Sprintf("%s/%s/%d",service.BaseURL,"admin", id)

	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(jsonValue))
	_, err = client.Do(req)
	var status addStatus
	if err != nil {
		status.Success = false
	}else {
		status.Success = true
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func HashPassword(password []byte)(string, error){
	hashedPassword,err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(hashedPassword), err
}
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
