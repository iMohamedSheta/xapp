package users

import (
	"context"
	"database/sql"
	"fmt"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser - hash user request password and create a new user in database
func (a *UserRepository) Create(c context.Context, user *models.User, tx *sql.Tx) (*models.User, error) {
	insertTime := time.Now()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	var tenantId *int64
	if user.TenantId == 0 {
		tenantId = nil
	} else {
		tenantId = &user.TenantId
	}

	status := enums.UserStatusActive
	if user.Status != 0 {
		status = user.Status
	}

	insertedId, err := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		InsertGetId([]map[string]any{
			{
				"tenant_id":         tenantId,
				"client_id":         user.ClientId,
				"client_type":       user.ClientType,
				"username":          user.Username,
				"name":              user.Name,
				"email":             user.Email,
				"password":          hashedPassword,
				"role":              user.Role,
				"status":            status,
				"provider":          user.Provider,
				"provider_id":       user.ProviderId,
				"avatar":            user.Avatar,
				"email_verified_at": user.EmailVerifiedAt,
				"created_at":        insertTime,
				"updated_at":        insertTime,
			},
		})
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:        insertedId,
		Name:      user.Name,
		Email:     user.Email,
		Role:      enums.UserRole(user.Role),
		Status:    status,
		CreatedAt: insertTime,
		UpdatedAt: insertTime,
	}, nil
}

func (r *UserRepository) ForceDelete(c *gin.Context, tenantId *int64, id int64, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx)

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	_, err := q.Where("id", "=", id).Delete()
	return err
}

func (r *UserRepository) SoftDelete(c context.Context, tenantId *int64, id int64, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", id)

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	_, err := q.Update(map[string]any{
		"deleted_at": time.Now(),
	})
	return err
}

func (r *UserRepository) SoftDeleteByClientIdAndType(c context.Context, tenantId *int64, clientId int64, clientType enums.UserClientType, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where("client_id", "=", clientId).
		Where("client_type", "=", clientType)

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	_, err := q.Update(map[string]any{
		"deleted_at": time.Now(),
	})
	return err
}

func (r *UserRepository) RestoreByClientIdAndType(c context.Context, tenantId *int64, clientId int64, clientType enums.UserClientType, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where("client_id", "=", clientId).
		Where("client_type", "=", clientType).
		WhereNotNull("deleted_at")

	if tenantId != nil {
		q.Where("tenant_id", "=", *tenantId)
	}

	_, err := q.Update(map[string]any{
		"deleted_at": nil,
	})
	return err
}

// Find a user by email.
func (h *UserRepository) FindUserByEmail(c context.Context, email string) (*models.User, error) {
	userData, err := xqb.Table("users").
		WithContext(c).
		WhereNull("deleted_at").
		Where(xqb.Lower("email", ""), "=", strings.ToLower(email)).
		First()
	if err != nil {
		return nil, err
	}

	var user models.User
	_ = xqb.Bind(userData, &user)

	return &user, nil
}

// Find a user by email.
func (h *UserRepository) FindUserByUsername(c context.Context, username string) (*models.User, error) {
	userData, err := xqb.Table("users").
		WithContext(c).
		WhereNull("deleted_at").
		Where("username", "=", username).
		First()
	if err != nil {
		return nil, err
	}

	var user models.User
	_ = xqb.Bind(userData, &user)

	return &user, nil
}

// Find a user by oauth provider
func (h *UserRepository) FindByOAuthProvider(c context.Context, provider string, providerId string) (*models.User, error) {
	userData, err := xqb.Table("users").
		WithContext(c).
		WhereNull("deleted_at").
		Where("provider", "=", provider).
		Where("provider_id", "=", providerId).
		First()
	if err != nil {
		return nil, err
	}

	var user models.User
	_ = xqb.Bind(userData, &user)

	return &user, nil
}

func (r *UserRepository) FindById(c context.Context, tenantId *int64, userId int64) (*models.User, error) {
	q := xqb.Model[models.User]().
		WithContext(c).
		WhereNull("deleted_at").
		Where("id", "=", userId)

	if tenantId != nil {
		q.Where("tenant_id", "=", tenantId)
	}

	return q.First()
}

