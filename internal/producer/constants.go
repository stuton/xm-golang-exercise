package producer

type EventType string

const (
	CreateCompanyEventType EventType = "CreateCompany"
	UpdateCompanyEventType EventType = "UpdateCompany"
	DeleteCompanyEventType EventType = "DeleteCompany"
)
