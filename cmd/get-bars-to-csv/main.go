package main

import (
	"bufio"
	"encoding/csv"
	"github.com/cyx/streampb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

const (
	timeFormat = time.RFC3339Nano
)

func main() {
	log.SetFlags(0)
	if err := parse(os.Stdout, bufio.NewReader(os.Stdin)); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func parse(dst io.Writer, src io.Reader) error {
	dec := streampb.NewDecoder(src)

	enc := csv.NewWriter(dst)
	defer enc.Flush()

	if err := enc.Write([]string{
		"timestamp",
		"status_code",
		"content_length",
		"duration",
		"primary_asset",
		"secondary_asset",
		"market_side",
		"start_time",
		"end_time",
		"resolution",
	}); err != nil {
		return err
	}

	for {
		var event schema.Event
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if event.GetApi().GetAccessLog().GetGetBars() == nil {
			continue
		}

		if err := enc.Write([]string{
			formatTimestamp(event.GetTimestamp()),
			formatInt(int64(event.GetApi().GetAccessLog().GetStatusCode()), 10),
			formatInt(event.GetApi().GetAccessLog().GetContentLength(), 10),
			formatDurationRat(event.GetApi().GetAccessLog().GetDuration(), time.Millisecond, 3),
			event.GetApi().GetAccessLog().GetGetBars().GetPrimaryAsset(),
			event.GetApi().GetAccessLog().GetGetBars().GetSecondaryAsset(),
			formatMarketSide(event.GetApi().GetAccessLog().GetGetBars().GetMarketSide()),
			formatTimestamp(event.GetApi().GetAccessLog().GetGetBars().GetStartTime()),
			formatTimestamp(event.GetApi().GetAccessLog().GetGetBars().GetEndTime()),
			formatDurationFloat(event.GetApi().GetAccessLog().GetGetBars().GetResolution(), time.Minute, -1),
		}); err != nil {
			return err
		}
	}

	enc.Flush()
	return enc.Error()
}

func formatInt(value int64, base int) string {
	if value == 0 {
		return ``
	}
	return strconv.FormatInt(value, base)
}

func formatTimestamp(t *timestamp.Timestamp) string {
	if t == nil {
		return ``
	}
	return t.AsTime().Format(timeFormat)
}

func formatMarketSide(side schema.MarketSide) string {
	switch side {
	case schema.MarketSide_BID:
		return "bid"
	case schema.MarketSide_ASK:
		return "ask"
	default:
		return ""
	}
}

func formatDurationRat(d *durationpb.Duration, u time.Duration, prec int) string {
	if d == nil {
		return ``
	}
	return big.NewRat(
		int64(d.AsDuration()),
		int64(u),
	).FloatString(prec)
}

func formatDurationFloat(d *durationpb.Duration, u time.Duration, prec int) string {
	if d == nil {
		return ``
	}
	return new(big.Float).SetRat(big.NewRat(
		int64(d.AsDuration()),
		int64(u),
	)).Text('f', prec)
}
