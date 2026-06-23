package invoices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/imohamedsheta/xapp/app/models"
)

type MonthlyGrowth struct {
	Total         float64              `json:"total"`
	GrowthPercent float64              `json:"growth_percent"`
	Trend         enums.StatisticTrend `json:"trend"` // "growing", "lowering", or "stable"
}

type InvoiceRepository struct{}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{}
}

func (r *InvoiceRepository) Paginate(c *gin.Context, tenantId *int64, filters *InvoiceFilters) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &InvoiceFilters{}
	}

	q := xqb.Table("invoices").
		WithContext(c).
		Select("invoices.*").
		AddSelect("agencies.name as tenant_name").
		Join("agencies", "agencies.id = invoices.tenant_id").
		LeftJoin("orders", "orders.id = invoices.order_id").
		LeftJoin("transactions", "transactions.id = invoices.transaction_id")

	q.LeftJoin("radius_users", "radius_users.id = invoices.user_id AND invoices.user_type = 'radius_user'").
		LeftJoin("radius_cards", "radius_cards.id = invoices.user_id AND invoices.user_type = 'radius_card'")

	q.AddSelect(`
			COALESCE(radius_users.fullname, radius_cards.username) AS customer_fullname
	`)

	if tenantId != nil {
		q.Where("invoices.tenant_id", "=", *tenantId)
	}

	if filters.Status != "" {
		q.Where("invoices.status", "=", filters.Status)
	}

	if filters.UserId != nil {
		q.Where("invoices.user_id", "=", *filters.UserId)
	}

	if filters.Type != "" {
		q.Where("invoices.type", "=", filters.Type)
	}

	if filters.UserType != "" {
		q.Where("invoices.user_type", "=", filters.UserType)
	}

	if filters.FromDate != "" && filters.ToDate != "" {
		from, err1 := time.Parse("2006-01-02", filters.FromDate)
		to, err2 := time.Parse("2006-01-02", filters.ToDate)
		if err1 == nil && err2 == nil {
			to = to.AddDate(0, 0, 1)
			q.WhereBetween("invoices.created_at", from, to)
		}
	} else if filters.FromDate != "" {
		if from, err := time.Parse("2006-01-02", filters.FromDate); err == nil {
			q.Where("invoices.created_at", ">=", from)
		}
	} else if filters.ToDate != "" {
		if to, err := time.Parse("2006-01-02", filters.ToDate); err == nil {
			to = to.AddDate(0, 0, 1)
			q.Where("invoices.created_at", "<=", to)
		}
	}

	if filters.Search != "" {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("invoices.invoice_number", "ILIKE", "%"+filters.Search+"%").
				OrWhere("COALESCE(radius_users.fullname, radius_cards.username)", "ILIKE", "%"+filters.Search+"%")
		})
	}

	mappedFields := map[string][]string{
		"status":         {"invoices.status"},
		"type":           {"invoices.type"},
		"invoice_number": {"invoices.invoice_number"},
		"amount":         {"invoices.amount"},
		"paid":           {"invoices.paid"},
		"paid_at":        {"invoices.paid_at"},
		"created_at":     {"invoices.created_at"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "invoices.created_at",
		SortOrder: "DESC",
	})

	filters.SetPaginationDefaults()

	return q.Paginate(filters.PerPage, filters.Page, "invoices.id")
}

func (r *InvoiceRepository) Create(c context.Context, invoice *models.Invoice, tx *sql.Tx) (*models.Invoice, error) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// get last invoice number for this tenant in this month
	var lastNumber int64
	value, err := xqb.Table("invoices").
		WithContext(c).
		WithTx(tx).
		Where("tenant_id", "=", invoice.TenantId).
		Where("EXTRACT(MONTH FROM created_at)", "=", month).
		OrderBy("id", "DESC").
		Value("invoice_seq")

	if err != nil {
		if !errors.Is(err, xqb.ErrNotFound) {
			return nil, err
		}
		lastNumber = 0
	} else {
		switch v := value.(type) {
		case int64:
			lastNumber = v
		case int32:
			lastNumber = int64(v)
		case int:
			lastNumber = int64(v)
		default:
			return nil, xerr.New("invalid invoice_seq type", enums.XErrServerError, nil)
		}
	}

	nextNumber := lastNumber + 1
	invoiceNumber := fmt.Sprintf("INV-%d-%02d-%06d-A%d", year, month, nextNumber, invoice.TenantId)

	id, err := xqb.Table("invoices").
		WithContext(c).
		WithTx(tx).
		InsertGetId([]map[string]any{
			{
				"tenant_id":        invoice.TenantId,
				"creator_id":       invoice.CreatorId,
				"creator_type":     invoice.CreatorType,
				"user_id":          invoice.UserId,
				"user_type":        invoice.UserType,
				"order_id":         invoice.OrderId,
				"transaction_id":   invoice.TransactionId,
				"invoiceable_type": invoice.InvoiceableType,
				"invoiceable_id":   invoice.InvoiceableId,
				"invoice_seq":      nextNumber,
				"invoice_number":   invoiceNumber,
				"type":             invoice.Type,
				"status":           invoice.Status,
				"amount":           invoice.Amount,
				"paid":             invoice.Paid,
				"currency":         invoice.Currency,
				"due_date":         invoice.DueDate,
				"paid_at":          invoice.PaidAt,
				"notes":            invoice.Notes,
				"metadata":         invoice.Metadata,
				"created_at":       now,
				"updated_at":       now,
			},
		})
	if err != nil {
		return nil, err
	}

	invoice.Id = id
	invoice.InvoiceNumber = invoiceNumber
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	return invoice, nil
}

