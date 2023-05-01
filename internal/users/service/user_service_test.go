package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stuton/xm-golang-exercise/internal/users/model"
	"github.com/stuton/xm-golang-exercise/internal/users/repository/mocks"
	"github.com/stuton/xm-golang-exercise/internal/users/service"
	"github.com/stuton/xm-golang-exercise/utils/jwt"
	"go.uber.org/zap"
)

func TestUserService_Login(t *testing.T) {
	repo := &mocks.UserRepository{}
	logger := zap.NewExample().Sugar()

	t.Run("successful", func(t *testing.T) {

	})

	t.Run("user is not found", func(t *testing.T) {

		repo.On("GetUserByUsername", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, sql.ErrNoRows).
			Once()

		service := service.NewUserService(logger, repo, jwt.JWT{})

		_, err := service.Login(context.Background(), model.UserLoginRequest{
			Username: "company1",
			Password: "password",
		})

		assert.ErrorContains(t, err, "user is not found")
		assert.NotNil(t, err)
	})

	t.Run("password is not correct", func(t *testing.T) {

		repo.On("GetUserByUsername", mock.Anything, mock.AnythingOfType("string")).
			Return(&model.User{Password: "password1"}, nil).
			Once()

		service := service.NewUserService(logger, repo, jwt.JWT{})

		_, err := service.Login(context.Background(), model.UserLoginRequest{
			Password: "password",
		})

		assert.ErrorContains(t, err, "password is not correct")
		assert.NotNil(t, err)
	})

}
