package enums

// UserRole is the role of the user
type UserRole string

const (
	RoleSuperManager UserRole = "super_manager"
	RoleManager      UserRole = "manager"
	RoleSuperAdmin   UserRole = "super_admin"
	RoleAdmin        UserRole = "admin"
	RoleClient       UserRole = "client"
)

// UserClientType is the type of the client if the user have another table to refer to
type UserClientType string
