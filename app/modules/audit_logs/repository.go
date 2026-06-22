package audit_logs

import (
	"context"
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/models"
)

type AuditLogRepository struct {
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

func (r *AuditLogRepository) Create(c context.Context, log *models.AuditLog, tx *sql.Tx) error {

	now := time.Now()
	data := map[string]any{
		"tenant_id":      log.TenantId,
		"user_id":        log.UserId,
		"user_type":      log.UserType,
		"auditable_id":   log.AuditableId,
		"auditable_type": log.AuditableType,
		"action":         log.Action,
		"summary":        log.Summary,
		"details":        log.Details,
		"created_at":     now,
		"updated_at":     now,
	}

	return xqb.Table("audit_logs").
		WithContext(c).
		WithTx(tx).
		Insert([]map[string]any{data})
}

func (r *AuditLogRepository) BulkCreate(c *gin.Context, logs []*models.AuditLog, tx *sql.Tx) error {

	var data []map[string]any
	for _, g := range logs {
		now := time.Now()
		data = append(data, map[string]any{
			"tenant_id":      g.TenantId,
			"user_id":        g.UserId,
			"user_type":      g.UserType,
			"auditable_id":   g.AuditableId,
			"auditable_type": g.AuditableType,
			"action":         g.Action,
			"summary":        g.Summary,
			"details":        g.Details,
			"created_at":     now,
			"updated_at":     now,
		})
	}

	return xqb.Table("audit_logs").
		WithContext(c).
		WithTx(tx).
		Insert(data)
}

func (r *AuditLogRepository) Paginate(c *gin.Context, tenantId *int64, filters *AuditLogFilters) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &AuditLogFilters{}
	}

	q := xqb.Table("audit_logs").
		WithContext(c).
		Select("audit_logs.*").
		AddSelect("users.name as actor_name").
		LeftJoin("users", "users.id = audit_logs.user_id")

	// Join target tables based on auditable_type
	if filters.AuditableType != "" {
		if filters.AuditableType == enums.AuditableTypeUser {
			q.LeftJoin("users as target_table",
				"target_table.id = audit_logs.auditable_id AND audit_logs.auditable_type = ?",
				enums.AuditableTypeUser).
				AddSelect("target_table.name as target_name")
		} else {
			q.AddSelect("NULL as target_name")
		}
	} else {
		// When no specific filter, join all candidate tables with type check
		q.LeftJoin("users as target_user",
			"target_user.id = audit_logs.auditable_id AND audit_logs.auditable_type = ?",
			enums.AuditableTypeUser)

		q.LeftJoin("nas as target_nas",
			"target_nas.id = audit_logs.auditable_id AND audit_logs.auditable_type = ?",
			enums.AuditableTypeNas)

		// Use CASE instead of COALESCE to avoid cross-database error
		q.AddSelectRaw(`CASE 
			WHEN target_user.name IS NOT NULL THEN target_user.name
			WHEN target_nas.name IS NOT NULL THEN target_nas.name
			ELSE NULL
		END AS target_name`)
	}

	// Apply filters
	if tenantId != nil {
		q.Where("audit_logs.tenant_id", "=", *tenantId)
	}

	if filters.Action != "" {
		q.Where("audit_logs.action", "=", filters.Action)
	}

	if filters.AuditableType != "" {
		q.Where("audit_logs.auditable_type", "=", filters.AuditableType)
	}

	if filters.AuditableId != 0 {
		q.Where("audit_logs.auditable_id", "=", filters.AuditableId)
	}

	if filters.UserId != 0 {
		q.Where("audit_logs.user_id", "=", filters.UserId)
	}

	if filters.Group == "subscription" {
		q.Where("audit_logs.auditable_type", "=", enums.AuditableTypeUser)
		q.WhereIn("audit_logs.action", []any{
			enums.AuditLogActionSubscribe,
			enums.AuditLogActionRenew,
		})
	}

	if filters.IsClientView {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.WhereGroup(func(sqb *xqb.QueryBuilder) {
				sqb.Where("audit_logs.auditable_type", "=", enums.AuditableTypeUser).
					Where("audit_logs.auditable_id", "=", filters.OwnRadiusUserId)
			}).OrWhere("audit_logs.user_id", "=", filters.OwnUserId)
		})
	}

	// Date range filtering
	if filters.FromDate != "" && filters.ToDate != "" {
		from, err1 := time.Parse("2006-01-02", filters.FromDate)
		to, err2 := time.Parse("2006-01-02", filters.ToDate)
		if err1 == nil && err2 == nil {
			// Add one day to 'to' date to include the entire end date
			to = to.AddDate(0, 0, 1)
			q = q.WhereBetween("audit_logs.created_at", from, to)
		}
	} else if filters.FromDate != "" {
		if from, err := time.Parse("2006-01-02", filters.FromDate); err == nil {
			q = q.Where("audit_logs.created_at", ">=", from)
		}
	} else if filters.ToDate != "" {
		if to, err := time.Parse("2006-01-02", filters.ToDate); err == nil {
			// Add one day to include the entire end date
			to = to.AddDate(0, 0, 1)
			q = q.Where("audit_logs.created_at", "<", to)
		}
	}

	// Search filtering
	if filters.Search != "" {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("users.name", "ILIKE", "%"+filters.Search+"%").
				OrWhere("audit_logs.summary", "ILIKE", "%"+filters.Search+"%").
				OrWhere("CAST(audit_logs.details AS TEXT)", "ILIKE", "%"+filters.Search+"%")
		})
	}

	// Field mapping for sorting
	mappedFields := map[string][]string{
		"action":         {"audit_logs.action"},
		"auditable_type": {"audit_logs.auditable_type"},
		"created_at":     {"audit_logs.created_at"},
		"actor_name":     {"users.name"},
		"summary":        {"audit_logs.summary"},
	}

	// Add target_name to mapped fields dynamically
	if filters.AuditableType == enums.AuditableTypeUser {
		mappedFields["target_name"] = []string{"target_table.name"}
	} else if filters.AuditableType == enums.AuditableTypeNas {
		mappedFields["target_name"] = []string{"target_table.name"}
	} else {
		// For non-filtered queries, we can't reliably sort by target_name due to COALESCE
		// But we'll provide the fields anyway
		mappedFields["target_name"] = []string{"target_user.name"}
	}

	// Apply sorting and filtering
	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "created_at",
		SortOrder: "DESC",
	})

	// Set pagination defaults
	filters.SetPaginationDefaults()

	// Execute pagination query
	return q.Paginate(filters.PerPage, filters.Page, "audit_logs.id")
}
func (r *AuditLogRepository) AuditLogStats(c *gin.Context, tenantId *int64, filters *AuditLogFilters) (map[string]any, error) {
	if filters == nil {
		filters = &AuditLogFilters{}
	}

	q := xqb.Table("audit_logs").WithContext(c).
		LeftJoin("users", "users.id = audit_logs.user_id")

	if tenantId != nil {
		q.Where("audit_logs.tenant_id", "=", *tenantId)
	}

	// Apply date filters to stats as well? Usually stats are "This Month" or based on filter.
	// The request says "this month" explicitly for one stat, but usually stats follow the page filter or have their own defaults.
	// Assuming stats follow the global time filter if provided, otherwise might need defaults.
	// For now, let's respect the filters passed in (which come from UI).

	if filters.FromDate != "" && filters.ToDate != "" {
		from, err1 := time.Parse("2006-01-02", filters.FromDate)
		to, err2 := time.Parse("2006-01-02", filters.ToDate)
		if err1 == nil && err2 == nil {
			q = q.WhereBetween("audit_logs.created_at", from, to)
		}
	} else if filters.FromDate != "" {
		if from, err := time.Parse("2006-01-02", filters.FromDate); err == nil {
			q = q.Where("audit_logs.created_at", ">=", from)
		}
	} else if filters.ToDate != "" {
		if to, err := time.Parse("2006-01-02", filters.ToDate); err == nil {
			q = q.Where("audit_logs.created_at", "<=", to)
		}
	}

	// Business Stats
	// 1. Total Renewals: action = 'renew'
	// 2. New Users: action = 'create' AND auditable_type = 'radius_user'
	// 3. Upgrades: action = 'update' AND auditable_type = 'radius_user' (Assuming offer changes) OR 'subscription'
	// 4. Subscriptions (Agency?): action = 'create' AND auditable_type = 'subscription'

	q.Select(xqb.Raw(`
		COUNT(*) AS total_count,
		COUNT(*) FILTER (WHERE auditable_type = ? AND action = ?) AS renewals_count,
		COUNT(*) FILTER (WHERE auditable_type = ? AND action IN (?, ?)) AS new_users_count,
		COUNT(*) FILTER (WHERE auditable_type = ? AND action = ?) AS upgrades_count,
		COUNT(*) FILTER (WHERE auditable_type = ? AND action = ?) AS subscriptions_count
	`,
		enums.AuditableTypeUser, enums.AuditLogActionRenew,
		enums.AuditableTypeUser, enums.AuditLogActionCreate, enums.AuditLogActionSubscribe,
		enums.AuditableTypeUser, enums.AuditLogActionUpdate,
		enums.AuditableTypeUser, enums.AuditLogActionSubscribe,
	))

	return q.First()
}

func (r *AuditLogRepository) GetByAuditable(c *gin.Context, agencyId *int64, auditableId int64, auditableType enums.AuditableType, limit int) ([]map[string]any, error) {
	q := xqb.Table("audit_logs").
		WithContext(c).
		Select("audit_logs.id", "audit_logs.action", "audit_logs.summary as description", "audit_logs.created_at").
		Where("audit_logs.auditable_id", "=", auditableId).
		Where("audit_logs.auditable_type", "=", auditableType).
		OrderBy("audit_logs.created_at", "DESC")

	// if agencyId != nil {
	// 	q.Where("audit_logs.agency_id", "=", *agencyId)
	// }

	if limit > 0 {
		q.Limit(limit)
	}

	return q.Get()
}
