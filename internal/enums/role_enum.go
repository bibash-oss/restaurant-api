package enums

type UserRole string

const (
    UserRoleSuperAdmin UserRole = "SUPER_ADMIN"
    UserRoleAdmin      UserRole = "ADMIN"
    UserRoleUser       UserRole = "USER"
)