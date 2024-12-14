package broker

import "gameApp/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}