func (r *UserRepository) FindClientUserByClientIdAndType(c context.Context, tenantId *int64, clientId int64, clientType enums.UserClientType) (*models.User, error) {
	q := xqb.Model[models.User]().
		WithContext(c).
		WhereNull("deleted_at").
		Where("client_id", "=", clientId).
		Where("client_type", "=", clientType)

	if tenantId != nil {
		q.Where("tenant_id", "=", tenantId)
	}

	return q.First()
}

func (r *UserRepository) FindTeamMemberWithDetails(c context.Context, tenantId *int64, userId int64) (map[string]any, error) {
	q := xqb.Table("users").
		WithContext(c).
		WhereNull("deleted_at").
		Where("users.id", "=", userId)

	if tenantId != nil && *tenantId > 0 {
		q.Where("users.tenant_id", "=", *tenantId)
	}

	q.Select("users.*")

	// Subquery for permissions
	q.AddSelect(xqb.Raw(fmt.Sprintf("COALESCE((SELECT settings->'permissions' FROM settings WHERE model_id = users.id AND model = 'user' AND type = '%s'), '[]'::jsonb) as permissions", enums.SettingPermissions)))

	return q.First()
}

func (r *UserRepository) FindByTenantId(c context.Context, tenantId int64, role enums.UserRole) (*models.User, error) {
	userData, err := xqb.Table("users").
		WithContext(c).
		WhereNull("deleted_at").
		// Where("tenant_id", "=", tenantId).
		Where("role", "=", role).
		First()
	if err != nil {
		return nil, err
	}
	var user models.User
	xqb.Bind(userData, &user)
	return &user, nil
}

func (r *UserRepository) Paginate(c *gin.Context, tenantId *int64, filters *UserFilters) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &UserFilters{}
	}

	q := xqb.Table("users").
		WithContext(c).
		WhereNull("users.deleted_at")

	if tenantId != nil && *tenantId > 0 {
		q.Where("users.tenant_id", "=", *tenantId)
	}

	// Global search across key columns
	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("users.name", "LIKE", like).
				OrWhere("users.email", "LIKE", like)
		})
	}

	if len(filters.UsersRole) > 0 {
		q.WhereIn("users.role", filters.UsersRole)
	}

	// Per-column search_by using mapped fields, like orders
	mappedFields := map[string][]string{
		"name":       {"users.name"},
		"email":      {"users.email"},
		"created_at": {"users.created_at"},
		"id":         {"users.id"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "users.id",
		SortOrder: "DESC",
	})

	// Explicitly select all user fields
	q.Select("users.*")

	// Subquery for permissions with COALESCE to ensure clean JSON array [] instead of null or "{}"
	q.AddSelect(xqb.Raw(fmt.Sprintf("COALESCE((SELECT settings->'permissions' FROM settings WHERE model_id = users.id AND model = 'user' AND type = '%s'), '[]'::jsonb) as permissions", enums.SettingPermissions)))

	filters.SetPaginationDefaults()

	q.OrderBy("users.id", "ASC")

	return q.Paginate(filters.PerPage, filters.Page, "users.id")
}

// PaginateForManager - For managers to view all users with role filtering and tenant details
func (r *UserRepository) PaginateForManager(c *gin.Context, filters *UserManagerFilters) (data []map[string]any, meta map[string]any, err error) {
	if filters == nil {
		filters = &UserManagerFilters{}
	}

	q := xqb.Table("users").
		WithContext(c).
		Select("users.*").
		AddSelect(models.Tenant{}.Cols()...).
		WhereNull("users.deleted_at").
		LeftJoin("agencies", "users.tenant_id = agencies.id")

	// Filter by role if specified
	if filters.Role != "" {
		q.Where("users.role", "=", filters.Role)
	}

	// Filter by status if specified
	if filters.Status != 0 {
		q.Where("users.status", "=", filters.Status)
	}

	if filters.TenantId != 0 {
		q.Where("users.tenant_id", "=", filters.TenantId)
	}

	// Global search across key columns
	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		q.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("CAST(users.name AS TEXT)", "LIKE", like).
				OrWhere("CAST(users.email AS TEXT)", "LIKE", like).
				OrWhere("CAST(agencies.name AS TEXT)", "LIKE", like)
		})
	}

	// Per-column search_by using mapped fields
	mappedFields := map[string][]string{
		"name":        {"users.name"},
		"email":       {"users.email"},
		"role":        {"users.role"},
		"status":      {"users.status"},
		"tenant_name": {"agencies.name"},
		"created_at":  {"users.created_at"},
		"id":          {"users.id"},
	}

	utils.ApplyFilters(q, &utils.FilterOptions{
		SearchBy:  filters.SearchBy,
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
	}, mappedFields, &utils.FilterOptions{
		SortBy:    "users.id",
		SortOrder: "DESC",
	})

	filters.SetPaginationDefaults()

	return q.Paginate(filters.PerPage, filters.Page, "users.id")
}

