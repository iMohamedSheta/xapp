package tenants

import (
	"context"
	"database/sql"

	"github.com/imohamedsheta/xapp/app/models"
)

type TenantService struct {
	tenantRepo *TenantRepository
}

func NewTenantService(
	tenantRepo *TenantRepository,
) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
	}
}

func (a *TenantService) Create(ctx context.Context, model *models.Tenant, tx *sql.Tx) (*models.Tenant, error) {
	return a.tenantRepo.Create(ctx, model, tx)
}

func (a *TenantService) FindById(ctx context.Context, id int64) (*models.Tenant, error) {
	return a.tenantRepo.FindById(ctx, id)
}
