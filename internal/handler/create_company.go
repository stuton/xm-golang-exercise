package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

type createCompanyHandler struct {
	companyService service.CompanyService
}

func NewCreateCompanyHandler(companyService service.CompanyService) createCompanyHandler {
	return createCompanyHandler{
		companyService: companyService,
	}
}

func (handler createCompanyHandler) CreateCompany() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyRequest model.CompanyRequest

		if err := c.ShouldBindJSON(&companyRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		if err := handler.companyService.CreateCompany(c, companyRequest); err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"success": true})
	}
}
