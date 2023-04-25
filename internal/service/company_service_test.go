package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apierror "github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/model"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"github.com/stuton/xm-golang-exercise/internal/repository/mocks"
	"github.com/stuton/xm-golang-exercise/internal/service"
)

func TestCompanyService_Get(t *testing.T) {
	repo := &mocks.CompanyRepository{}

	t.Run("getting non-existent company", func(t *testing.T) {

		repo.On("GetCompanyByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, apierror.NewNotFoundError("", "", "company is not found")).
			Once()

		service := service.NewCompanyService(repo, producer.ProducerProcessing{})

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

		service := service.NewCompanyService(repo, producer.ProducerProcessing{})

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

		service := service.NewCompanyService(repo, producer.ProducerProcessing{})

		company, err := service.GetCompanyByID(context.Background(), model.CompanyIDRequest{
			ID: "company1",
		})

		assert.Nil(t, err)
		assert.Equal(t, &model.Company{Name: "company1", Type: "NonProfit", Description: "", AmountEmployees: 20, Registered: false}, company)
	})

}

func TestCompanyService_Create(t *testing.T) {
	repo := &mocks.CompanyRepository{}

	t.Run("company is already exist", func(t *testing.T) {

		repo.On("GetCompanyByName", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, nil).
			Once()

		service := service.NewCompanyService(repo, producer.ProducerProcessing{})

		err := service.CreateCompany(context.Background(), model.CompanyRequest{
			Name: "company1",
		})

		assert.Error(t, err, "unable to create company")
		assert.NotNil(t, err)
	})

	// t.Run("company is already exist", func(t *testing.T) {

	// 	repo.On("GetCompanyByName", mock.Anything, mock.AnythingOfType("string")).
	// 		Return(nil, nil).
	// 		Once()

	// 	repo.On("CreateCompany", mock.Anything, mock.AnythingOfType("model.Company")).
	// 		Return(nil).
	// 		Once()

	// 	service := service.NewCompanyService(repo, producer.ProducerProcessing{})

	// 	err := service.CreateCompany(context.Background(), model.CompanyRequest{
	// 		Name:            "company1",
	// 		Description:     "",
	// 		AmountEmployees: 10,
	// 		Registered:      false,
	// 		Type:            "test",
	// 	})

	// 	assert.Error(t, err)
	// 	assert.NotNil(t, err)
	// })

}
