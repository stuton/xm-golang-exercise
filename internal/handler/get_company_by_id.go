package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

type getCompanyByIDHandler struct {
	companyService service.CompanyService
}

func NewGetCompanyByIDHandler(companyService service.CompanyService) getCompanyByIDHandler {
	return getCompanyByIDHandler{
		companyService: companyService,
	}
}

func (handler getCompanyByIDHandler) GetCompanyByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		company, err := handler.companyService.GetCompanyByID(c, companyIDRequest)

		if err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(company.ResponseMapping())
	}
}