func (r *InvoiceRepository) LatestTenantSubscriptionInvoices(c *gin.Context, tenantId *int64, n int) ([]map[string]any, error) {
	if n == 0 {
		return nil, nil
	}

	qb := xqb.Table("invoices").
		WithContext(c).
		Where("invoiceable_type", "=", enums.InvoiceableTypeSubscription)

	if tenantId != nil {
		qb = qb.Where("tenant_id", "=", *tenantId)
	}

	invoices, err := qb.OrderBy("id", "DESC").
		Take(n).
		Get()
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *InvoiceRepository) TotalPaidThisMonth(c *gin.Context, tenantId *int64) (float64, error) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	qb := xqb.Table("invoices").
		WithContext(c).
		Where("EXTRACT(YEAR FROM created_at)", "=", year).
		Where("EXTRACT(MONTH FROM created_at)", "=", month).
		Where("status", "=", enums.InvoiceStatusPaid)

	if tenantId != nil {
		qb.Where("tenant_id", "=", tenantId)
	}

	return qb.Sum("paid")
}

func (r *InvoiceRepository) MonthlyTenantPaymentGrowth(c *gin.Context, tenantId *int64) (*MonthlyGrowth, error) {
	now := time.Now()
	currentYear, currentMonth := now.Year(), int(now.Month())

	qb := xqb.Table("invoices").
		WithContext(c).
		Where("status", "=", enums.InvoiceStatusPaid)

	if tenantId != nil {
		qb.Where("tenant_id", "=", *tenantId)
	}

	// Aggregate: sum per month
	rows, err := qb.
		Select(
			"EXTRACT(YEAR FROM created_at) AS year",
			"EXTRACT(MONTH FROM created_at) AS month",
			"SUM(paid) AS total",
		).
		GroupBy("year", "month").
		Get()
	if err != nil {
		return nil, err
	}

	var totalCurrent float64
	var prevTotal float64
	var prevMonthsCount float64

	for _, row := range rows {
		yStr := fmt.Sprintf("%v", row["year"])
		mStr := fmt.Sprintf("%v", row["month"])
		sumStr := fmt.Sprintf("%v", row["total"])

		y, _ := strconv.Atoi(yStr)
		m, _ := strconv.Atoi(mStr)
		sum, _ := strconv.ParseFloat(sumStr, 64)

		if y == currentYear && m == currentMonth {
			totalCurrent = sum
		} else {
			prevTotal += sum
			prevMonthsCount++
		}
	}

	var growth float64
	trend := enums.StatisticTrendStable

	var prevAvg float64
	if prevMonthsCount > 0 {
		prevAvg = prevTotal / prevMonthsCount
	}

	if prevAvg == 0 {
		if totalCurrent > 0 {
			growth = 100
			trend = enums.StatisticTrendGrowing
		}
	} else {
		growth = ((totalCurrent - prevAvg) / prevAvg) * 100
		if growth > 0 {
			trend = enums.StatisticTrendGrowing
		} else if growth < 0 {
			trend = enums.StatisticTrendLowering
		}
	}

	return &MonthlyGrowth{
		Total:         totalCurrent,
		GrowthPercent: growth,
		Trend:         trend,
	}, nil
}

