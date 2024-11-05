package matchingservice

import (
	"gameApp/entity"
	"gameApp/param"
	"gameApp/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWatingList(userID uint, category entity.Category) error
}

type Service struct {
	config Config
	repo   Repo
}

type Config struct {
	WatingTimeout time.Duration `koanf:"wating_timeout"`
}

func New(config Config, repo Repo) Service {
	return Service{
		config: config,
		repo:   repo,
	}
}

func (s Service) AddToWatingList(req param.AddToWatingListRequest) (param.AddToWatingListResponse, error) {
	const op = "matchingservice.AddToWatingList"

	err := s.repo.AddToWatingList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWatingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	return param.AddToWatingListResponse{Timeout: s.config.WatingTimeout}, nil
}

func (s Service) MatchWaitedUsers(req param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	return param.MatchWaitedUsersResponse{}, nil
}
