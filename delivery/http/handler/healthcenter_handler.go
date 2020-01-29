package handler

import (
	"net/http"
	"html/template"
	"fmt"
	"strconv"
	// "github.com/TenaHub/api/entity"
	//"encoding/json"
	//"bytes"
	"github.com/tenahubclientdocker/service"
	"os"
	"io"
	"path/filepath"
	"encoding/json"
	"bytes"
	"github.com/tenahubclientdocker/entity"
	"github.com/tenahubclientdocker/rtoken"
	"net/url"
	"github.com/tenahubclientdocker/form"
	"github.com/tenahubclientdocker/session"
	"strings"
	"context"
	"github.com/tenahubclientdocker/permission"
)


type HealthCenterHandler struct {
	temp *template.Template
	UserSess     *entity.Session
	LoggedInUser *entity.HealthCenter
	CsrfSignKey  []byte
}
func NewHealthCenterHandler(T *template.Template, usrSess *entity.Session, csKey []byte) *HealthCenterHandler {
	return &HealthCenterHandler{temp: T, UserSess:usrSess, CsrfSignKey:csKey}
}
type healthcenterData struct {
	HealthCenter *entity.HealthCenter
	FeedBack []entity.Comment
	Service []entity.Service
	Form form.Input
}

func (uh *HealthCenterHandler) LoggedIn(r *http.Request) bool {
	if uh.UserSess == nil {
		return false
	}
	UserSess := uh.UserSess
	fmt.Printf("usersess: %q\n", UserSess)
	c, err := r.Cookie(UserSess.UUID)
	//fmt.Println(c)
	if err != nil {
		return false
	}

	fmt.Printf("logged In: %s" ,UserSess.SigningKey)
	ok, err := session.Valid(c.Value, UserSess.SigningKey)
	fmt.Println(err)
	if !ok || (err != nil) {
		return false
	}
	return true
}

