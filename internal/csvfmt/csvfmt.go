package csvfmt

import (
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/big"
	"strconv"
	"time"
)

const (
	TimeFormat = time.RFC3339Nano
)

func FormatInt(value int64, base int) string {
	if value == 0 {
		return ``
	}
	return strconv.FormatInt(value, base)
}

func FormatTimestamp(t *timestamppb.Timestamp) string {
	if t == nil {
		return ``
	}
	return t.AsTime().Format(TimeFormat)
}

func FormatMarketSide(side schema.MarketSide) string {
	switch side {
	case schema.MarketSide_BID:
		return "bid"
	case schema.MarketSide_ASK:
		return "ask"
	default:
		return ""
	}
}

func FormatDurationRat(d *durationpb.Duration, u time.Duration, prec int) string {
	if d == nil {
		return ``
	}
	return big.NewRat(
		int64(d.AsDuration()),
		int64(u),
	).FloatString(prec)
}

func FormatDurationFloat(d *durationpb.Duration, u time.Duration, prec int) string {
	if d == nil {
		return ``
	}
	return new(big.Float).SetRat(big.NewRat(
		int64(d.AsDuration()),
		int64(u),
	)).Text('f', prec)
}

var EventHeader = []string{
	"timestamp",
}

func FormatEvent(msg *schema.Event) []string {
	return []string{
		FormatTimestamp(msg.GetTimestamp()),
	}
}

var APIAccessLogHeader = []string{
	"status_code",
	"content_length",
	"duration",
}

func FormatAPIAccessLog(msg *schema.ApiAccessLog) []string {
	return []string{
		FormatInt(int64(msg.GetStatusCode()), 10),
		FormatInt(msg.GetContentLength(), 10),
		FormatDurationRat(msg.GetDuration(), time.Millisecond, 3),
	}
}
