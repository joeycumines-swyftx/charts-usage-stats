package main

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"testing"
	"time"
)

func Test_less(t *testing.T) {
	n := func(v *int64) *timestamppb.Timestamp {
		if v == nil {
			return nil
		}
		return timestamppb.New(time.Unix(0, *v))
	}
	i := func(v int64) *int64 { return &v }
	for _, tc := range [...]struct {
		Name string
		A    *int64
		B    *int64
		Want bool
	}{
		{
			Name: "nil nil",
		},
		{
			Name: "non-nil nil",
			A:    i(math.MinInt64),
		},
		{
			Name: "nil non-nil",
			B:    i(math.MinInt64),
			Want: true,
		},
		{
			Name: "0 nil",
			A:    i(0),
		},
		{
			Name: "nil 0",
			B:    i(0),
			Want: true,
		},
		{
			Name: "0 0",
			A:    i(0),
			B:    i(0),
		},
		{
			Name: "0 1",
			A:    i(0),
			B:    i(1),
			Want: true,
		},
		{
			Name: "1 0",
			A:    i(1),
			B:    i(0),
		},
		{
			Name: "-1 -1",
			A:    i(-1),
			B:    i(-1),
		},
		{
			Name: "-1 1",
			A:    i(-1),
			B:    i(1),
			Want: true,
		},
		{
			Name: "1 -1",
			A:    i(1),
			B:    i(-1),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			got := less(n(tc.A), n(tc.B))
			if got != tc.Want {
				t.Error(got)
			}
		})
	}
}
