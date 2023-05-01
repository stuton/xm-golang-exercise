package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stuton/xm-golang-exercise/internal/companies/model"
	"github.com/stuton/xm-golang-exercise/internal/companies/repository/mocks"
	"github.com/stuton/xm-golang-exercise/internal/companies/service"
	apierror "github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"go.uber.org/zap"
)

func TestCompanyService_Get(t *testing.T) {
	repo := &mocks.CompanyRepository{}
	logger := zap.NewExample().Sugar()

	t.Run("getting non-existent company", func(t *testing.T) {

		repo.On("GetCompanyByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, apierror.NewNotFoundError("", "", "company is not found")).
			Once()

		service := service.NewCompanyService(logger, repo, producer.ProducerProcessing{})

		_, err := service.GetCompanyByID(context.Background(), model.CompanyIDRequest{
			ID: "company1",
		})

		assert.Error(t, err, "company is not found")
		assert.NotNil(t, err)
	})

	t.Run("getting unexpected error", func(t *testing.T) {

		repo.On("GetCompanyByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, apierror.NewInternalServerError("", "", "unable to get company")).
			Once()

		service := service.NewCompanyService(logger, repo, producer.ProducerProcessing{})

		_, err := service.GetCompanyByID(context.Background(), model.CompanyIDRequest{
			ID: "company1",
		})

		assert.Error(t, err, "unable to get company")
		assert.NotNil(t, err)
	})

	t.Run("getting existing company", func(t *testing.T) {

		repo.On("GetCompanyByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&model.Company{Name: "company1", Type: "NonProfit", Description: "", AmountEmployees: 20, Registered: false}, nil).
			Once()

		service := service.NewCompanyService(logger, repo, producer.ProducerProcessing{})

		company, err := service.GetCompanyByID(context.Background(), model.CompanyIDRequest{
			ID: "company1",
		})

		assert.Nil(t, err)
		assert.Equal(t, &model.Company{Name: "company1", Type: "NonProfit", Description: "", AmountEmployees: 20, Registered: false}, company)
	})

}

func TestCompanyService_Create(t *testing.T) {
	repo := &mocks.CompanyRepository{}
	logger := zap.NewExample().Sugar()

	t.Run("company is already exist", func(t *testing.T) {

		repo.On("CreateCompany", mock.Anything, mock.AnythingOfType("model.Company")).
			Return("", apierror.NewBadRequestError("", "", "company is already exist")).
			Once()

		service := service.NewCompanyService(logger, repo, producer.ProducerProcessing{})

		var registered bool

		_, err := service.CreateCompany(context.Background(), model.CompanyRequest{
			Name:       "company1",
			Registered: &registered,
		})

		assert.Error(t, err, "company is already exists")
		assert.NotNil(t, err)
	})

}
