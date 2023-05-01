package event

import "github.com/stuton/xm-golang-exercise/internal/companies/model"

type EventType string

const (
	CreateCompanyEventType EventType = "CreateCompany"
	UpdateCompanyEventType EventType = "UpdateCompany"
	DeleteCompanyEventType EventType = "DeleteCompany"
)

type EventMessage struct {
	EventType EventType
	Company   model.Company
}
