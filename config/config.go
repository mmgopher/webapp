package config

import (
	"html/template"
	"webapp/session"
	_ "webapp/session/memory"
	"io/ioutil"
	"os"
	"webapp/config/logging"
)

var Template *template.Template
var SessionManager *session.Manager
func init() {
	Template = template.Must(template.ParseGlob("templates/*.gohtml"))
	logging.InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	SessionManager = session.NewManager("memory", "gosessionid", 3600)
	go SessionManager.GC()

}
