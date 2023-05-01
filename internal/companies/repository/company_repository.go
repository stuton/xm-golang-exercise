package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"

	"github.com/stuton/xm-golang-exercise/internal/companies/model"
)

var (
	ErrKeyExists = errors.New("company is already exists")
)

//go:generate mockery --name CompanyRepository --filename company_repository.go
type CompanyRepository interface {
	GetCompanyByID(ctx context.Context, id string) (*model.Company, error)
	CreateCompany(ctx context.Context, params model.Company) (string, error)
	UpdateCompany(ctx context.Context, params model.Company) error
	DeleteCompanyByID(ctx context.Context, id string) error
}

type companyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) CompanyRepository {
	return companyRepository{
		db: db,
	}
}

// GetCompanyByID implements Company
func (repo companyRepository) GetCompanyByID(ctx context.Context, id string) (*model.Company, error) {
	c := new(model.Company)

	if err := repo.db.Get(c, "SELECT * FROM companies WHERE id = $1", id); err != nil {
		return nil, err
	}

	return c, nil
}

// Create implements Company
func (repo companyRepository) CreateCompany(ctx context.Context, params model.Company) (string, error) {
	result, err := repo.db.NamedQuery(`
		INSERT INTO companies (
			name, 
			description, 
			amount_employees, 
			registered, 
			type
		) 
		VALUES (
			:name, 
			:description, 
			:amount_employees, 
			:registered, 
			:type
		) RETURNING id`, params)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", ErrKeyExists
		}
	}

	var id string

	for result.Next() {
		if err = result.Scan(&id); err != nil {
			return "", result.Err()
		}
	}

	return id, result.Err()

}

// Patch implements Company
func (repo companyRepository) UpdateCompany(ctx context.Context, params model.Company) error {
	_, err := repo.db.NamedExec(`
		UPDATE companies SET (
			description, 
			amount_employees, 
			registered, 
			type
		) = (
			:description, 
			:amount_employees, 
			:registered, 
			:type
		) WHERE id = :id`,
		map[string]interface{}{
			"id":               params.ID,
			"description":      params.Description,
			"amount_employees": params.AmountEmployees,
			"registered":       params.Registered,
			"type":             params.Type,
		})

	return err
}

// Delete implements Company
func (repo companyRepository) DeleteCompanyByID(ctx context.Context, id string) error {
	_, err := repo.db.NamedExec(`DELETE FROM companies WHERE id = :id`, map[string]interface{}{
		"id": id,
	})

	return err
}
