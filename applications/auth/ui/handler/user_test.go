package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
	"github.com/goda6565/ptf-backends/applications/auth/ui/gen"
)

type mockUserService struct {
	mock.Mock
}

func NewMockUserService() *mockUserService {
	return &mockUserService{}
}

func (m *mockUserService) UserRegister(email string, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) UserLogin(email string, password string) (token string, err error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *mockUserService) GetUser(token string) (*models.User, error) {
	args := m.Called(token)
	return args.Get(0).(*models.User), args.Error(1)
}

type UserHandlerTestSuite struct {
	suite.Suite
	userHandler *UserHandler
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) SetupTest() {
}

func (suite *UserHandlerTestSuite) TestRegister() {
	userService := NewMockUserService()
	userService.On("UserRegister", "email", "password").Return(&models.User{
		ID:        1,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now(),
	}, nil)
	suite.userHandler = NewUserHandler(userService)

	request, _ := api.NewUserRegisterRequest("api", api.UserRegisterJSONRequestBody{
		Email:    "email",
		Password: "password",
	})

	w := httptest.NewRecorder()               // レスポンスを記録するためのレコーダーを作成
	ginContext, _ := gin.CreateTestContext(w) // テスト用のgin.Contextを作成
	ginContext.Request = request              // リクエストをgin.Contextに設定

	suite.userHandler.UserRegister(ginContext) // テスト対象のメソッドを実行

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body) // レスポンスのボディを読み込む
	var userRegisterResponse api.RegisterResponse
	err := json.Unmarshal(bodyBytes, &userRegisterResponse) // レスポンスのボディを構造体に変換
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal(1, userRegisterResponse.Id)
	suite.Assert().Equal("email", userRegisterResponse.Email)
}

func (suite *UserHandlerTestSuite) TestLogin() {
	userService := NewMockUserService()
	userService.On("UserLogin", "email", "password").Return("token", nil)
	suite.userHandler = NewUserHandler(userService)

	request, _ := api.NewUserLoginRequest("api", api.UserLoginJSONRequestBody{
		Email:    "email",
		Password: "password",
	})

	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.UserLogin(ginContext)

	suite.Assert().Equal(http.StatusOK, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var token string
	err := json.Unmarshal(bodyBytes, &token)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("token", token)
}

func (suite *UserHandlerTestSuite) TestGetUser() {
	userService := NewMockUserService()
	userService.On("GetUser", "token").Return(&models.User{
		ID:        1,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
	}, nil)
	suite.userHandler = NewUserHandler(userService)

	request, _ := api.NewGetUserRequest("api")
	request.Header.Set("Authorization", "Bearer token")

	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.GetUser(ginContext)

	suite.Assert().Equal(http.StatusOK, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var userResponse api.UserResponse
	err := json.Unmarshal(bodyBytes, &userResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal(1, userResponse.Id)
	suite.Assert().Equal("email", userResponse.Email)
	suite.Assert().Equal(time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), userResponse.CreatedAt)
}
