package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

type updateCompanyHandler struct {
	companyService service.CompanyService
}

func NewUpdateCompanyHandler(companyService service.CompanyService) updateCompanyHandler {
	return updateCompanyHandler{
		companyService: companyService,
	}
}

func (handler updateCompanyHandler) UpdateCompany() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest
		var companyRequest model.CompanyRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			fmt.Println(err)
			cc.BadRequest(err)
			return
		}

		if err := c.ShouldBindJSON(&companyRequest); err != nil {
			fmt.Println(err)
			cc.BadRequest(err)
			return
		}

		companyRequest.ID = companyIDRequest.ID

		if err := handler.companyService.UpdateCompany(c, companyRequest); err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"success": true})

	}
}
