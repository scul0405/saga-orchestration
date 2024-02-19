package sonyflake

import (
	"errors"

	"github.com/sony/sonyflake"
)

// IDGenerator is the interface for generating unique ID
type IDGenerator interface {
	NextID() (uint64, error)
}

// NewSonyFlake returns new SonyFlake ID generator
func NewSonyFlake() (IDGenerator, error) {
	var st sonyflake.Settings
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return nil, errors.New("sonyflake not created")
	}
	return sf, nil
}
