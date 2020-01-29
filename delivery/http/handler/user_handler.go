package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/tenahubclientdocker/service"
	"github.com/tenahubclientdocker/session"

	"github.com/tenahubclientdocker/entity"
	"github.com/tenahubclientdocker/rtoken"
	"github.com/tenahubclientdocker/form"
	"golang.org/x/crypto/bcrypt"
	"github.com/tenahubclientdocker/permission"
)

// UserHandler handles user related http requests
type UserHandler struct {
	Templ        *template.Template
	UserSess     *entity.Session
	LoggedInUser *entity.User
	CsrfSignKey  []byte
}

// NewUserHandler creates object of UserHandler
func NewUserHandler(tmpl *template.Template,	usrSess *entity.Session, csKey []byte) *UserHandler {
	return &UserHandler{Templ: tmpl, UserSess:usrSess, CsrfSignKey:csKey}
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")

func (uh *UserHandler) GetLoggedInUser()*entity.User{
	return uh.LoggedInUser
}

// Authenticated checks if a user is authenticated to access a given route
func (uh *UserHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := uh.LoggedIn(r)
		fmt.Println(ok)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, uh.UserSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Authorized checks if a user has proper authority to access a give route
func (uh *UserHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := ""
		if !(uh.LoggedInUser == nil) {
			role = uh.LoggedInUser.Role
		}
		//roles, errs := uh.userService.UserRoles(uh.LoggedInUser)
		//if len(errs) > 0 {
		//	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		//	return
		//}

		//for _, role := range roles {
		//	permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
		//	if !permitted {
		//		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		//		return
		//	}
		//}

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
			ok, err := rtoken.ValidCSRF(r.FormValue("_csrf"), uh.CsrfSignKey)
			if !ok || (err != nil) {
				fmt.Println("also here")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (uh *UserHandler) LoggedIn(r *http.Request) bool {
	if uh.UserSess == nil {
		return false
	}
	UserSess := uh.UserSess
	fmt.Printf("usersess: %q", UserSess)
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

// Index handles GET /
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	hcs, err := service.GetTop(4)

	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	fmt.Println(hcs)

	uh.Templ.ExecuteTemplate(w, "user.index.default.layout", hcs)
}

// Login handles Get /login and POST /login
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Referer())
	token, err := rtoken.CSRFToken(uh.CsrfSignKey)

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
		uh.Templ.ExecuteTemplate(w, "user.login.layout", loginForm)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		user := entity.User{Email: email, Password: password}
		fmt.Println(user)

		resp, err := service.Authenticate(&user)

		if err != nil {
			if err.Error() == "error" {
				loginForm.VErrors.Add("generic", "Your email address or password is wrong")
				uh.Templ.ExecuteTemplate(w, "user.login.layout", loginForm)
				return
			}
		} else {

			//fmt.Println(resp)
			uh.LoggedInUser = resp
			claims := rtoken.Claims(resp.Email, uh.UserSess.Expires)
			session.Create(claims, uh.UserSess.UUID, uh.UserSess.SigningKey, w)
			fmt.Printf("sess: %s", uh.UserSess.SigningKey)
			newSess, err := service.StoreSession(uh.UserSess)
			fmt.Println("here", uh.UserSess.SigningKey)
			if err!=nil{
				loginForm.VErrors.Add("generic", "Failed to store session")
				uh.Templ.ExecuteTemplate(w, "user.login.layout", loginForm)
				return
			}
			uh.UserSess = newSess

			fmt.Printf("referer: %s\n", r.Referer())
			var urll string
			fmt.Println(uh.LoggedInUser)
			fmt.Println("role: "+uh.LoggedInUser.Role)
			fmt.Printf("id: %d\n", uh.LoggedInUser.ID)
			switch uh.LoggedInUser.Role {
			case "admin":
				urll = "http://localhost:8282/admin"
			case "agent":
				urll = "http://localhost:8282/agent"
			case "user":
				urll = "http://localhost:8282/home"
			}
			fmt.Println(urll)
			//r.Form.Add("id", strconv.Itoa(int(uh.LoggedInUser.ID)))
			r.URL.Query().Add("id", strconv.Itoa(int(uh.LoggedInUser.ID)))
			//r.Header.Add("id", strconv.Itoa(int(uh.LoggedInUser.ID)))
			//r.Method = http.MethodPost
			Url, err := url.Parse(urll)
			parameters := url.Values{}
			parameters.Add("id", strconv.Itoa(int(uh.LoggedInUser.ID)))
			Url.RawQuery = parameters.Encode()

			if err != nil {
				http.Redirect(w, r, "http://localhost:8282/login", http.StatusSeeOther)
			}

			if r.Referer() ==  "http://localhost:8282/login" {
				http.Redirect(w, r, Url.String(), http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			// uh.Templ.ExecuteTemplate(w, "user.index.auth.layout", resp)
		}

	}
}

// Home handles GET /home
func (uh *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	hcs, err := service.GetTop(4)
	fmt.Printf("home: %s\n", uh.UserSess.SigningKey)
	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	//fmt.Println(uh.LoggedIn(r))

	//c, err := r.Cookie("user")
	if !uh.LoggedIn(r) {
		uh.Templ.ExecuteTemplate(w, "user.index.default.layout", hcs)
		return
	} else {
		//fmt.Println(c.Value)
		//fmt.Println(c.MaxAge)
	}
	uh.Templ.ExecuteTemplate(w, "user.index.auth.layout", hcs)
}

// SignUp handles GET /signup and POST /signup
func (uh *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	token, err := rtoken.CSRFToken(uh.CsrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		signUpForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		uh.Templ.ExecuteTemplate(w, "user.signup.layout", signUpForm)
		return
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Println(r.PostForm)
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.Required("firstname","lastname", "email", "password", "confirmpassword")
		signUpForm.MatchesPattern("email", form.EmailRX)
		signUpForm.MatchesPattern("phone", form.PhoneRX)
		signUpForm.MinLength("password", 8)
		signUpForm.PasswordMatches("password", "confirmpassword")
		signUpForm.CSRF = token

		if !signUpForm.Valid() {
			fmt.Println(signUpForm.VErrors)
			uh.Templ.ExecuteTemplate(w, "user.signup.layout", signUpForm)
			return
		}

		firstname := r.PostFormValue("firstname")
		lastname := r.PostFormValue("lastname")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		phonenum := r.PostFormValue("phonenum")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		user := entity.User{FirstName: firstname, LastName: lastname, Email: email, Password: string(hashedPassword), PhoneNumber: phonenum, Role: "user"}
		fmt.Println(user)

		err = service.PostUser(&user)
		fmt.Println(err)
		if err != nil {
			if strings.Compare(err.Error(), "duplicate") == 0 {
				fmt.Println("duplicate")
				signUpForm.VErrors.Add("generic", "email or phone number is already taken")
				uh.Templ.ExecuteTemplate(w, "user.signup.layout", signUpForm)
				return
			}
			w.Write([]byte("failed"))
			return
		} else {
			//w.Write([]byte("success"))
			//w.Header().Set("Location:", "http://locahost:8282/login")
			http.Redirect(w, r, "http://localhost:8282/login", http.StatusSeeOther)
		}
	}
}

// Search handles GET /search
func (uh *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	searchkey := r.URL.Query().Get("search-key")
	column := r.URL.Query().Get("column")
	if column == "" {
		column = "name"
	}
	fmt.Println(searchkey)

	healthcenters, err := service.GetHealthcenters(searchkey, column)

	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	length := len(healthcenters)

	data := struct {
		Length  int
		Content []entity.Hcrating
	}{
		Length:  length,
		Content: healthcenters,
	}

	//c, err := r.Cookie("user")

	if !uh.LoggedIn(r) {
		uh.Templ.ExecuteTemplate(w, "user.result.default.layout", data)
		return
	} else {
		//fmt.Println(c.Value)
		//fmt.Println(c.MaxAge)
	}

	uh.Templ.ExecuteTemplate(w, "user.result.auth.layout", data)
}

// Healthcenters handles GET /healthcenters
func (uh *UserHandler) Healthcenters(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	//fmt.Println(id)

	hc, errr := service.GetHealthcenter(uint(id))

	if errr != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	services, err := service.GetServices(uint(id))

	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	rating, err := service.GetRating(uint(id))
	frating, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", rating), 64)
	fmt.Println("rating: ", rating)

	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	comments, err := service.GetFeedback(hc.ID)

	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.error.layout", nil)
		return
	}

	fmt.Println(comments)
	token, err := rtoken.CSRFToken(uh.CsrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	loginForm := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		CSRF:    token,
	}

	data := struct {
		Rating       float64
		Healthcenter entity.HealthCenter
		Services     []entity.Service
		Comments     []entity.UserComment
		Isvalid      string
		FormValue	form.Input
	}{
		Rating:       frating,
		Healthcenter: *hc,
		Services:     services,
		Comments:     comments,
		FormValue: loginForm,
	}



	//c, err := r.Cookie("user")

	if !uh.LoggedIn(r) {
		uh.Templ.ExecuteTemplate(w, "user.hc.default.layout", data)
		return
	} else {
		//fmt.Println(c.Value)
		//fmt.Println(c.MaxAge)
	}
	//uid, _ := strconv.Atoi(c.Value)
	uid := uh.LoggedInUser.ID
	validity, err := service.CheckValidity(uid, hc.ID)
	fmt.Println(validity)
	fmt.Println(err)


	if err != nil {
		uh.Templ.ExecuteTemplate(w, "user.hc.default.layout", data)
		return
	}
	data.Isvalid = validity
	uh.Templ.ExecuteTemplate(w, "user.hc.auth.layout", data)

}

// Logout hanldes the POST /logout requests
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//UserSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	session.Remove(uh.UserSess.UUID, w)
	service.DeleteSession(uh.UserSess.UUID)
	uh.LoggedInUser = nil
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// Feedback handles POST /feedback
func (uh *UserHandler) Feedback(w http.ResponseWriter, r *http.Request) {
	rating := r.PostFormValue("userrating")
	comment := r.PostFormValue("usercomment")
	hid := r.PostFormValue("hcid")
	fmt.Println("the request has come this far")
	//c, err := r.Cookie("user")
	if !uh.LoggedIn(r) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	fuid := uh.LoggedInUser.ID
	frating, err := strconv.Atoi(rating)

	//fuid, err := strconv.Atoi(uid)
	fhid, err := strconv.Atoi(strings.Trim(hid, " "))
	if err != nil {
		fmt.Println(err)
		return
	}
	feedback := entity.Comment{
		Rating:         uint(frating),
		Comment:        comment,
		UserID:         fuid,
		HealthCenterID: uint(fhid),
	}

	fmt.Println(feedback)
	err = service.PostFeedback(&feedback)

	if err != nil {
		w.Write([]byte("error"))
		return
	}
	w.Write([]byte("success"))
}


func (adh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}
	id,_ := strconv.Atoi(r.FormValue("hidden_id"))
	fmt.Println("deleting")
	URL := fmt.Sprintf("http://localhost:8181/v1/users/%d",id)

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

