package entity

type Question struct {
	ID               uint
	Text             string
	PossibleAnswers  []PossibleAnswer
	CorrectAnswersID uint
	CategoryID       uint
	Difficulty       QuestionDifficulty
}
type PossibleAnswer struct {
	ID     uint
	Text   string
	Choice PossibleAnswereChoice
}

type PossibleAnswereChoice uint8

func (p PossibleAnswereChoice) Isvalid() bool {
	if p >= PossibleAnswereA && p <= PossibleAnswereD {
		return true
	} else {
		return false
	}
}

const (
	PossibleAnswereA = iota + 1
	PossibleAnswereB
	PossibleAnswereC
	PossibleAnswereD
)

type QuestionDifficulty uint8

const (
	QuestionDifficultyEasy QuestionDifficulty = iota + 1
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (q QuestionDifficulty) Isvalid() bool {
	if q >= QuestionDifficultyEasy && q <= QuestionDifficultyHard {
		return true
	} else {
		return false
	}
}
