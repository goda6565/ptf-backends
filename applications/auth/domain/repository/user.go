package repository

import (
	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
