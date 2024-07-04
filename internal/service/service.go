package service

import "effective_mobile_testing/internal/model"

type UserTaskService struct {
	repo Repository
}

func NewUserTaskService(repo Repository) *UserTaskService {
	return &UserTaskService{repo: repo}
}

type Repository interface {
	CreateUser(surname, name, patronymic, address, passportNumber string) (*model.User, error)
	StartTask(userId int64, taskName string) error
	StopTask(userId int64, taskName string) error
	GetLaborCosts(userId int64) ([]model.ResponseLobarCost, error)
	CheckUserIDPerson(userID int64) error
	CheckUserIDTask(id int64) error
	GetUserByFilters(limit, offset int, id int64, surname, name, patronymic, address, passportNumber string) (*[]model.User, error)
	DeleteUser(id int64) error
	UpdateUser(id int64, surname, name, patronymic, address, passportNumber string) (*model.User, error)
}
