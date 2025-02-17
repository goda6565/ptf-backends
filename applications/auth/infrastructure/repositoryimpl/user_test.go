package repositoryimpl_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goda6565/ptf-backends/applications/auth/domain/repository"
	"github.com/goda6565/ptf-backends/applications/auth/infrastructure/repositoryimpl"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/tester"
)

type UserTestSuite struct {
	tester.DBSQLiteSuite
	userRepo repository.UserRepository
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite() // 親スイートの SetupSuite を実行し、SQLite を初期化
	suite.userRepo = repositoryimpl.NewUserRepository(suite.DB)
}

func (suite *UserTestSuite) TestCreateUser() {
	createdUser, err := suite.userRepo.CreateUser("email1", "password")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(createdUser)
	suite.Assert().Equal("email1", createdUser.Email)
	suite.Assert().Equal("password", createdUser.Password)
	suite.Assert().NotEmpty(createdUser.ID)
	suite.Assert().NotEmpty(createdUser.CreatedAt)
}

func (suite *UserTestSuite) TestGetUserByEmail() {
	createdUser, err := suite.userRepo.CreateUser("email2", "password")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(createdUser)

	user, err := suite.userRepo.GetUserByEmail("email2")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(user)
	suite.Assert().Equal("email2", user.Email)
	suite.Assert().Equal("password", user.Password)
	suite.Assert().NotEmpty(user.ID)
	suite.Assert().NotEmpty(user.CreatedAt)
}
