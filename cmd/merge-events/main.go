package main

import (
	"bufio"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"io"
	"log"
	"os"
)

const (
	inputBufSize = 4096
)

type (
	inputDecoder struct {
		dec  *streampb.Decoder
		next *schema.Event
	}
)

func main() {
	log.SetFlags(0)
	if err := func() error {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		if err := merge(w, os.Args[1:]...); err != nil {
			return err
		}
		return w.Flush()
	}(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func merge(dst io.Writer, files ...string) error {
	enc := streampb.NewEncoder(dst)

	var inputs []*inputDecoder
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		inputs = append(inputs, &inputDecoder{dec: streampb.NewDecoder(bufio.NewReaderSize(f, inputBufSize))})
	}

	for {
		// find the input with the earliest event
		var earliest *inputDecoder

		for _, input := range inputs {
			event, err := input.peek()
			if err != nil {
				if err == io.EOF {
					continue
				}
				return err
			}
			if earliest == nil || msgutil.TimestampLess(event.GetTimestamp(), earliest.next.GetTimestamp()) {
				earliest = input
			}
		}

		if earliest == nil {
			// no more inputs
			break
		}

		// write out the earliest event
		if err := enc.Encode(earliest.next); err != nil {
			return err
		}

		// consume the event
		earliest.next = nil
	}

	return nil
}

func (x *inputDecoder) peek() (*schema.Event, error) {
	if x.next != nil {
		return x.next, nil
	}
	var event schema.Event
	if err := x.dec.Decode(&event); err != nil {
		return nil, err
	}
	x.next = &event
	return x.next, nil
}
