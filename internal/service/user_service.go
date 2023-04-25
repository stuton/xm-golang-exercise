package service

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/repository"
	"github.com/stuton/xm-golang-exercise/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (s UserService) Login(ctx context.Context, request UserLoginRequest) (string, error) {

	user, err := s.userRepository.GetUserByUsername(ctx, request.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.NewNotFoundError("", nil, "user is not found")
		}
		log.Error().Msgf("Unable to get user by username: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to login user")
	}

	if CheckPasswordHash(request.Password, user.Password) {
		return "", errors.NewBadRequestError("", nil, "login or password is not correct")
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		log.Error().Msgf("Unable to generate token: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to login user")
	}

	return token, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil && err == bcrypt.ErrMismatchedHashAndPassword
}