func (hch *HealthCenterHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := hch.LoggedIn(r)
		fmt.Println(ok)
		if !ok {
			http.Redirect(w, r, "/healthcenter/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, hch.UserSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Authorized checks if a user has proper authority to access a give route
func (hch *HealthCenterHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := ""
		if !(hch.LoggedInUser == nil) {
			role = "HEALTH_CENTER"
		}

		permitted := permission.HasPermission(r.URL.Path, strings.ToUpper(role) , r.Method)
		fmt.Printf("permitted: %t\n", permitted)
		if !permitted {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		fmt.Println("here 1")
		if r.Method == http.MethodPost {
			fmt.Println("here")
			fmt.Println(r.PostFormValue("_csrf"))
			ok, err := rtoken.ValidCSRF(r.FormValue("_csrf"), hch.CsrfSignKey)
			if !ok || (err != nil) {
				fmt.Println("also here")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (adh *HealthCenterHandler) EditHealthCenter(w http.ResponseWriter, r *http.Request) {
	//c, err := r.Cookie("healthcenter")
	//id, _ := strconv.Atoi(c.Value)
	fmt.Println("editing healthcenter")
	token, err := rtoken.CSRFToken(adh.CsrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	formData := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
	formData.Required("name","email", "address", "phone", "password", "confirmpassword")
	formData.MatchesPattern("email", form.EmailRX)
	//formData.MatchesPattern("phone", form.PhoneRX)
	formData.MinLength("password", 8)
	formData.PasswordMatches("password", "confirmpassword")
	formData.CSRF = token

	if !formData.Valid() {
		fmt.Println("From Is Valid")
		fmt.Println(formData.VErrors)
		id := adh.LoggedInUser.ID
		healthcenter, err := service.FetchHealthCenter(uint(id))
		if err != nil {
			http.Redirect(w, r, "/healthcenter/login", http.StatusSeeOther)
		}
		services, err := service.FetchService(uint(id))
		fmt.Println("service is ", services)
		feedbacks, err := service.FetchFeedbacks(uint(id))

		data := healthcenterData{HealthCenter:healthcenter, FeedBack:feedbacks, Service:services, Form:formData}
		fmt.Println("data is ", data)
		adh.temp.ExecuteTemplate(w, "healthcenter_home.layout", data)
		//adh.temp.ExecuteTemplate(w, "healthcenter_edit_profile.layout", struct {Form form.Input}{formData})
		return
	}
	fmt.Println("From Is Valid")

	id := adh.LoggedInUser.ID

	Name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	city := r.FormValue("address")
	password := r.FormValue("password")
	confirm := r.FormValue("confirmpassword")

	if password != confirm {
		fmt.Println("Passwords Doesn't match")
		return
	}
	var fileName string
	f, _, _ := r.FormFile("upload_image")
	if f != nil {
		fileName, err = FileUpload(r,"healthcenter_uploads")
		if err != nil{
			fmt.Println(err)
		}
	}
	fileName = ""
	data := entity.HealthCenter{ID:uint(id),Name:Name, Email:email,PhoneNumber:phone,City:city,Password:password,ProfilePic:fileName}
	jsonValue, _ := json.Marshal(data)
	URL := fmt.Sprintf("http://localhost:8181/v1/healthcenter/%d", id)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(jsonValue))
	resp, err := client.Do(req)

	fmt.Println(resp)
	var status addStatus
	if err != nil {
		status.Success = false
		fmt.Println(err)
	}else {
		status.Success = true
	}
	fmt.Println(err)

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	//adh.temp.ExecuteTemplate(w, "admin_home.layout", status)
}


func (adh *HealthCenterHandler) DeleteHealthCenter(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	id,_ := strconv.Atoi(r.FormValue("hidden_id"))
	URL := fmt.Sprintf("http://localhost:8181/v1/healthcenter/%d",id)

	req, err := http.NewRequest(http.MethodDelete,URL,nil)

	res, err := client.Do(req)
	var status addStatus
	if err != nil {
		status.Success = false
		fmt.Println(err)
	}else {
		status.Success = true
		fmt.Println(res.StatusCode)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}


func (ah *HealthCenterHandler) HealthCenterPage(w http.ResponseWriter, r *http.Request) {
	token, err := rtoken.CSRFToken(ah.CsrfSignKey)
	formData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		CSRF:    token,
	}

	id := ah.LoggedInUser.ID
	healthcenter, err := service.FetchHealthCenter(uint(id))
	if err != nil {
		http.Redirect(w, r, "/healthcenter/login", http.StatusSeeOther)
	}
	services, err := service.FetchService(uint(id))
	fmt.Println("service is ", services)
	feedbacks, err := service.FetchFeedbacks(uint(id))

	data := healthcenterData{HealthCenter:healthcenter, FeedBack:feedbacks, Service:services, Form:formData}
	fmt.Println("data is ", data)

	ah.temp.ExecuteTemplate(w, "healthcenter_home.layout", data)
}

// Login handles Get /login and POST /login
func (ah *HealthCenterHandler) HealthCenterLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Referer())
	token, err := rtoken.CSRFToken(ah.CsrfSignKey)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		loginForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		ah.temp.ExecuteTemplate(w, "healthcenter.login.layout", loginForm)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		loginForm.Required("password")
		healthcenter := entity.HealthCenter{Email: email,Password : password}
		resp, err := service.HealthCenterAuthenticate(&healthcenter)
		//fmt.Printf("responsee: %q\n", resp)
		if err != nil {
			if err.Error() == "error" {
				loginForm.VErrors.Add("generic", "email or password not correct")
				ah.temp.ExecuteTemplate(w, "healthcenter.login.layout", loginForm)
				return
			}
		} else{
			ah.LoggedInUser = resp
			claims := rtoken.Claims(resp.Email, ah.UserSess.Expires)
			session.Create(claims, ah.UserSess.UUID, ah.UserSess.SigningKey, w)
			fmt.Printf("sess: %s", ah.UserSess.SigningKey)
			newSess, err := service.StoreSession(ah.UserSess)
			fmt.Println("here", ah.UserSess.SigningKey)
			if err!=nil{
				loginForm.VErrors.Add("generic", "Failed to store session")
				ah.temp.ExecuteTemplate(w, "healthcenter.login.layout", loginForm)
				return
			}
			ah.UserSess = newSess
			http.Redirect(w, r, "http://localhost:8282/healthcenter", http.StatusSeeOther)
		}
	}
}
// Logout handles GET /logout
func (uh *HealthCenterHandler) HealthCenterLogout(w http.ResponseWriter, r *http.Request) {

	session.Remove(uh.UserSess.UUID, w)
	service.DeleteSession(uh.UserSess.UUID)
	uh.LoggedInUser = nil
	fmt.Println("logging out")
	http.Redirect(w, r, "http://localhost:8282/healthcenter/login", http.StatusSeeOther)
}

func FileUpload(r *http.Request, folderName string) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("upload_image")
	if err != nil {
		panic(err)
		return "",err
	}
	defer file.Close()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "client","ui", "assets", "img", "uploads",folderName, header.Filename)

	f, err := os.OpenFile(path,os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		return "",err
	}
	defer f.Close()
	io.Copy(f, file)
	return header.Filename, nil
}