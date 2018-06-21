package controllers

import (
	"github.com/julienschmidt/httprouter"
	"webapp/config"
	"webapp/config/db"
	"webapp/dao"
	"webapp/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"webapp/config/logging"
)

type UserController struct {
	ud dao.IUserDAO
}

func NewUserController(userDao dao.IUserDAO) *UserController {
	return &UserController{userDao}
}

func (uc UserController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := p.ByName("id")
	u, err := uc.ud.GetById(id)
	if validateDbError(w, r, err) {
		config.Template.ExecuteTemplate(w, "update_user.gohtml", u)
	}

}

func validateDbError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch {
	case err == db.ErrNotFound:
		http.NotFound(w, r)
		return false
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return false
	}

	return true
}

func (uc UserController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := r.FormValue("id")
	if !alreadyLoggedIn(r) && id != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := models.User{}

	u.Login = r.FormValue("login")
	u.First = r.FormValue("first")
	u.Last = r.FormValue("last")
	u.Age, _ = strconv.Atoi(r.FormValue("age"))
	password := r.FormValue("password")
	logging.Info("PASSS ",password)
	if password != "" {
		bp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		u.Password = bp
	}
	ufd, err := uc.ud.GetByLogin(u.Login)
	if id == "" {
		if err == nil {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		uc.ud.Create(u)
		config.Template.ExecuteTemplate(w, "login.gohtml", "Account successfully created")
	} else {
		u.Id = bson.ObjectIdHex(id)
		if err == nil {
			if ufd.Id != u.Id {
				http.Error(w, "Username already taken", http.StatusForbidden)
				return
			}
		}

		uc.ud.Update(u)
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (uc UserController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := p.ByName("id")
	err := uc.ud.Delete(id)
	if validateDbError(w, r, err) {
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}

}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	config.Template.ExecuteTemplate(w, "create_user.gohtml", nil)
}
func (uc UserController) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := uc.ud.FindAll()
	if validateDbError(w, r, err) {
		config.Template.ExecuteTemplate(w, "users.gohtml", u)
	}
}
