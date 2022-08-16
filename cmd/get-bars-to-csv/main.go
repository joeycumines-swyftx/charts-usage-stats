package main

import (
	"bufio"
	"encoding/csv"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/csvfmt"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"io"
	"log"
	"os"
	"time"
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
			csvfmt.FormatTimestamp(event.GetTimestamp()),
			csvfmt.FormatInt(int64(event.GetApi().GetAccessLog().GetStatusCode()), 10),
			csvfmt.FormatInt(event.GetApi().GetAccessLog().GetContentLength(), 10),
			csvfmt.FormatDurationRat(event.GetApi().GetAccessLog().GetDuration(), time.Millisecond, 3),
			event.GetApi().GetAccessLog().GetGetBars().GetPrimaryAsset(),
			event.GetApi().GetAccessLog().GetGetBars().GetSecondaryAsset(),
			csvfmt.FormatMarketSide(event.GetApi().GetAccessLog().GetGetBars().GetMarketSide()),
			csvfmt.FormatTimestamp(event.GetApi().GetAccessLog().GetGetBars().GetStartTime()),
			csvfmt.FormatTimestamp(event.GetApi().GetAccessLog().GetGetBars().GetEndTime()),
			csvfmt.FormatDurationFloat(event.GetApi().GetAccessLog().GetGetBars().GetResolution(), time.Minute, -1),
		}); err != nil {
			return err
		}
	}

	enc.Flush()
	return enc.Error()
}
