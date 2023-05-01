package service

import (
	"context"
	"database/sql"

	"github.com/stuton/xm-golang-exercise/internal/companies/event"
	"github.com/stuton/xm-golang-exercise/internal/companies/model"
	"github.com/stuton/xm-golang-exercise/internal/companies/repository"
	"github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"go.uber.org/zap"
)

type CompanyService interface {
	GetCompanyByID(ctx context.Context, request model.CompanyIDRequest) (*model.Company, error)
	CreateCompany(ctx context.Context, request model.CompanyRequest) (string, error)
	UpdateCompany(ctx context.Context, request model.CompanyRequest) error
	DeleteCompanyByID(ctx context.Context, request model.CompanyIDRequest) error
}

type CompanyServiceStruct struct {
	logger             *zap.SugaredLogger
	companyRepository  repository.CompanyRepository
	producerProcessing producer.ProducerProcessing
}

func NewCompanyService(
	logger *zap.SugaredLogger,
	companyRepository repository.CompanyRepository,
	producerProcessing producer.ProducerProcessing,
) CompanyService {
	return CompanyServiceStruct{
		logger:             logger,
		companyRepository:  companyRepository,
		producerProcessing: producerProcessing,
	}
}

func (s CompanyServiceStruct) GetCompanyByID(ctx context.Context, request model.CompanyIDRequest) (*model.Company, error) {

	company, err := s.companyRepository.GetCompanyByID(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("", nil, "company is not found")
		}
		s.logger.Errorf("Unable to get company: %v", err)
		return nil, errors.NewInternalServerError("", nil, "unable to get company")
	}

	return company, nil
}

func (s CompanyServiceStruct) CreateCompany(ctx context.Context, request model.CompanyRequest) (string, error) {
	companyModel := model.Company{
		Name:            request.Name,
		Type:            request.Type,
		Description:     request.Description,
		AmountEmployees: request.AmountEmployees,
		Registered:      *request.Registered,
	}

	// TODO: check db code if company is already exist
	uuid, err := s.companyRepository.CreateCompany(ctx, companyModel)

	if err != nil {
		if err == repository.ErrKeyExists {
			return "", errors.NewBadRequestError("", nil, err.Error())
		}
		s.logger.Errorf("unable to create company: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to create company")
	}

	createCompanyMessage := &event.EventMessage{
		EventType: event.CreateCompanyEventType,
		Company:   companyModel,
	}

	if err := s.producerProcessing.WriteMessages(createCompanyMessage); err != nil {
		s.logger.Errorf("Unable to processing message: %v", err)
		return "", errors.NewInternalServerError("", nil, "unable to create company")
	}

	return uuid, nil
}

func (s CompanyServiceStruct) UpdateCompany(ctx context.Context, request model.CompanyRequest) error {

	companyModel := model.Company{
		ID:              request.ID,
		Type:            request.Type,
		Description:     request.Description,
		AmountEmployees: request.AmountEmployees,
		Registered:      *request.Registered,
	}

	if err := s.companyRepository.UpdateCompany(ctx, companyModel); err != nil {
		s.logger.Errorf("Unable to update company: %v", err)
		return errors.NewInternalServerError("", nil, "unable to update company")
	}

	updateCompanyMessage := &event.EventMessage{
		EventType: event.UpdateCompanyEventType,
		Company:   companyModel,
	}

	if err := s.producerProcessing.WriteMessages(updateCompanyMessage); err != nil {
		s.logger.Errorf("Unable to processing message: %v", err)
		return errors.NewInternalServerError("", nil, "unable to update company")
	}

	return nil
}

func (s CompanyServiceStruct) DeleteCompanyByID(ctx context.Context, request model.CompanyIDRequest) error {

	if err := s.companyRepository.DeleteCompanyByID(ctx, request.ID); err != nil {
		s.logger.Errorf("Unable to delete company: %v", err)
		return errors.NewInternalServerError("", nil, "unable to delete company")
	}

	deleteCompanyMessage := &event.EventMessage{
		EventType: event.DeleteCompanyEventType,
		Company: model.Company{
			ID: request.ID,
		},
	}

	if err := s.producerProcessing.WriteMessages(deleteCompanyMessage); err != nil {
		s.logger.Errorf("Unable to processing message: %v", err)
		return errors.NewInternalServerError("", nil, "unable to delete company")
	}

	return nil
}
