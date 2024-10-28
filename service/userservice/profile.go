package userservice

import (
	"context"
	"gameApp/param"
	"gameApp/pkg/richerror"
)

func (s Service) Profile(ctx context.Context, req param.ProfileRequest) (param.ProfileResponse, error) {
	const op = "userService.Profile"
	user, err := s.repo.GetUserById(ctx, req.UserID)
	if err != nil {

		return param.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

	}

	return param.ProfileResponse{Name: user.Name}, nil
}
