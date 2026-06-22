package plans

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/models"
)

type PlanRepository struct {
}

func NewPlanRepository() *PlanRepository {
	return &PlanRepository{}
}

func (r *PlanRepository) Paginate(c *gin.Context, perPage, page int, filters *PlanFilters) ([]map[string]any, map[string]any, error) {
	q := xqb.Table("plans").
		WithContext(c).
		Select("plans.*").
		AddSelect(models.PlanSetting{}.Cols()...).
		LeftJoin("plan_settings", "plan_settings.plan_id = plans.id")

	if filters == nil {
		filters = &PlanFilters{}
	}

	if isActive := filters.IsActive; isActive != nil {
		q.Where("plans.is_active", "=", isActive)
	}

	if filters.Search != "" {
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("plans.name", "ILIKE", "%"+filters.Search+"%").
				OrWhere("plans.features", "ILIKE", "%"+filters.Search+"%")
		})
	}

	mappedFields := map[string][]string{
		"name":       {"plans.name"},
		"features":   {"plans.features"},
		"limits":     {"plan_settings.devices_limit", "plan_settings.messages_limit"},
		"is_active":  {"plans.is_active"},
		"created_at": {"plans.created_at"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		// Search:    filters.Search,
		SortOrder: filters.SortOrder,
		SortBy:    filters.SortBy,
		SearchBy:  filters.SearchBy,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "plans.id",
		SortOrder: "DESC",
	})

	filters.SetPaginationDefaults()

	data, meta, err := q.Paginate(perPage, page, "plans.id")
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 0 {
		return []map[string]any{}, meta, nil
	}

	planIDs := make([]any, 0, len(data))
	for _, row := range data {
		planIDs = append(planIDs, row["id"])
	}

	prices, err := xqb.Table("plan_prices").
		WithContext(c).
		Select("plan_id", "price", "currency").
		WhereIn("plan_id", planIDs).
		OrderBy("created_at", "ASC").
		Get()
	if err != nil {
		return nil, nil, err
	}

	grouped := make(map[any][]map[string]any)
	for _, p := range prices {
		pid := p["plan_id"]
		grouped[pid] = append(grouped[pid], p)
	}

	for _, plan := range data {
		pid := plan["id"]
		plan["plan_prices"] = grouped[pid]
	}

	return data, meta, nil
}

func (t *PlanRepository) Create(c *gin.Context, plans []*models.Plan, tx *sql.Tx) (insertedId int64, err error) {
	if len(plans) == 0 {
		return 0, xerr.New("can not create plan with empty plan list", enums.XErrBadRequestError, nil)
	}

	insertedValues := make([]map[string]any, 0)
	for _, plan := range plans {
		insertTime := time.Now()
		insertedValues = append(insertedValues, map[string]any{
			"name":       plan.Name,
			"features":   plan.Features,
			"is_active":  plan.IsActive,
			"popular":    plan.Popular,
			"created_at": insertTime,
			"updated_at": insertTime,
		})
	}

	return xqb.Table("plans").WithContext(c).WithTx(tx).InsertGetId(insertedValues)
}

func (t *PlanRepository) Update(c *gin.Context, planId int64, updatedFields map[string]any, tx *sql.Tx) error {
	if planId == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	if len(updatedFields) == 0 {
		return nil
	}

	_, err := xqb.Table("plans").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", planId).
		Update(updatedFields)
	return err
}

func (t *PlanRepository) Delete(c *gin.Context, id int64, tx *sql.Tx) error {
	if id == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	_, err := xqb.Table("plans").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", id).
		Delete()

	if err != nil {
		return err
	}

	return nil
}

func (t *PlanRepository) GetAll(c *gin.Context) ([]map[string]any, error) {
	return xqb.Table("plans").
		WithContext(c).
		Select("plans.*").
		AddSelect(models.PlanSetting{}.Cols()...).
		Join("plan_settings", "plan_settings.plan_id = plans.id").
		Get()
}

func (t *PlanRepository) GetAllActive(c *gin.Context) ([]map[string]any, error) {
	return xqb.Table("plans").
		WithContext(c).
		Select("plans.*").
		AddSelect(models.PlanSetting{}.Cols()...).
		WhereTrue("plans.is_active").
		WhereFalse("plans.hidden").
		Join("plan_settings", "plan_settings.plan_id = plans.id").
		Get()
}

func (r *PlanRepository) FindById(c *gin.Context, id int64) (*models.Plan, error) {
	planData, err := xqb.Table("plans").
		Select("plans.*").
		AddSelect(models.PlanSetting{}.Cols()...).
		WithContext(c).
		Where("plans.id", "=", id).
		Join("plan_settings", "plan_settings.plan_id = plans.id").
		First()

	if err != nil {
		return nil, err
	}

	var plan models.Plan
	if err = xqb.Bind(planData, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}
