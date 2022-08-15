package main

import (
	"bufio"
	"encoding/json"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"log"
	"os"
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
	dec := streampb.NewDecoder(src)

	enc := json.NewEncoder(dst)

	for {
		var event schema.Event
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		b, err := protojson.Marshal(&event)
		if err != nil {
			return err
		}

		if err := enc.Encode(json.RawMessage(b)); err != nil {
			return err
		}
	}

	return nil
}
