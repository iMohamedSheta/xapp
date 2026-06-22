package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
)

type User struct {
	Id              int64            `xqb:"id" json:"id"`
	TenantId        int64            `xqb:"tenant_id" json:"tenant_id"`
	ClientId        sql.NullInt64    `xqb:"client_id" json:"client_id"`
	ClientType      sql.NullString   `xqb:"client_type" json:"client_type"` // Can be USER || CARD only
	Username        sql.NullString   `xqb:"username" json:"username"`
	Name            string           `xqb:"name" json:"name"`
	Email           sql.NullString   `xqb:"email" json:"email"`
	Password        string           `xqb:"password" json:"-"`
	Provider        sql.NullString   `xqb:"provider" json:"provider"`
	ProviderId      sql.NullString   `xqb:"provider_id" json:"provider_id"`
	Avatar          sql.NullString   `xqb:"avatar" json:"avatar"`
	EmailVerifiedAt sql.NullTime     `xqb:"email_verified_at" json:"email_verified_at"`
	Role            enums.UserRole   `xqb:"role" json:"role"`
	Status          enums.UserStatus `xqb:"status" json:"status"`
	DeletedAt       sql.NullTime     `xqb:"deleted_at" json:"deleted_at"`
	CreatedAt       time.Time        `xqb:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `xqb:"updated_at" json:"updated_at"`

	// Relationships
	Tenant *Tenant `table:"tenants" json:"tenant"`
	BaseModel
}

func (User) Table() string {
	return "users"
}

func (u *User) GetNotifiableID() int64 {
	return u.Id
}

func (User) GetNotifiableType() string {
	return "user"
}

func (m User) Cols() []any {
	return m.BaseModel.Cols(m, "user")
}

func (u *User) LoadTenant(c context.Context) error {
	if u.Tenant != nil && u.Tenant.Id != 0 {
		return nil
	}

	tenantData, err := xqb.Model[Tenant]().
		WithContext(c).
		Where("id", "=", u.TenantId).
		First()

	if err != nil {
		return err
	}

	u.Tenant = tenantData
	return nil
}

func (u *User) IsBlockedFromLogin() bool {
	return u.Status == enums.UserStatusBlocked
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User

	return json.Marshal(&struct {
		ClientId        any `json:"client_id,omitempty"`
		ClientType      any `json:"client_type,omitempty"`
		Username        any `json:"username,omitempty"`
		Email           any `json:"email,omitempty"`
		Provider        any `json:"provider,omitempty"`
		ProviderId      any `json:"provider_id,omitempty"`
		Avatar          any `json:"avatar,omitempty"`
		EmailVerifiedAt any `json:"email_verified_at,omitempty"`
		DeletedAt       any `json:"deleted_at,omitempty"`
		*Alias
	}{
		ClientId:        nullInt64ToAny(u.ClientId),
		ClientType:      nullStringToAny(u.ClientType),
		Username:        nullStringToAny(u.Username),
		Email:           nullStringToAny(u.Email),
		Provider:        nullStringToAny(u.Provider),
		ProviderId:      nullStringToAny(u.ProviderId),
		Avatar:          nullStringToAny(u.Avatar),
		EmailVerifiedAt: nullTimeToAny(u.EmailVerifiedAt),
		DeletedAt:       nullTimeToAny(u.DeletedAt),
		Alias:           (*Alias)(&u),
	})
}
