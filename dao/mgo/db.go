package mgo

import (
	"gopkg.in/mgo.v2"
	"webapp/config/db"
)

func GetSession() *mgo.Session {
	s, err := mgo.Dial(db.DB_URL)

	if err != nil {
		panic(err)
	}
	return s
}