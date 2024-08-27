package userservice

import (
	"gameApp/dto"
	"gameApp/pkg/richerror"
)

func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userService.Profile"
	user, err := s.repo.GetUserById(req.UserID)
	if err != nil {

		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

	}

	return dto.ProfileResponse{Name: user.Name}, nil
}
