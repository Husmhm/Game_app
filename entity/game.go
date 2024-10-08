package entity

import "time"

type Game struct {
	ID          uint
	Category    Category
	QuestionIDs []uint
	PlayerIDs   []uint
	WinnerID    uint
	StaetTime   time.Time
}

type Player struct {
	ID     uint
	UserID uint
	GameID uint
	Score  uint
	Answer []PlayerAnswer
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Chice      PossibleAnswereChoice
}
