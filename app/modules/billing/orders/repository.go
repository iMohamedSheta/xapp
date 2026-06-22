package orders

import (
	"context"
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/models"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (a *OrderRepository) Create(c context.Context, order *models.Order, tx *sql.Tx) (*models.Order, error) {
	insertTime := time.Now()

	insertedId, err := xqb.Table("orders").
		WithContext(c).
		WithTx(tx).
		InsertGetId([]map[string]any{
			{
				"tenant_id":      order.TenantId,
				"user_id":        order.UserId,
				"orderable_id":   order.OrderableId,
				"orderable_type": order.OrderableType,
				"quantity":       order.Quantity,
				"unit_price":     order.UnitPrice,
				"total_price":    order.TotalPrice,
				"currency":       order.Currency,
				"status":         order.Status,
				"created_at":     insertTime,
				"updated_at":     insertTime,
			},
		})
	if err != nil {
		return nil, err
	}

	order.Id = insertedId
	order.CreatedAt = insertTime
	order.UpdatedAt = insertTime

	return order, nil
}

func (r *OrderRepository) Paginate(c *gin.Context, tenantId *int64, filters *OrderFilters, types ...enums.OrderableType) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &OrderFilters{}
	}

	q := xqb.Table("orders").
		WithContext(c).
		Select("orders.*").
		AddSelect(models.User{}.Cols()...)

	if tenantId != nil {
		q.Where("orders.tenant_id", "=", tenantId)
	}

	if len(types) > 0 {
		q.WhereIn("orders.orderable_type", types)
	}

	q.Join("users", "users.id = orders.user_id")

	if filters.Status != 0 {
		q.Where("orders.status", "=", filters.Status)
	}

	if filters.Search != "" {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("users.name", "ILIKE", "%"+filters.Search+"%").
				OrWhere("users.email", "ILIKE", "%"+filters.Search+"%")
		})
	}

	mappedFields := map[string][]string{
		"quantity":    {"orders.quantity"},
		"unit_price":  {"orders.unit_price"},
		"total_price": {"orders.total_price"},
		"user_name":   {"users.name"},
		"user_email":  {"users.email"},
		"status":      {"orders.status"},
		"created_at":  {"orders.created_at"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "orders.created_at",
		SortOrder: "DESC",
	})

	filters.SetPaginationDefaults()

	return q.Paginate(filters.PerPage, filters.Page, "orders.id")
}

func (r *OrderRepository) Update(c context.Context, orderId int64, fields map[string]any, tx *sql.Tx) error {
	if fields["updated_at"] == nil {
		fields["updated_at"] = time.Now()
	}

	_, err := xqb.Table("orders").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", orderId).
		Update(fields)

	if err != nil {
		return err
	}

	return nil

}

