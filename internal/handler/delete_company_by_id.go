package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

type deleteCompanyByIDHandler struct {
	companyService service.CompanyService
}

func NewDeleteCompanyByIDHandler(companyService service.CompanyService) deleteCompanyByIDHandler {
	return deleteCompanyByIDHandler{
		companyService: companyService,
	}
}

func (handler deleteCompanyByIDHandler) DeleteCompanyByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		if err := handler.companyService.DeleteCompanyByID(c, companyIDRequest); err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"success": true})
	}
}
