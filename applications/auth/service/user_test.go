package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/utils"
)

type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
}

func (m *mockUserRepository) CreateUser(email string, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

type UserServiceTestSuite struct {
	suite.Suite
	userService UserService
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupTest() {
}

func (suite *UserServiceTestSuite) TestRegister() {
	userRepository := NewMockUserRepository()
	suite.userService = NewUserService(userRepository)
	userRepository.On("CreateUser", "email", mock.AnythingOfType("string")).Return(&models.User{
		Email:    "email",
		Password: "password",
	}, nil)

	user, err := suite.userService.UserRegister("email", "password")
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal("email", user.Email)
	suite.NotEqual(0, user.ID) // 0 でなければ作成されたと見なす
}

func (suite *UserServiceTestSuite) TestLogin() {
	userRepository := NewMockUserRepository()
	suite.userService = NewUserService(userRepository)
	hashed, err := utils.HashPassword("password")
	suite.NoError(err)
	userRepository.On("GetUserByEmail", "email").Return(&models.User{
		Email:    "email",
		Password: hashed,
	}, nil)
	token, err := suite.userService.UserLogin("email", "password")
	suite.NoError(err)
	suite.NotEmpty(token)
}

func (suite *UserServiceTestSuite) TestGetUser() {
	userRepository := NewMockUserRepository()
	suite.userService = NewUserService(userRepository)
	hashed, err := utils.HashPassword("password")
	suite.NoError(err)
	userRepository.On("GetUserByEmail", "email").Return(&models.User{
		Email:    "email",
		Password: hashed,
	}, nil)
	token, err := suite.userService.UserLogin("email", "password")
	suite.NoError(err)
	suite.NotEmpty(token)

	userRepository.On("GetUserByEmail", "email").Return(&models.User{
		Email:    "email",
		Password: hashed,
	}, nil)
	user, err := suite.userService.GetUser(token)
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal("email", user.Email)
}
