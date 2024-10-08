package entity

// Access control only keep allowed permision
type AccessControl struct {
	ID           uint
	ActorID      uint
	ActorType    ActorType
	PermissionID uint
}

type ActorType string

const (
	RoleActorType ActorType = "role"
	UserActorType ActorType = "user"
)