func (r *InvoiceRepository) Update(c *gin.Context, tenantId *int64, invoice *models.Invoice, tx *sql.Tx) error {
	now := time.Now()
	q := xqb.Table("invoices").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", invoice.Id)

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	_, err := q.Update(map[string]any{
		"status":     invoice.Status,
		"paid":       invoice.Paid,
		"currency":   invoice.Currency,
		"due_date":   invoice.DueDate,
		"paid_at":    invoice.PaidAt,
		"metadata":   invoice.Metadata,
		"notes":      invoice.Notes,
		"updated_at": now,
	})
	return err
}

func (r *InvoiceRepository) AdminInvoiceStats(c *gin.Context, tenantId *int64, filters *InvoiceFilters) (map[string]any, error) {
	if filters == nil {
		filters = &InvoiceFilters{}
	}

	q := xqb.Table("invoices").WithContext(c)

	q.LeftJoin("radius_users", "radius_users.id = invoices.user_id AND invoices.user_type = 'radius_user'").
		LeftJoin("radius_cards", "radius_cards.id = invoices.user_id AND invoices.user_type = 'radius_card'").
		LeftJoin("users", "users.id = invoices.user_id AND (invoices.user_type = 'distributor' OR invoices.user_type = 'user')")

	q.AddSelect(`
			COALESCE(radius_users.fullname, radius_cards.username, users.name) AS customer_fullname
	`)

	// if tenantId != nil {
	// 	q.Where("invoices.tenant_id", "=", *tenantId)
	// }

	if filters.UserId != nil {
		q.Where("invoices.user_id", "=", *filters.UserId)
	}

	if filters.Type != "" {
		q.Where("invoices.type", "=", filters.Type)
	}

	if filters.Status != "" {
		q.Where("invoices.status", "=", filters.Status)
	}

	if filters.UserType != "" {
		q.Where("invoices.user_type", "=", filters.UserType)
	}

	if filters.Search != "" {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("invoices.invoice_number", "ILIKE", "%"+filters.Search+"%").
				OrWhere("COALESCE(radius_users.fullname, radius_cards.username)", "ILIKE", "%"+filters.Search+"%")
		})
	}

	if filters.FromDate != "" && filters.ToDate != "" {
		from, err1 := time.Parse("2006-01-02", filters.FromDate)
		to, err2 := time.Parse("2006-01-02", filters.ToDate)
		if err1 == nil && err2 == nil {
			q.WhereBetween("invoices.created_at", from, to)
		}
	} else if filters.FromDate != "" {
		if from, err := time.Parse("2006-01-02", filters.FromDate); err == nil {
			q.Where("invoices.created_at", ">=", from)
		}
	} else if filters.ToDate != "" {
		if to, err := time.Parse("2006-01-02", filters.ToDate); err == nil {
			q.Where("invoices.created_at", "<=", to)
		}
	}

	pgRaw := xqb.Raw(`
    SUM(amount) AS total_amount,
    SUM(amount - paid) AS total_unpaid,
    SUM(paid) AS total_paid,
    COUNT(*) FILTER (WHERE invoices.status = 'pending' OR invoices.status = 'partially_paid') AS unpaid_count,
    SUM(paid) FILTER (WHERE DATE(paid_at) = CURRENT_DATE) AS paid_today_amount,
    COUNT(*) FILTER (WHERE DATE(paid_at) = CURRENT_DATE AND (invoices.status = 'paid' OR invoices.status = 'partially_paid')) AS paid_today_count,
    SUM(amount - paid) FILTER (WHERE DATE(invoices.created_at) = CURRENT_DATE AND (invoices.status = 'pending' OR invoices.status = 'partially_paid')) AS unpaid_today_amount
`)

	mysqlRaw := xqb.Raw(`
    SUM(amount) AS total_amount,
    SUM(amount - paid) AS total_unpaid,
    SUM(paid) AS total_paid,
    COUNT(CASE WHEN invoices.status = 'pending' OR invoices.status = 'partially_paid' THEN 1 END) AS unpaid_count,
    SUM(CASE WHEN DATE(paid_at) = CURRENT_DATE THEN paid ELSE 0 END) AS paid_today_amount,
    COUNT(CASE WHEN DATE(paid_at) = CURRENT_DATE AND (invoices.status = 'paid' OR invoices.status = 'partially_paid') THEN 1 END) AS paid_today_count,
    SUM(CASE WHEN DATE(invoices.created_at) = CURRENT_DATE AND (invoices.status = 'pending' OR invoices.status = 'partially_paid') THEN amount - paid ELSE 0 END) AS unpaid_today_amount
`)

	raw := xqb.RawDialect(string(xqb.DialectPostgres), map[string]*types.Expression{
		string(xqb.DialectPostgres): pgRaw,
		string(xqb.DialectMySql):    mysqlRaw,
	})

	q.Select(raw)

	mappedFields := map[string][]string{
		"invoice_number": {"invoices.invoice_number"},
		"type":           {"invoices.type"},
		"status":         {"invoices.status"},
		"amount":         {"invoices.amount"},
		"paid":           {"invoices.paid"},
		"paid_at":        {"invoices.paid_at"},
		"created_at":     {"invoices.created_at"},
	}

	if filters.SearchBy != nil {
		for key, value := range filters.SearchBy {
			if fieldsToSearch, ok := mappedFields[key]; ok {
				q.WhereGroup(func(qb *xqb.QueryBuilder) {
					for i, field := range fieldsToSearch {
						cond := "CAST(" + field + " AS TEXT)"
						if i == 0 {
							qb.Where(cond, "LIKE", "%"+value+"%")
						} else {
							qb.OrWhere(cond, "LIKE", "%"+value+"%")
						}
					}
				})
			}
		}
	}

	return q.First()
}

