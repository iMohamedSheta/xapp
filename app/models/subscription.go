package models

import (
	"database/sql"
	"math"
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
)

type Subscription struct {
	Id              int64                    `xqb:"id" json:"id"`
	TenantId        int64                    `xqb:"tenant_id" json:"tenant_id"`
	PlanId          int64                    `xqb:"plan_id" json:"plan_id"`
	StartDate       time.Time                `xqb:"start_date" json:"start_date"`
	EndDate         time.Time                `xqb:"end_date" json:"end_date"`
	Status          enums.SubscriptionStatus `xqb:"status" json:"status"`
	Price           float64                  `xqb:"price" json:"price"`
	OriginalPrice   float64                  `xqb:"original_price" json:"original_price"`
	Currency        string                   `xqb:"currency" json:"currency"`
	AutoRenew       bool                     `xqb:"auto_renew" json:"auto_renew"`
	BillingCycle    enums.BillingCycle       `xqb:"billing_cycle" json:"billing_cycle"`
	PlanLimits      PlanLimits               `xqb:"plan_limits" json:"plan_limits"`
	ExpireAction    enums.PlanExpireAction   `xqb:"expire_action" json:"expire_action"`
	DowngradeToPlan sql.NullInt64            `xqb:"downgrade_to_plan" json:"downgrade_to_plan"`
	GracePeriodDays int64                    `xqb:"grace_period_days" json:"grace_period_days"`
	CreatedAt       time.Time                `xqb:"created_at" json:"created_at"`
	UpdatedAt       time.Time                `xqb:"updated_at" json:"updated_at"`
	Plan            *Plan                    `table:"plans" json:"plan,omitempty"`
	Tenant          *Tenant                  ` table:"tenants" json:"tenant,omitempty"`
	BaseModel
}

func (Subscription) Table() string {
	return "subscriptions"
}

func (m Subscription) Cols() []any {
	return m.BaseModel.Cols(m, "subscription")
}

// helpers to return truncated dates
func (s *Subscription) GetStartDate() time.Time {
	y, m, d := s.StartDate.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (s *Subscription) GetEndDate() time.Time {
	y, m, d := s.EndDate.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func nowDate() time.Time {
	y, m, d := time.Now().UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// IsActive checks if the subscription is currently active
func (s *Subscription) IsActive() bool {
	now := nowDate()
	start := s.GetStartDate()
	end := s.GetEndDate()

	return (s.Status == enums.SubscriptionStatusActive || s.Status == enums.SubscriptionStatusTrial) &&
		!s.IsCanceled() &&
		(now.Equal(start) || now.After(start)) &&
		(now.Equal(end) || now.Before(end))
}

// IsExpired checks if subscription naturally expired (time ran out)
func (s *Subscription) IsExpired() bool {
	now := nowDate()
	end := s.GetEndDate()
	return now.After(end) || s.IsCanceled()
}

// IsCanceled checks if subscription was manually canceled
func (s *Subscription) IsCanceled() bool {
	return s.Status == enums.SubscriptionStatusCanceled
}

// IsEnded checks if subscription is no longer valid (expired or canceled)
func (s *Subscription) IsEnded() bool {
	return s.IsExpired() || s.IsCanceled()
}

// IsPending checks if subscription is pending (future start date)
func (s *Subscription) IsPending() bool {
	now := nowDate()
	start := s.GetStartDate()
	return now.Before(start) &&
		(s.Status == enums.SubscriptionStatusActive || s.Status == enums.SubscriptionStatusTrial)
}

// CanRenew checks if subscription is eligible for renewal
func (s *Subscription) CanRenew() bool {
	return !s.IsCanceled()
}

// RemainingDays returns number of days left until expiry
func (s *Subscription) RemainingDays() int {
	if s.IsEnded() {
		return 0
	}
	end := s.GetEndDate()
	now := nowDate()
	return int(math.Ceil(end.Sub(now).Hours() / 24))
}

// IsTrial checks if subscription is a trial
func (s *Subscription) IsTrial() bool {
	return s.Status == enums.SubscriptionStatusTrial
}

// RenewDate returns the expected next renewal date
func (s *Subscription) RenewDate() *time.Time {
	if s.IsCanceled() {
		return nil
	}
	end := s.GetEndDate()
	return &end
}

/*
* CalculateUpgradeCost calculates the cost of upgrading the subscription
 */
func (s *Subscription) CalculateUpgradeCost(newPrice float64) float64 {
	totalDays := s.EndDate.Sub(s.GetStartDate()).Hours() / 24
	daysPassed := time.Since(s.GetStartDate()).Hours() / 24

	if totalDays <= 0 {
		return newPrice
	}

	if daysPassed < 0 {
		daysPassed = 0
	} else if daysPassed > totalDays {
		daysPassed = totalDays
	}

	usedValue := (s.Price / totalDays) * daysPassed
	remainingCredit := s.Price - usedValue
	upgradeCost := newPrice - remainingCredit

	// Never negative
	if upgradeCost < 0 {
		upgradeCost = 0
	}

	return math.Ceil(upgradeCost)
}