func (r *OrderRepository) FindById(c context.Context, id int64, tenantId *int64) (*models.Order, error) {
	q := xqb.Table("orders").
		WithContext(c).
		Where("id", "=", id)

	if tenantId != nil {
		q.Where("tenant_id", "=", tenantId)
	}

	orderData, err := q.First()
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := xqb.Bind(orderData, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// GetDashboardMetrics - Dashboard metrics in one query
func (r *OrderRepository) GetDashboardMetrics(c *gin.Context, tenantId *int64, filters *requests.DashboardFilters) (map[string]any, error) {
	q := xqb.Table("orders").
		WithContext(c)

	if tenantId != nil {
		q.Where("tenant_id", "=", tenantId)
	}

	// Apply date filter if period is specified and not "all"
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

	// Get filtered metrics
	row, err := q.Select(
		xqb.Raw(`
			COUNT(*) AS total_orders,
			COALESCE(SUM(total_price) FILTER (WHERE status = ?), 0) AS total_revenue,
			COUNT(*) FILTER (WHERE status = ?) AS paid_orders
		`, enums.OrderPaid, enums.OrderPaid),
	).First()

	if err != nil {
		return nil, err
	}

	// extract values safely
	totalOrders := int64(0)
	if v, ok := row["total_orders"].(int64); ok {
		totalOrders = v
	}

	paidOrders := int64(0)
	if v, ok := row["paid_orders"].(int64); ok {
		paidOrders = v
	}

	// build metrics map
	metrics := map[string]any{
		"total_orders":    row["total_orders"],
		"total_revenue":   row["total_revenue"],
		"conversion_rate": 0.0,
	}

	if totalOrders > 0 {
		metrics["conversion_rate"] = float64(paidOrders) / float64(totalOrders) * 100
	}

	return metrics, nil
}

func (r *OrderRepository) GetRecentOrders(c *gin.Context, tenantId int64, limit int, filters *requests.DashboardFilters) ([]map[string]any, error) {
	q := xqb.Table("orders").
		WithContext(c).
		Select("orders.*").
		AddSelect("users.name as user_name").
		AddSelect("users.email as user_email").
		AddSelect("events.id as event_id").
		AddSelect("events.title as event_title").
		AddSelect("events.image as event_image").
		AddSelect("events.event_date as event_date").
		AddSelect("events.location as event_location").
		// Where("orders.agency_id", "=", tenantId).
		Where("orders.orderable_type", "=", "ticket").
		LeftJoin("tickets", "tickets.id = orders.orderable_id").
		LeftJoin("events", "events.id = tickets.event_id").
		LeftJoin("users", "users.id = orders.user_id").
		OrderBy("orders.created_at", "DESC").
		Limit(limit)

	// Apply date filter if period is specified and not "all"
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

		q.Where("orders.created_at", ">=", startDate)
	}

	return q.Get()
}

// --- Change order status ---
func (a *OrderRepository) ChangeStatus(c *gin.Context, orderId int64, status enums.OrderStatus, tx *sql.Tx) error {
	_, err := xqb.Table("orders").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", orderId).
		Update(map[string]any{"status": status})
	return err
}

func (r *OrderRepository) GetManagerMetrics(c *gin.Context, filters *requests.DashboardFilters) (map[string]any, error) {
	q := xqb.Table("orders").WithContext(c)

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

	row, err := q.Select(
		xqb.Raw(`
			COUNT(*) AS total_orders,
			COALESCE(SUM(total_price), 0) AS total_revenue,
			COALESCE(SUM(total_price) FILTER (WHERE status = ?), 0) AS paid_revenue,
			COUNT(*) FILTER (WHERE status = ?) AS paid_orders,
			COUNT(*) FILTER (WHERE status = ?) AS pending_orders,
			COUNT(*) FILTER (WHERE status = ?) AS cancelled_orders,
			COALESCE(AVG(total_price) FILTER (WHERE status = ?), 0) AS avg_order_value
		`, enums.OrderPaid, enums.OrderPaid, enums.OrderPending, enums.OrderCancelled, enums.OrderPaid),
	).First()

	if err != nil {
		return nil, err
	}

	totalOrders := int64(0)
	if v, ok := row["total_orders"].(int64); ok {
		totalOrders = v
	}

	paidOrders := int64(0)
	if v, ok := row["paid_orders"].(int64); ok {
		paidOrders = v
	}

	row["conversion_rate"] = 0.0
	if totalOrders > 0 {
		row["conversion_rate"] = float64(paidOrders) / float64(totalOrders) * 100
	}

	return row, nil
}

func (r *OrderRepository) GetRevenueByAgency(c *gin.Context, filters *requests.DashboardFilters) ([]map[string]any, error) {
	q := xqb.Table("orders AS o").
		WithContext(c).
		Select("a.id", "a.name").
		AddSelect(xqb.Raw("COUNT(*) AS total_orders")).
		AddSelect(xqb.Raw("COALESCE(SUM(o.total_price) FILTER (WHERE o.status = ?), 0) AS total_revenue", enums.OrderPaid)).
		Join("agencies AS a", "a.id = o.agency_id").
		Where("o.status", "=", enums.OrderPaid)

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
		Get()
}

func (r *OrderRepository) GetRevenueTrends(c *gin.Context, filters *requests.DashboardFilters) ([]map[string]any, error) {
	period := "month"
	if filters != nil && filters.Period != "" && filters.Period != "all" {
		period = filters.Period
	}

	var dateFormat string

	switch period {
	case "today":
		dateFormat = "HH24:00"
	case "week":
		dateFormat = "YYYY-MM-DD"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	now := time.Now()
	var startDate time.Time

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startDate = now.AddDate(0, 0, -7)
	default:
		startDate = now.AddDate(0, -1, 0)
	}

	return xqb.Table("orders").
		WithContext(c).
		Select(xqb.Raw("TO_CHAR(created_at, ?) AS date", dateFormat)).
		AddSelect(xqb.Raw("COUNT(*) AS total_orders")).
		AddSelect(xqb.Raw("COALESCE(SUM(total_price) FILTER (WHERE status = ?), 0) AS revenue", enums.OrderPaid)).
		Where("created_at", ">=", startDate).
		GroupBy("date").
		OrderBy("date", "ASC").
		Get()
}

// GetOrderDetails - Get complete order details with all relations
func (r *OrderRepository) GetOrderDetails(c context.Context, orderId int64, tenantId int64) (map[string]any, error) {
	q := xqb.Table("orders").
		WithContext(c).
		Select("orders.*").
		// User columns
		AddSelect(models.User{}.Cols()...).
		// Invoice columns
		AddSelect(models.Invoice{}.Cols()...).
		// Transaction columns
		AddSelect(models.Transaction{}.Cols()...).
		// Agency columns
		AddSelect(models.Tenant{}.Cols()...).
		Where("orders.id", "=", orderId)

	// Apply agency and filters
	// if tenantId != 0 {
	// 	q.Where("orders.agency_id", "=", tenantId)
	// }

	// Join related tables
	q.LeftJoin("users", "users.id = orders.user_id").
		LeftJoin("invoices", "invoices.order_id = orders.id").
		LeftJoin("transactions", "transactions.order_id = orders.id").
		LeftJoin("tenants", "tenants.id = orders.tenant_id")

	// Execute query
	result, err := q.First()
	if err != nil {
		return nil, err
	}

	return result, nil
}
