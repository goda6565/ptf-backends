package repositoryimpl

import (
	"gorm.io/gorm"

	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
	"github.com/goda6565/ptf-backends/applications/auth/domain/repository"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(email string, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Password: password,
	}
	tx := r.db.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	tx := r.db.Where("email = ?", email).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
