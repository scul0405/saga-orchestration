package timeconvert

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func Time2pbTimestamp(t time.Time) *timestamp.Timestamp {
	s := t.Unix()
	n := int32(t.Nanosecond())

	return &timestamp.Timestamp{
		Seconds: s,
		Nanos:   n,
	}
}
