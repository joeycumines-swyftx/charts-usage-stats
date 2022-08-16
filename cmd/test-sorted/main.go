package main

import (
	"bufio"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	ok, err := parse(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	if !ok {
		os.Exit(1)
	}
}

func parse(src io.Reader) (bool, error) {
	dec := streampb.NewDecoder(src)

	for last := (*timestamppb.Timestamp)(nil); ; {
		var event schema.Event
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return false, err
		}

		if msgutil.TimestampLess(event.GetTimestamp(), last) {
			return false, nil
		}

		last = event.GetTimestamp()
	}

	return true, nil
}
