package events

import "github.com/imohamedsheta/xapp/app/shared/enums"

/*
| This file will contain the payload of the events.
*/

type UserLoggedInPayload struct {
	UserId          int64               `json:"user_id"`
	AuditableType   enums.AuditableType `json:"auditable_type"`
	ImpersionatedBy *int64              `json:"impersonated_by"`
	Summary         string              `json:"summary"`
}

type UserRegisterPayload struct {
	UserId        int64               `json:"user_id"`
	AuditableType enums.AuditableType `json:"auditable_type"`
	Summary       string              `json:"summary"`
}
