package authorizationservice

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/richerror"
)

type Repository interface {
	GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo}
}

func (s Service) CheckAccess(UserID uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {
	const op = "authorizationservice.CheckAccess"
	permissionTitles, err := s.repo.GetUserPermissionTitles(UserID, role)
	fmt.Println(role)
	if err != nil {
		return false, richerror.New(op).WithErr(err)
	}
	for _, pt := range permissionTitles {
		for _, p := range permissions {
			if p == pt {
				return true, nil
			}
		}
	}
	return false, nil
}
