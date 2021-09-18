package dao

import (
	"com/justin/micro/kit/library-user-service/models"
	"com/justin/micro/kit/pkg/databases"
)

type UserDao interface {
	SelectByID(id uint64) (*models.User, error)
	SelectByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}

type UserDaoImpl struct{}

func NewUserDaoImpl() UserDao {
	return &UserDaoImpl{}
}

func (u *UserDaoImpl) SelectByID(id uint64) (*models.User, error) {
	user := &models.User{}
	err := databases.DB.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDaoImpl) SelectByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := databases.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDaoImpl) Save(user *models.User) error {
	return databases.DB.Create(user).Error
}
