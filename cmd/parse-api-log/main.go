package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/cloudwatch"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"io"
	"log"
	"os"
)

var (
	errUnhandledEvent = errors.New(`unhandled event`)
)

type (
	// AccessLog is a log event from https://docs.fluentbit.io/manual/installation/kubernetes
	AccessLog struct {
		Stream string `json:"stream"`
		Log    string `json:"log"`
	}
)

func main() {
	log.SetFlags(0)
	if err := func() error {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		if err := parse(w, bufio.NewReader(os.Stdin)); err != nil {
			return err
		}
		return w.Flush()
	}(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func parse(dst io.Writer, src io.Reader) error {
	dec := json.NewDecoder(src)
	dec.UseNumber()

	enc := streampb.NewEncoder(dst)

	for {
		var logEvent cloudwatch.LogEvent
		if err := dec.Decode(&logEvent); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		acsEvent, err := parseEvent(&logEvent)
		if err != nil {
			if errors.Is(err, errUnhandledEvent) {
				continue
			}
			return err
		}

		if err := enc.Encode(acsEvent); err != nil {
			return err
		}
	}

	return nil
}

func parseEvent(event *cloudwatch.LogEvent) (parsed *schema.Event, err error) {
	if !event.Timestamp.Valid {
		return nil, errors.New(`missing timestamp`)
	}
	if parsed, err = parseAccessLog(event); !errors.Is(err, errUnhandledEvent) {
		return parsed, err
	}
	return nil, errUnhandledEvent
}
