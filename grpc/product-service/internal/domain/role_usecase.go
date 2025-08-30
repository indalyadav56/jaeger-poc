package domain

type RoleUsecase interface {
	CreateRole(role *Role) error
	GetRole(id string) (*Role, error)
	UpdateRole(role *Role) error
	DeleteRole(id string) error
}
