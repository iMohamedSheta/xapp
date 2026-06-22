package plans

import (
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/domain/enums"
)

type PlanPriceRepository struct {
}

func NewPlanPriceRepository() *PlanPriceRepository {
	return &PlanPriceRepository{}
}

func (t *PlanPriceRepository) Create(c *gin.Context, prices []*models.PlanPrice, tx *sql.Tx) (insertedId int64, err error) {
	if len(prices) == 0 {
		return 0, xerr.New("can not create plan with empty plan list", enums.XErrBadRequestError, nil)
	}

	insertedValues := make([]map[string]any, 0)
	for _, price := range prices {
		insertTime := time.Now()
		insertedValues = append(insertedValues, map[string]any{
			"plan_id":    price.PlanId,
			"price":      price.Price,
			"currency":   price.Currency,
			"created_at": insertTime,
			"updated_at": insertTime,
		})
	}

	return xqb.Table("plan_prices").WithContext(c).WithTx(tx).InsertGetId(insertedValues)
}

func (t *PlanPriceRepository) DeleteByPlanId(c *gin.Context, planId int64, tx *sql.Tx) error {
	if planId == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	_, err := xqb.Table("plan_prices").
		WithContext(c).
		WithTx(tx).
		Where("plan_id", "=", planId).
		Delete()

	if err != nil {
		return err
	}

	return nil
}

func (r *PlanPriceRepository) UpsertByPlanID(c *gin.Context, planId int64, updatedFields []map[string]any, tx *sql.Tx) error {
	if planId == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	if len(updatedFields) == 0 {
		return nil
	}

	for i := range updatedFields {
		updatedFields[i]["plan_id"] = planId
	}

	// insert or update by unique key (plan_id + currency)
	_, err := xqb.Table("plan_prices").
		WithContext(c).
		WithTx(tx).
		Upsert(
			updatedFields,
			[]string{"plan_id", "currency"},
			[]string{"price"},
		)

	return err
}

func (t *PlanPriceRepository) GetAll(c *gin.Context) ([]map[string]any, error) {
	return xqb.Table("plan_prices").
		WithContext(c).
		Get()
}
