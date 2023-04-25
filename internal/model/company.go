package model

type CompanyType string

const (
	CorporationsType       CompanyType = "Corporations"
	NonProfitType          CompanyType = "NonProfit"
	CooperativeType        CompanyType = "Cooperative"
	SoleProprietorshipType CompanyType = "SoleProprietorship"
)

type Company struct {
	ID              string      `db:"id"`
	Name            string      `db:"name"`
	Type            CompanyType `db:"type"`
	Description     string      `db:"description"`
	AmountEmployees int         `db:"amount_employees"`
	Registered      bool        `db:"registered"`
}

type CompanyResponse struct {
	ID              string      `db:"id"`
	Name            string      `db:"name"`
	Type            CompanyType `db:"type"`
	Description     string      `db:"description"`
	AmountEmployees int         `db:"amount_employees"`
}

// Mapping response if some fields we would like to hide (e.x Registered)
func (c *Company) ResponseMapping() *CompanyResponse {
	if c == nil {
		return nil
	}
	return &CompanyResponse{
		ID:              c.ID,
		Name:            c.Name,
		Type:            c.Type,
		Description:     c.Description,
		AmountEmployees: c.AmountEmployees,
	}
}

type CompanyIDRequest struct {
	ID string `uri:"id" binding:"uuid4_rfc4122"`
}

type CompanyRequest struct {
	ID              string      `json:"-"`
	Name            string      `json:"name" binding:"required,max=15"`
	Description     string      `json:"description" binding:"max=3000"`
	AmountEmployees int         `json:"amount_employees" binding:"required,gte=0"`
	Registered      *bool       `json:"registered" binding:"required"`
	Type            CompanyType `json:"type" binding:"required"`
}
