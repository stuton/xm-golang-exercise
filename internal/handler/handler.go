package handler

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	Login() func(ctx *gin.Context)

	GetCompanyByID() func(ctx *gin.Context)
	UpdateCompany() func(ctx *gin.Context)
	CreateCompany() func(ctx *gin.Context)
	DeleteCompanyByID() func(ctx *gin.Context)
}
