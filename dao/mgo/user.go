package mgo

import (
	"webapp/config/db"
	"webapp/config/logging"
	"webapp/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserDAO struct {
	session *mgo.Session
}

func NewUserDAO(s *mgo.Session) *UserDAO {
	return &UserDAO{s}
}

func (ud UserDAO) GetById(id string) (models.User, error) {
	u := models.User{}
	if !bson.IsObjectIdHex(id) {
		logging.Error(db.ErrParseId)
		return u, db.ErrParseId
	}
	oid := bson.ObjectIdHex(id)

	s := ud.session.Copy()
	defer s.Close()
	err := s.DB(db.DATABASE).C(db.USER_COLLECTION).FindId(oid).One(&u)
	if err != nil {
		logging.Error(err)
	}
	return u, err
}

func (ud UserDAO) GetByLogin(login string) (models.User, error) {
	s := ud.session.Copy()
	defer s.Close()
	u := models.User{}
	err := s.DB(db.DATABASE).C(db.USER_COLLECTION).Find(bson.M{"login": login}).One(&u)
	if err != nil {
		logging.Error(err)
	}
	return u, err


}

func (ud UserDAO) Create(u models.User) (models.User, error) {
	u.Id = bson.NewObjectId()
	s := ud.session.Copy()
	defer s.Close()
	err := s.DB(db.DATABASE).C(db.USER_COLLECTION).Insert(u)
	if err != nil {
		logging.Error(err)
	}
	return u, err
}

func (ud UserDAO) Update(u models.User) (models.User, error) {
	s := ud.session.Copy()
	defer s.Close()
	var err error
	if u.Password == nil {
		err = s.DB(db.DATABASE).C(db.USER_COLLECTION).UpdateId(u.Id, bson.M{"$set": bson.M{"login": u.Login, "first": u.First, "last": u.Last, "age": u.Age}})
	} else {
		err = s.DB(db.DATABASE).C(db.USER_COLLECTION).UpdateId(u.Id, u)
	}
	if err != nil {
		logging.Error(err)
	}
	return u, err
}

func (ud UserDAO) Delete(id string) error {
	if !bson.IsObjectIdHex(id) {
		logging.Error(db.ErrParseId)
		return db.ErrParseId
	}
	oid := bson.ObjectIdHex(id)
	s := ud.session.Copy()
	defer s.Close()
	err := s.DB(db.DATABASE).C(db.USER_COLLECTION).RemoveId(oid)
	if err != nil {
		logging.Error(err)
	}
	return err
}

func (ud UserDAO) FindAll() ([]models.User, error) {
	s := ud.session.Copy()
	defer s.Close()
	var users []models.User
	err := s.DB(db.DATABASE).C(db.USER_COLLECTION).Find(nil).All(&users)
	if err != nil {
		logging.Error(err)
	}
	return users, err
}
