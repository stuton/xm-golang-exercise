package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stuton/xm-golang-exercise/internal/companies/model"
	"github.com/stuton/xm-golang-exercise/internal/companies/service"
	"github.com/stuton/xm-golang-exercise/internal/server/http/request"
	"go.uber.org/zap"
)

type CompaniesHandlers struct {
	GetCompanyByIDHandler    GetCompanyByIDHandler
	CreateCompanyHandler     CreateCompanyHandler
	UpdateCompanyHandler     UpdateCompanyHandler
	DeleteCompanyByIDHandler DeleteCompanyByIDHandler
}

type GetCompanyByIDHandler struct {
	service service.CompanyService
	logger  *zap.SugaredLogger
}

type CreateCompanyHandler struct {
	service service.CompanyService
	logger  *zap.SugaredLogger
}
type UpdateCompanyHandler struct {
	service service.CompanyService
	logger  *zap.SugaredLogger
}

type DeleteCompanyByIDHandler struct {
	service service.CompanyService
	logger  *zap.SugaredLogger
}

func NewCompaniesHandlers(service service.CompanyService, logger *zap.SugaredLogger) CompaniesHandlers {
	return CompaniesHandlers{
		GetCompanyByIDHandler:    NewGetCompanyByIDHandler(service, logger),
		CreateCompanyHandler:     NewCreateCompanyHandler(service, logger),
		UpdateCompanyHandler:     NewUpdateCompanyHandler(service, logger),
		DeleteCompanyByIDHandler: NewDeleteCompanyByIDHandler(service, logger),
	}
}

func NewGetCompanyByIDHandler(service service.CompanyService, logger *zap.SugaredLogger) GetCompanyByIDHandler {
	return GetCompanyByIDHandler{
		service: service,
		logger:  logger,
	}
}

func NewCreateCompanyHandler(service service.CompanyService, logger *zap.SugaredLogger) CreateCompanyHandler {
	return CreateCompanyHandler{
		service: service,
		logger:  logger,
	}
}

func NewUpdateCompanyHandler(service service.CompanyService, logger *zap.SugaredLogger) UpdateCompanyHandler {
	return UpdateCompanyHandler{
		service: service,
		logger:  logger,
	}
}

func NewDeleteCompanyByIDHandler(service service.CompanyService, logger *zap.SugaredLogger) DeleteCompanyByIDHandler {
	return DeleteCompanyByIDHandler{
		service: service,
		logger:  logger,
	}
}

func (h GetCompanyByIDHandler) GetCompanyByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		company, err := h.service.GetCompanyByID(c, companyIDRequest)

		if err != nil {
			cc.ResponseError(err)
			return
		}

		cc.Ok(company.ResponseMapping())
	}
}

func (h CreateCompanyHandler) CreateCompany() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyRequest model.CompanyRequest

		if err := c.ShouldBindJSON(&companyRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		uuid, err := h.service.CreateCompany(c, companyRequest)

		if err != nil {
			h.logger.Error("unable to create company", zap.Error(err))
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"uuid": uuid})
	}
}

func (h UpdateCompanyHandler) UpdateCompany() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest
		var companyRequest model.CompanyRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		if err := c.ShouldBindJSON(&companyRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		companyRequest.ID = companyIDRequest.ID

		if err := h.service.UpdateCompany(c, companyRequest); err != nil {
			h.logger.Error("unable to update company", zap.Error(err))
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"success": true})

	}
}

func (h DeleteCompanyByIDHandler) DeleteCompanyByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := request.FromContext(c)

		var companyIDRequest model.CompanyIDRequest

		if err := c.ShouldBindUri(&companyIDRequest); err != nil {
			cc.BadRequest(err)
			return
		}

		if err := h.service.DeleteCompanyByID(c, companyIDRequest); err != nil {
			h.logger.Error("unable to delete company", zap.Error(err))
			cc.ResponseError(err)
			return
		}

		cc.Ok(gin.H{"success": true})
	}
}
