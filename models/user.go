package models

import "gopkg.in/mgo.v2/bson"

type (
	User struct {
		Id       bson.ObjectId `json:"id" bson:"_id"`
		Login    string        `json:"login" bson:"login"`
		Password []byte        `json:"password" bson:"password"`
		First    string        `json:"first" bson:"first"`
		Last     string        `json:"last" bson:"last"`
		Age      int           `json:"age" bson:"age"`
	}
)
