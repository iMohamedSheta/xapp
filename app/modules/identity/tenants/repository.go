package tenants

import (
	"context"
	"database/sql"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/models"
)

type TenantRepository struct{}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{}
}

func (a *TenantRepository) Paginate(c *gin.Context, filters *TenantManagerFilters) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &TenantManagerFilters{}
	}

	q := xqb.Table("tenants").
		WithContext(c).
		Select("tenants.*").
		GroupBy("tenants.id")

	if filters.Status != 0 {
		q.Where("tenants.status", "=", filters.Status)
	}

	mappedFields := map[string][]string{
		"name":       {"tenants.name"},
		"status":     {"tenants.status"},
		"created_at": {"tenants.created_at"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "created_at",
		SortOrder: "DESC",
	})

	filters.SetPaginationDefaults()

	return q.Paginate(filters.PerPage, filters.Page, "tenants.id")
}

// CreateTenant - Create new tenant record in database
func (a *TenantRepository) Create(c context.Context, tenant *models.Tenant, tx *sql.Tx) (*models.Tenant, error) {
	insertTime := time.Now()

	insertedId, err := xqb.Table("tenants").
		WithContext(c).
		WithTx(tx).
		InsertGetId([]map[string]any{
			{
				"name":       tenant.Name,
				"status":     tenant.Status,
				"created_at": insertTime,
				"updated_at": insertTime,
			},
		})
	if err != nil {
		return nil, err
	}
	return &models.Tenant{
		Id:        insertedId,
		Name:      tenant.Name,
		Status:    tenant.Status,
		CreatedAt: insertTime,
		UpdatedAt: insertTime,
	}, nil
}

// UpdateTenant - Update existing tenant record in database
func (a *TenantRepository) Update(c context.Context, tenantId int64, fields map[string]any, tx *sql.Tx) error {

	if fields["updated_at"] == nil {
		fields["updated_at"] = time.Now()
	}

	_, err := xqb.Table("tenants").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", tenantId).
		Update(fields)

	if err != nil {
		return err
	}

	return nil
}

// AddBalance - Atomically add balance to tenant
func (a *TenantRepository) AddBalance(c context.Context, tenantId int64, amount float64, tx *sql.Tx) error {
	_, err := xqb.Table("tenants").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", tenantId).
		Update(map[string]any{
			"balance":    xqb.Raw("balance + ?", amount),
			"updated_at": time.Now(),
		})

	if err != nil {
		return err
	}
	return nil
}

// ChargeFromBalance - Atomically charge from tenant balance
func (a *TenantRepository) ChargeFromBalance(c context.Context, tenantId int64, amount float64, tx *sql.Tx) error {
	rowsAffected, err := xqb.Table("tenants").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", tenantId).
		Where("balance", ">=", amount).
		Update(map[string]any{
			"balance":    xqb.Raw("balance - ?", amount),
			"updated_at": time.Now(),
		})
	if err != nil {
		return err
	}

	// Ensure that one row was updated (atomic success)
	if rowsAffected == 0 {
		return xerr.New("Insufficient balance to pay the subscription", enums.XErrValidationError, nil).
			WithDetails(map[string]any{
				"plan_id": fmt.Sprintf("الرصيد الخاص بك غير كافي لدفع الاشتراك برجاء شحن رصيد كافي الرصيد المطلوب هو  (%v)", amount),
			})
	}

	return nil
}

func (a *TenantRepository) FindById(c context.Context, tenantId int64) (*models.Tenant, error) {
	data, err := xqb.Table("tenants").
		WithContext(c).
		Where("id", "=", tenantId).
		First()
	if err != nil {
		return nil, err
	}

	var tenant models.Tenant
	if err := xqb.Bind(data, &tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) GetManagerMetrics(c *gin.Context, filters *requests.DashboardFilters) (map[string]any, error) {
	q := xqb.Table("tenants").WithContext(c)

	if filters != nil && filters.Period != "" && filters.Period != "all" {
		now := time.Now()
		var startDate time.Time

		switch filters.Period {
		case "today":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		case "week":
			startDate = now.AddDate(0, 0, -7)
		case "month":
			startDate = now.AddDate(0, -1, 0)
		}

		q.Where("created_at", ">=", startDate)
	}

	return q.Select(
		xqb.Raw(`
			COUNT(*) AS total_tenants,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') AS new_this_week,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') AS new_this_month
		`),
	).First()
}

func (r *TenantRepository) GetTopPerformers(c *gin.Context, limit int, filters *requests.DashboardFilters) ([]map[string]any, error) {
	q := xqb.Table("tenants AS a").
		WithContext(c).
		Select("a.id", "a.name").
		AddSelect(xqb.Raw("COUNT(DISTINCT e.id) AS total_events")).
		AddSelect(xqb.Raw("COUNT(DISTINCT o.id) AS total_orders")).
		AddSelect(xqb.Raw("COALESCE(SUM(o.total_price) FILTER (WHERE o.status = ?), 0) AS total_revenue", enums.OrderPaid)).
		LeftJoin("events AS e", "e.tenant_id = a.id").
		LeftJoin("orders AS o", "o.tenant_id = a.id")

	if filters != nil && filters.Period != "" && filters.Period != "all" {
		now := time.Now()
		var startDate time.Time

		switch filters.Period {
		case "today":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		case "week":
			startDate = now.AddDate(0, 0, -7)
		case "month":
			startDate = now.AddDate(0, -1, 0)
		}

		q.Where("o.created_at", ">=", startDate)
	}

	return q.GroupBy("a.id", "a.name").
		OrderBy("total_revenue", "DESC").
		Limit(limit).
		Get()
}

func (r *TenantRepository) GetRecent(c *gin.Context, limit int) ([]map[string]any, error) {
	return xqb.Table("tenants").
		WithContext(c).
		Select("id", "name", "created_at").
		OrderBy("created_at", "DESC").
		Limit(limit).
		Get()
}

func (r *TenantRepository) GetAgenciesWithoutEvents(c *gin.Context) (int64, error) {
	row, err := xqb.Table("tenants AS a").
		WithContext(c).
		LeftJoin("events AS e", "e.tenant_id = a.id").
		Select(xqb.Raw("COUNT(DISTINCT a.id) FILTER (WHERE e.id IS NULL) AS tenants_without_events")).
		First()

	if err != nil {
		return 0, err
	}

	if v, ok := row["tenants_without_events"].(int64); ok {
		return v, nil
	}

	return 0, nil
}
