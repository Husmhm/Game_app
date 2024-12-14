package protobufencoder

import (
	"encoding/base64"
	"gameApp/contract/protogolang/matching"
	"gameApp/entity"
	"gameApp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func EncodeEvent(event entity.Event, data any) string {
	var payload []byte
	switch event {
	case entity.MatchingUsersMatchedEvent:
		mu, ok := data.(entity.MatchedUsers)
		if !ok {
			//TODO - log
			//TODO - update metrics
			return ""
		}
		pbmu := matching.MatchedUsers{
			Category: string(mu.Category),
			UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
		}
		var err error
		payload, err = proto.Marshal(&pbmu)
		if err != nil {
			//TODO - log
			//TODO - update metrics
			return ""
		}
	}

	payloadStr := base64.StdEncoding.EncodeToString(payload)
	return payloadStr
}

func DecodeEvent(event entity.Event, data string) any {
	switch event {
	case entity.MatchingUsersMatchedEvent:
		payload, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			//TODO - log
			//TODO - update metrics
			return ""
		}
		pbMu := matching.MatchedUsers{}
		if err := proto.Unmarshal(payload, &pbMu); err != nil {
			//TODO - log
			//TODO - update metrics
			return ""
		}
		return entity.MatchedUsers{
			Category: entity.Category(pbMu.Category),
			UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
		}
	}
	return nil
}
