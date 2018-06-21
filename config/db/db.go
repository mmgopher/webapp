package db

import (
	"gopkg.in/mgo.v2"
	"errors"
)

const (
	DB_URL          = "mongodb://localhost"
	DATABASE        = "gl-web-app"
	USER_COLLECTION = "users"
)

var (
	ErrNotFound = mgo.ErrNotFound
	ErrParseId = errors.New("id: cannot parse as HEX")
)

