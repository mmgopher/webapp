package controllers

import (
	"webapp/dao"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"webapp/config/logging"
	"webapp/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRestController struct {
	ud dao.IUserDAO
}
var success = map[string]string{"result": "success"}
func NewUserRestController(userDao dao.IUserDAO) *UserRestController {
	return &UserRestController{userDao}
}

func (uc UserRestController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	u, err := uc.ud.GetById(id)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	responseWithJson(w,http.StatusOK,u)

}

func (uc UserRestController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	err := uc.ud.Delete(id)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	responseWithJson(w,http.StatusOK,success)
}

func (uc UserRestController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logging.Info("REST CREATE")
	u := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logging.Error(err)
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	bp, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		responseWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	u.Password = bp

	u, errd := uc.ud.Create(u)
	if(errd != nil) {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	responseWithJson(w,http.StatusCreated,u)
}

func (uc UserRestController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logging.Info("REST UPDATES")
	u := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logging.Error(err)
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	u, errd := uc.ud.Update(u)
	if(errd != nil) {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	responseWithJson(w,http.StatusCreated,u)
}

func responseWithError(w http.ResponseWriter, code int,  msg string) {
	responseWithJson(w, code, map[string]string{"error": msg})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if(err != nil) {
		logging.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}