// GetManagerUserStats - Get user statistics for manager dashboard
func (r *UserRepository) GetManagerUserStats(c *gin.Context, filters *UserManagerFilters) (map[string]any, error) {
	q := xqb.Table("users").WithContext(c).WhereNull("users.deleted_at")

	// Apply same filters as pagination
	if filters != nil {
		if filters.Role != "" {
			q.Where("role", "=", filters.Role)
		}
		if filters.Status != 0 {
			q.Where("status", "=", filters.Status)
		}
		// if filters.TenantId != 0 {
		// 	q.Where("tenant_id", "=", filters.TenantId)
		// }
	}

	row, err := q.Select(
		xqb.Raw(`
			COUNT(*) AS total_users,
			COUNT(*) FILTER (WHERE role = ?) AS total_clients,
			COUNT(*) FILTER (WHERE role IN (?, ?)) AS total_admins,
			COUNT(*) FILTER (WHERE role IN (?, ?)) AS total_managers,
			COUNT(*) FILTER (WHERE status = ?) AS active_users,
			COUNT(*) FILTER (WHERE status = ?) AS blocked_users,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') AS new_users_week,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') AS new_users_month,
			COUNT(DISTINCT tenant_id) AS total_agencies
		`,
			enums.RoleClient,
			enums.RoleSuperAdmin, enums.RoleAdmin,
			enums.RoleSuperManager, enums.RoleManager,
			enums.UserStatusActive,
			enums.UserStatusBlocked,
		),
	).First()

	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetTeamStats - Get user statistics for team management
func (r *UserRepository) GetTeamStats(c *gin.Context, tenantId int64) (map[string]any, error) {
	q := xqb.Table("users").
		WithContext(c).
		WhereNull("users.deleted_at")

	// if tenantId > 0 {
	// 	q = q.Where("tenant_id", "=", tenantId)
	// }

	row, err := q.Select(
		xqb.Raw(`
				COUNT(*) FILTER (WHERE role IN (?,?)) AS total_users,
				COUNT(*) FILTER (WHERE role IN (?,?))  AS total_admins,
			`,
			enums.RoleSuperAdmin, enums.RoleAdmin,
			enums.RoleSuperAdmin, enums.RoleAdmin,
		),
	).First()

	if err != nil {
		return nil, err
	}

	return row, nil
}

// UpdateTenant - Update existing tenant record in database
func (a *UserRepository) Update(c context.Context, userId int64, fields map[string]any, tx *sql.Tx) error {

	if fields["updated_at"] == nil {
		fields["updated_at"] = time.Now()
	}

	_, err := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where("id", "=", userId).
		Update(fields)

	if err != nil {
		return err
	}

	return nil
}

// Dashboard methods for users (tenant admins)
func (r *UserRepository) GetDashboardMetrics(c *gin.Context, tenantId int64, filters *requests.DashboardFilters) (map[string]any, error) {
	q := xqb.Table("users").
		WithContext(c).
		WhereNull("users.deleted_at")
		// Where("tenant_id", "=", tenantId)

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

	row, err := q.Select(
		xqb.Raw(`
			COUNT(*) AS total_admins,
			COUNT(*) FILTER (WHERE status = ?) AS active_admins
		`, enums.UserStatusActive),
	).First()

	if err != nil {
		return nil, err
	}

	metrics := map[string]any{
		"total_admins":  row["total_admins"],
		"active_admins": row["active_admins"],
	}

	return metrics, nil
}

func (r *UserRepository) GetManagerMetrics(c *gin.Context, filters *requests.DashboardFilters) (map[string]any, error) {
	q := xqb.Table("users").WithContext(c).WhereNull("users.deleted_at")

	// Apply date filter
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
			COUNT(*) AS total_users,
			COUNT(*) FILTER (WHERE role = ?) AS total_clients,
			COUNT(*) FILTER (WHERE role IN (?, ?)) AS total_admins,
			COUNT(*) FILTER (WHERE role IN (?, ?)) AS total_managers,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') AS new_users_week
		`, enums.RoleClient, enums.RoleSuperAdmin, enums.RoleAdmin, enums.RoleSuperManager, enums.RoleManager),
	).First()
}

func (r *UserRepository) DeactivateByTenantId(c *gin.Context, tenantId int64, tx *sql.Tx) error {
	_, err := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		WhereNull("users.deleted_at").
		// Where("tenant_id", "=", tenantId).
		Update(map[string]any{
			"status": enums.UserStatusBlocked,
		})
	return err
}

func (r *UserRepository) ReactivateByTenantId(c *gin.Context, tenantId int64, tx *sql.Tx) error {
	_, err := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		WhereNull("users.deleted_at").
		// Where("tenant_id", "=", tenantId).
		Update(map[string]any{
			"status": enums.UserStatusActive,
		})
	return err
}

func (r *UserRepository) GenerateUsername(c context.Context, prefix string, maxAttempts int, len int, tx *sql.Tx) (string, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {

		// Generate 10 candidate usernames
		candidates := make([]string, 0, 10)
		for range 10 {
			suffix, _ := utils.RandomString(len, true, true, false)
			candidates = append(candidates, prefix+suffix)
		}

		// Check existing usernames
		rows, err := xqb.Table("users").
			WithContext(c).
			Select("username").
			WhereIn("username", candidates).
			Get()

		if err != nil {
			return "", err
		}

		exists := map[string]bool{}
		for _, row := range rows {
			if val, ok := row["username"].(string); ok {
				exists[val] = true
			}
		}

		// pick first available
		for _, username := range candidates {
			if !exists[username] {
				return username, nil
			}
		}
	}

	return "", fmt.Errorf("failed to generate unique username after max attempts")
}

// GetStaffUsers - Get all staff users (not clients) for an tenant
func (r *UserRepository) GetStaffUsers(c context.Context, tenantId int64) ([]models.User, error) {
	q := xqb.Model[models.User]().
		WithContext(c).
		WhereNull("users.deleted_at").
		Where("tenant_id", "=", tenantId).
		Where("role", "!=", enums.RoleClient).
		Where("status", "=", enums.UserStatusActive)

	return q.Get()
}
func (r *UserRepository) ValidateUniqueEmail(c context.Context, email string, exceptId *int64, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where(xqb.Lower("email", ""), "=", strings.ToLower(email)).
		WhereNull("deleted_at")

	if exceptId != nil {
		q.Where("id", "!=", *exceptId)
	}

	count, err := q.Count("id")
	if err != nil {
		return err
	}

	if count > 0 {
		return xerr.New("email already exists", enums.XErrValidationError, nil).WithDetails(map[string]any{
			"email": "البريد الإلكتروني مستخدم من قبل مستخدم اخر",
		})
	}

	return nil
}

func (r *UserRepository) ValidateUniqueUsername(c context.Context, username string, exceptId *int64, tx *sql.Tx) error {
	q := xqb.Table("users").
		WithContext(c).
		WithTx(tx).
		Where("username", "=", username).
		WhereNull("deleted_at")

	if exceptId != nil {
		q.Where("id", "!=", *exceptId)
	}

	count, err := q.Count("id")
	if err != nil {
		return err
	}

	if count > 0 {
		return xerr.New("username already exists", enums.XErrValidationError, nil).WithDetails(map[string]any{
			"username": "اسم المستخدم مستخدم من قبل مستخدم اخر",
		})
	}

	return nil
}

// ListTenantAdmins returns active super-admin and admin users for an tenant.
func (r *UserRepository) ListTenantAdmins(ctx context.Context, tenantID int64) ([]models.User, error) {
	if tenantID <= 0 {
		return nil, nil
	}
	return xqb.Model[models.User]().
		WithContext(ctx).
		Where("tenant_id", "=", tenantID).
		WhereNull("deleted_at").
		WhereIn("role", []any{enums.RoleSuperAdmin, enums.RoleAdmin}).
		Get()
}

func (r *UserRepository) ListManagers(ctx context.Context) ([]models.User, error) {
	return xqb.Model[models.User]().
		WithContext(ctx).
		WhereNull("deleted_at").
		WhereIn("role", []any{enums.RoleManager, enums.RoleSuperManager}).
		Get()
}
