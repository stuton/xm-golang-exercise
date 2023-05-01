package service

import (
	"context"
	"database/sql"

	"github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/users/model"
	"github.com/stuton/xm-golang-exercise/internal/users/repository"
	"github.com/stuton/xm-golang-exercise/utils/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger         *zap.SugaredLogger
	userRepository repository.UserRepository
	token          jwt.JWT
}

func NewUserService(
	logger *zap.SugaredLogger,
	userRepository repository.UserRepository,
	token jwt.JWT,
) UserService {
	return UserService{
		logger:         logger,
		userRepository: userRepository,
		token:          token,
	}
}

func (s UserService) Login(ctx context.Context, request model.UserLoginRequest) (string, error) {

	user, err := s.userRepository.GetUserByUsername(ctx, request.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.NewNotFoundError("", nil, "user is not found")
		}
		s.logger.Errorf("unable to get user by username: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to login user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return "", errors.NewBadRequestError("", nil, "password is not correct")
	}

	token, err := s.token.GenerateToken(user.ID)

	if err != nil {
		s.logger.Errorf("unable to generate token: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to login user")
	}

	return token, nil
}
