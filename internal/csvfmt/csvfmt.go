package csvfmt

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
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

func FormatTimestamp(t *timestamp.Timestamp) string {
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
