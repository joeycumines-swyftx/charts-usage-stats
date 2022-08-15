package msgutil

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"
)

func TestEpoch_UnmarshalJSON(t *testing.T) {
	for _, tc := range [...]struct {
		Name   string
		Input  string
		Output Epoch
		Error  error
	}{
		{
			Name:  `null`,
			Input: `null`,
		},
		{
			Name:  `max int64`,
			Input: strconv.FormatInt(int64(math.MaxInt64/time.Millisecond), 10),
			Output: Epoch{
				Value: time.Unix(0, int64(math.MaxInt64/time.Millisecond*time.Millisecond)),
				Valid: true,
			},
		},
		{
			Name:  `max int64 overflow`,
			Input: strconv.FormatInt(int64(math.MaxInt64/time.Millisecond)+1, 10),
			Error: errors.New(`invalid ms epoch: 9223372036855`),
		},
		{
			Name:  `min int64`,
			Input: strconv.FormatInt(int64(math.MinInt64/time.Millisecond), 10),
			Output: Epoch{
				Value: time.Unix(0, int64(math.MinInt64/time.Millisecond*time.Millisecond)),
				Valid: true,
			},
		},
		{
			Name:  `min int64 overflow`,
			Input: strconv.FormatInt(int64(math.MinInt64/time.Millisecond)-1, 10),
			Error: errors.New(`invalid ms epoch: -9223372036855`),
		},
		{
			Name:  `max int64 float`,
			Input: strconv.FormatInt(int64(math.MaxInt64/time.Millisecond), 10) + `.0`,
			Output: Epoch{
				Value: time.Unix(0, 9223372036854000000),
				Valid: true,
			},
		},
		{
			Name:  `max int64 float overflow`,
			Input: strconv.FormatInt(int64(math.MaxInt64/time.Millisecond)+1, 10) + `.0`,
			Error: errors.New(`invalid ms epoch: 9223372036855.0`),
		},
		{
			Name:  `min int64 float`,
			Input: strconv.FormatInt(int64(math.MinInt64/time.Millisecond), 10) + `.0`,
			Output: Epoch{
				Value: time.Unix(0, -9223372036854000000),
				Valid: true,
			},
		},
		{
			Name:  `min int64 float overflow`,
			Input: strconv.FormatInt(int64(math.MinInt64/time.Millisecond)-1, 10) + `.0`,
			Error: errors.New(`invalid ms epoch: -9223372036855.0`),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var output Epoch
			err := json.Unmarshal([]byte(tc.Input), &output)
			if (err == nil) != (tc.Error == nil) || (err != nil && err.Error() != tc.Error.Error()) {
				t.Errorf("unexpected error: %v", err)
			}
			if output != tc.Output {
				t.Errorf("unexpected output: %+v / %d", output, output.Value.UnixNano())
			}
			if err != nil || tc.Error != nil {
				return
			}
			b, err := json.Marshal(output)
			if err != nil {
				t.Fatal(err)
			}
			output = Epoch{}
			if err := json.Unmarshal(b, &output); err != nil {
				t.Fatal(err)
			}
			if output != tc.Output {
				t.Errorf("unexpected output (after marshal): %+v / %d", output, output.Value.UnixNano())
			}
		})
	}
}

func TestEpoch_UnmarshalJSON_nullOverwrites(t *testing.T) {
	epoch := Epoch{Value: time.Unix(1234125, 0), Valid: true}
	if err := json.Unmarshal([]byte(`null`), &epoch); err != nil {
		t.Fatal(err)
	}
	if epoch != (Epoch{}) {
		t.Error(epoch)
	}
}
