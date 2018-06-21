package dao

import "webapp/models"

type IUserDAO interface {
	GetById(id string) (models.User, error)
	Create(u models.User) (models.User, error)
	Update(u models.User) (models.User, error)
	Delete(id string) (error)
	FindAll() ([]models.User,error)
	GetByLogin(login string) (models.User, error)
}

