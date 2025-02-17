package service

import (
	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
	"github.com/goda6565/ptf-backends/applications/auth/domain/repository"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/utils"
)

type UserService interface {
	UserRegister(email string, password string) (*models.User, error)
	UserLogin(email string, password string) (token string, err error)
	GetUser(token string) (*models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) UserRegister(email string, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.userRepository.CreateUser(email, hashedPassword)
}

func (s *userService) UserLogin(email string, password string) (token string, err error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	err = utils.CheckPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	return utils.GenerateSignedString(user.ID, user.Email)
}

func (s *userService) GetUser(token string) (*models.User, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.userRepository.GetUserByEmail(claims.Email)
}
