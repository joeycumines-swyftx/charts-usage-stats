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

	var record []string
	record = append(record, csvfmt.EventHeader...)
	record = append(record, csvfmt.APIAccessLogHeader...)
	record = append(
		record,
		"primary_asset",
		"secondary_asset",
		"market_side",
		"start_time",
		"end_time",
		"resolution",
	)
	if err := enc.Write(record); err != nil {
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

		msg := event.GetApi().GetAccessLog().GetGetBars()
		if msg == nil {
			continue
		}

		record = record[:0]
		record = append(record, csvfmt.FormatEvent(&event)...)
		record = append(record, csvfmt.FormatAPIAccessLog(event.GetApi().GetAccessLog())...)
		record = append(
			record,
			msg.GetPrimaryAsset(),
			msg.GetSecondaryAsset(),
			csvfmt.FormatMarketSide(msg.GetMarketSide()),
			csvfmt.FormatTimestamp(msg.GetStartTime()),
			csvfmt.FormatTimestamp(msg.GetEndTime()),
			csvfmt.FormatDurationFloat(msg.GetResolution(), time.Minute, -1),
		)
		if err := enc.Write(record); err != nil {
			return err
		}
	}

	enc.Flush()
	return enc.Error()
}
