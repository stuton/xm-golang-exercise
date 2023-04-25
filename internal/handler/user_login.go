package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

type userLoginHandler struct {
	userService service.UserService
}

func NewUserLoginHandler(userService service.UserService) userLoginHandler {
	return userLoginHandler{
		userService: userService,
	}
}

func (handler userLoginHandler) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var userLoginRequest service.UserLoginRequest

		if err := c.ShouldBindJSON(&userLoginRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		token, err := handler.userService.Login(c, userLoginRequest)

		if err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"token": token})
	}
}
