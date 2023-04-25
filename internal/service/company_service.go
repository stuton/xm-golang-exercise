package service

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"github.com/stuton/xm-golang-exercise/internal/repository"
)

type CompanyService interface {
	GetCompanyByID(ctx context.Context, request model.CompanyIDRequest) (*model.Company, error)
	CreateCompany(ctx context.Context, request model.CompanyRequest) error
	UpdateCompany(ctx context.Context, request model.CompanyRequest) error
	DeleteCompanyByID(ctx context.Context, request model.CompanyIDRequest) error
}

type CompanyServiceStruct struct {
	companyRepository  repository.CompanyRepository
	producerProcessing producer.ProducerProcessing
}

func NewCompanyService(
	companyRepository repository.CompanyRepository,
	producerProcessing producer.ProducerProcessing,
) CompanyService {
	return CompanyServiceStruct{
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
		log.Error().Msgf("Unable to get company: %v", err)
		return nil, errors.NewInternalServerError("", nil, "unable to get company")
	}

	return company, nil
}

func (s CompanyServiceStruct) CreateCompany(ctx context.Context, request model.CompanyRequest) error {

	if _, err := s.companyRepository.GetCompanyByName(ctx, request.Name); err == nil {
		return errors.NewBadRequestError("", nil, "company is already exist")
	}

	companyModel := model.Company{
		Name:            request.Name,
		Type:            request.Type,
		Description:     request.Description,
		AmountEmployees: request.AmountEmployees,
		Registered:      *request.Registered,
	}

	if err := s.companyRepository.CreateCompany(ctx, companyModel); err != nil {
		log.Error().Msgf("Unable to create company: %v", err)
		return errors.NewInternalServerError("", nil, "unable to create company")
	}

	createCompanyMessage := &producer.EventMessage{
		EventType: producer.CreateCompanyEventType,
		Company:   companyModel,
	}

	if err := s.producerProcessing.WriteMessages(viper.GetString("KAFKA_TOPIC"), createCompanyMessage); err != nil {
		log.Error().Msgf("Unable to processing message: %v", err)
		return errors.NewInternalServerError("", nil, "unable to create company")
	}

	return nil
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
		log.Error().Msgf("Unable to update company: %v", err)
		return errors.NewInternalServerError("", nil, "unable to update company")
	}

	updateCompanyMessage := &producer.EventMessage{
		EventType: producer.UpdateCompanyEventType,
		Company:   companyModel,
	}

	if err := s.producerProcessing.WriteMessages(viper.GetString("KAFKA_TOPIC"), updateCompanyMessage); err != nil {
		log.Error().Msgf("Unable to processing message: %v", err)
		return errors.NewInternalServerError("", nil, "unable to update company")
	}

	return nil
}

func (s CompanyServiceStruct) DeleteCompanyByID(ctx context.Context, request model.CompanyIDRequest) error {

	if err := s.companyRepository.DeleteCompanyByID(ctx, request.ID); err != nil {
		log.Error().Msgf("Unable to delete company: %v", err)
		return errors.NewInternalServerError("", nil, "unable to delete company")
	}

	createCompanyMessage := &producer.EventMessage{
		EventType: producer.DeleteCompanyEventType,
		Company: model.Company{
			ID: request.ID,
		},
	}

	if err := s.producerProcessing.WriteMessages(viper.GetString("KAFKA_TOPIC"), createCompanyMessage); err != nil {
		log.Error().Msgf("Unable to processing message: %v", err)
		return errors.NewInternalServerError("", nil, "unable to delete company")
	}

	return nil
}
