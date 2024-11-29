package matchingservice

import (
	"context"
	"fmt"
	"gameApp/entity"
	"gameApp/param"
	"gameApp/pkg/protobufencoder"
	"gameApp/pkg/richerror"
	"gameApp/pkg/timestamp"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

type Publisher interface {
	Publish(event entity.Event, payload string)
}

type Repo interface {
	AddToWatingList(userID uint, category entity.Category) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	RemoveUsersFromWaitingList(category entity.Category, userIDs []uint)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error)
}

type Service struct {
	config         Config
	repo           Repo
	PresenceClient PresenceClient
	pub            Publisher
}

type Config struct {
	WatingTimeout time.Duration `koanf:"wating_timeout"`
}

func New(config Config, repo Repo, PresenceClient PresenceClient, pub Publisher) Service {
	return Service{
		config:         config,
		repo:           repo,
		PresenceClient: PresenceClient,
		pub:            pub,
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

func (s Service) MatchWaitedUsers(ctx context.Context, _ param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	const op = "matchingservice.MatchWaitedUsers"

	var wg sync.WaitGroup
	for _, category := range entity.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}
	wg.Wait()
	return param.MatchWaitedUsersResponse{}, nil
}

func (s Service) match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	defer wg.Done()
	if err != nil {
		// TODO- log err
		// TODO- upate metrics
		log.Errorf("GetWaitingListByCategory err: %v", err)
		return
	}
	userIDs := make([]uint, 0)
	for _, user := range list {
		userIDs = append(userIDs, user.UserID)
	}

	if len(userIDs) < 2 {
		return
	}

	presenceList, err := s.PresenceClient.GetPresence(ctx, param.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {
		// TODO- log err
		// TODO- upate metrics
		log.Errorf("Get PresenceClient.GetPresence err: %v", err)
		return
	}

	presenceUserIDs := make([]uint, 0, len(list))
	for _, l := range presenceList.Items {
		userIDs = append(presenceUserIDs, l.UserID)
	}
	var toBeRemovedUsers = make([]uint, 0)
	var finalList = make([]entity.WaitingMember, 0)
	for _, l := range list {
		LastOnlineTimestamp, ok := getPresenceitem(presenceList, l.UserID)
		// TODO - add 20 , 300 to config
		if ok && LastOnlineTimestamp < timestamp.Add(-20*time.Second) && l.TimeStamp > timestamp.Add(-300*time.Second) {
			finalList = append(finalList, l)
		} else {
			toBeRemovedUsers = append(toBeRemovedUsers, l.UserID)
		}
	}

	go s.repo.RemoveUsersFromWaitingList(category, toBeRemovedUsers)

	matchedUsersToBeRemoved := make([]uint, 0)

	for i := 0; i < len(finalList)-1; i = i + 2 {

		mu := entity.MatchedUsers{
			Category: category,
			UserIDs:  []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		fmt.Println("mu", mu)
		go s.pub.Publish(entity.MatchingUsersMatchedEvent, protobufencoder.EncodeEvent(entity.MatchingUsersMatchedEvent, mu))
		// publish a new event for mu

		//remove mu users from wiating list
		matchedUsersToBeRemoved = append(matchedUsersToBeRemoved, mu.UserIDs...)
	}
	go s.repo.RemoveUsersFromWaitingList(category, matchedUsersToBeRemoved)
}

func getPresenceitem(presenceList param.GetPresenceResponse, userID uint) (int64, bool) {
	for _, itemt := range presenceList.Items {
		if itemt.UserID == userID {
			return itemt.TimeStamp, true
		}
	}
	return 0, false
}