func (r *InvoiceRepository) FindById(c *gin.Context, tenantId *int64, id int64, tx *sql.Tx) (*models.Invoice, error) {
	q := xqb.Model[models.Invoice]().
		WithContext(c).
		WithTx(tx).
		Where("id", "=", id)

	if tenantId != nil {
		q.Where("tenant_id", "=", tenantId)
	}

	return q.First()
}

func (r *InvoiceRepository) FindWithTenantNameAndCustomerNameById(c *gin.Context, tenantId *int64, id int64, tx *sql.Tx) (map[string]any, error) {
	q := xqb.Table("invoices").
		WithContext(c).
		WithTx(tx).
		Select("invoices.*").
		Where("invoices.id", "=", id)

	q.LeftJoin("radius_users", "radius_users.id = invoices.user_id AND invoices.user_type = 'radius_user'").
		LeftJoin("radius_cards", "radius_cards.id = invoices.user_id AND invoices.user_type = 'radius_card'")

	q.AddSelect(`
			COALESCE(radius_users.fullname, radius_cards.username) AS customer_fullname
	`)

	q.LeftJoin("agencies", "agencies.id = invoices.tenant_id").
		AddSelect(`
			agencies.name AS tenant_name
	`)

	// if tenantId != nil {
	// 	q.Where("invoices.tenant_id", "=", tenantId)
	// }

	return q.First()
}

func (r *InvoiceRepository) GetRevenueStats(c context.Context, tenantId *int64, period string) (map[string]any, error) {
	q := xqb.Table("invoices").
		WithContext(c).
		Select(
			xqb.Sum("paid", "total_revenue"),
			xqb.Sum("amount - paid", "outstanding_amount"),
			xqb.Count("*", "total_invoices"),
		).
		AddSelect(xqb.Raw("COUNT(*) FILTER (WHERE status = 'pending' OR status = 'partially_paid') AS outstanding_count"))

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	switch period {
	case "today":
		q.Where(xqb.Date("invoices.created_at", ""), "=", time.Now().Format("2006-01-02"))
	case "week":
		q.Where("invoices.created_at", ">=", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	case "month":
		q.Where("invoices.created_at", ">=", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	case "year":
		q.Where("invoices.created_at", ">=", time.Now().AddDate(-1, 0, 0).Format("2006-01-02"))
	}

	return q.First()
}

func (r *InvoiceRepository) GetRevenueByOffer(c context.Context, tenantId *int64, period string, limit int) ([]map[string]any, error) {
	q := xqb.Table("invoices").
		WithContext(c).
		Join("radius_users", "radius_users.id = invoices.user_id AND invoices.user_type = 'radius_user'").
		Join("offers", "offers.id = radius_users.offer_id").
		Select("offers.name as offer_name", "offers.id as offer_id").
		AddSelect(xqb.Sum("invoices.paid", "total_revenue")).
		AddSelect(xqb.Count("invoices.id", "invoice_count")).
		GroupBy("offers.id", "offers.name").
		OrderBy("total_revenue", "DESC").
		Take(limit)

	if tenantId != nil {
		q.Where("invoices.tenant_id", "=", *tenantId)
	}

	switch period {
	case "today":
		q.Where(xqb.Date("invoices.created_at", ""), "=", time.Now().Format("2006-01-02"))
	case "week":
		q.Where("invoices.created_at", ">=", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	case "month":
		q.Where("invoices.created_at", ">=", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	}

	return q.Get()
}
