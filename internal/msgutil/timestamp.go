package msgutil

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampLess(a, b *timestamppb.Timestamp) bool {
	if b == nil {
		return false
	}
	if a == nil {
		return true
	}
	return a.AsTime().Before(b.AsTime())
}
