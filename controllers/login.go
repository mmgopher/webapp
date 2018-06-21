package controllers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/julienschmidt/httprouter"
	"webapp/dao"
	 "webapp/config"
)

type LoginController struct {
	ud dao.IUserDAO
}

func NewLoginController(userDao dao.IUserDAO) *LoginController {
	return &LoginController{userDao}
}
func (lc LoginController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	config.Template.ExecuteTemplate(w, "login.gohtml", nil)

}

func (lc LoginController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	config.SessionManager.DestroySession(w,r)
	config.Template.ExecuteTemplate(w, "login.gohtml", nil)
}

func (lc LoginController) LoginProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	l := r.FormValue("login")
	p := r.FormValue("password")
	u,err := lc.ud.GetByLogin(l)
	if( err != nil) {
		http.Error(w, "Username and/or password do not match", http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	if err != nil {
		http.Error(w, "Username and/or password do not match", http.StatusForbidden)
		return
	}

	config.SessionManager.CreateSession(w)
	http.Redirect(w, r, "/users", http.StatusSeeOther)

}

func alreadyLoggedIn(r *http.Request) bool {
	_, err := config.SessionManager.GetSession(r)
	if err != nil {
		return false
	}
	return true
}
