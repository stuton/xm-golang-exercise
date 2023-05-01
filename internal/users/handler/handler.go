package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/users/model"
	"github.com/stuton/xm-golang-exercise/internal/users/service"
	"go.uber.org/zap"
)

type UsersHandlers struct {
	UserLoginHandler UserLoginHandler
}

type UserLoginHandler struct {
	service service.UserService
	logger  *zap.SugaredLogger
}

func NewUsersHandlers(service service.UserService, logger *zap.SugaredLogger) UsersHandlers {
	return UsersHandlers{
		UserLoginHandler: NewUserLoginHandler(service, logger),
	}
}

func NewUserLoginHandler(userService service.UserService, logger *zap.SugaredLogger) UserLoginHandler {
	return UserLoginHandler{
		service: userService,
		logger:  logger,
	}
}

func (h UserLoginHandler) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var userLoginRequest model.UserLoginRequest

		if err := c.ShouldBindJSON(&userLoginRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		token, err := h.service.Login(c, userLoginRequest)

		if err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"token": token})
	}
}
