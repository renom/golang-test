package task

import "errors"

type Task struct {
	ID            string `sql:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Payload       string `gorm:"size:1024;not null" json:"payload"`
	HashRoundsCnt int    `gorm:"not null" json:"hash_rounds_cnt"`
	Status        Status `gorm:"not null" json:"status"`
	Hash          string `gorm:"size:128" json:"hash,omitempty"`
}

type Status int

const (
	StatusInProgress Status = iota
	StatusFinished
)

func (s Status) MarshalText() (text []byte, err error) {
	switch s {
	case StatusInProgress:
		return []byte("in progress"), nil
	case StatusFinished:
		return []byte("finished"), nil
	default:
		return []byte{}, errors.New("Incorrect value of the status field")
	}
}
