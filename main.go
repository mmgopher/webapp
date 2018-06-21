package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"webapp/controllers"
	"webapp/dao/mgo"
	"log"
	"webapp/config"
)

func main() {
	session := mgo.GetSession()
	defer session.Close()
	userDao := mgo.NewUserDAO(session)
	uc := controllers.NewUserController(userDao)
	lc := controllers.NewLoginController(userDao)
	urc := controllers.NewUserRestController(userDao)

	r := httprouter.New()
	r.GET("/users/create", uc.Create)
	r.GET("/users/update/:id", uc.Get)
	r.GET("/users", uc.FindAll)
	r.GET("/users/delete/:id", uc.Delete)
	r.GET("/login", lc.Login)
	r.POST("/login",lc.LoginProcess)
	r.GET("/logout", lc.Logout)
	r.POST("/users/update",uc.Update)
	r.GET("/rest/user/:id", urc.Get)
	r.DELETE("/rest/user/:id", urc.Delete)
	r.POST("/rest/user", urc.Create)
	r.PUT("/rest/user", urc.Update)
	r.GET("/", index)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}


func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config.Template.ExecuteTemplate(w, "index.gohtml", nil)
//	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
