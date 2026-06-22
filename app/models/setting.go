package models

import (
	"database/sql"
	"time"
)

type TenantSettings struct {
	SupportPhones []string `json:"support_phones"`
}

type Setting struct {
	Id        int64          `xqb:"id" json:"id"`
	Type      string         `xqb:"type" json:"type"`
	Model     sql.NullString `xqb:"model" json:"model"`
	ModelId   sql.NullInt64  `xqb:"model_id" json:"model_id"`
	Settings  map[string]any `xqb:"settings" json:"settings"`
	CreatedAt time.Time      `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time      `xqb:"updated_at" json:"updated_at"`
	BaseModel
}

func (Setting) Table() string {
	return "settings"
}

func (m Setting) Cols() []any {
	return m.BaseModel.Cols(m, "settings")
}
