package param

import (
	"gameApp/entity"
	"time"
)

type AddToWatingListRequest struct {
	UserID   uint            `json:"user_id"`
	Category entity.Category `json:"category"`
}

type AddToWatingListResponse struct {
	Timeout time.Duration `json:"timeout_in_nano_second"`
}
