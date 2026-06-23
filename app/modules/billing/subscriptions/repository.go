package subscriptions

import (
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

type SubscriptionRepository struct {
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{}
}

func (r *SubscriptionRepository) First(c *gin.Context) (*models.Subscription, error) {
	return xqb.Model[models.Subscription]().
		WithContext(c).
		First()
}

func (r *SubscriptionRepository) FirstActiveSubscription(c *gin.Context) (*models.Subscription, error) {
	return xqb.Model[models.Subscription]().
		WithContext(c).
		Where("status", "=", enums.SubscriptionStatusActive).
		OrderBy("id", "desc").
		First()
}

func (t *SubscriptionRepository) Create(c *gin.Context, subscription *models.Subscription, tx *sql.Tx) (*models.Subscription, error) {
	if subscription == nil {
		return nil, xerr.New("cannot create subscription with nil subscription", enums.XErrBadRequestError, nil)
	}

	insertTime := time.Now()

	values := map[string]any{
		"tenant_id":      subscription.TenantId,
		"plan_id":        subscription.PlanId,
		"start_date":     subscription.StartDate,
		"end_date":       subscription.EndDate,
		"status":         subscription.Status,
		"price":          subscription.Price,
		"original_price": subscription.OriginalPrice,
		"currency":       subscription.Currency,
		"billing_cycle":  subscription.BillingCycle,
		"auto_renew":     subscription.AutoRenew,

		// Snapshot settings
		"plan_limits": subscription.PlanLimits,

		"expire_action":     subscription.ExpireAction,
		"grace_period_days": subscription.GracePeriodDays,
		"downgrade_to_plan": subscription.DowngradeToPlan,

		"created_at": insertTime,
		"updated_at": insertTime,
	}

	id, err := xqb.Table("subscriptions").
		WithContext(c).
		WithTx(tx).
		InsertGetId([]map[string]any{values})

	if err != nil {
		return nil, err
	}

	subscription.Id = id
	subscription.CreatedAt = insertTime
	subscription.UpdatedAt = insertTime

	return subscription, nil
}
