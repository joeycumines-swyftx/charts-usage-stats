package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cyx/streampb"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/cloudwatch"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"math/big"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errUnhandledEvent = errors.New(`unhandled event`)
	// :method :url :status :res[content-length] - :response-time ms
	morganTinyRegex = regexp.MustCompile(`^(?P<method>[A-Z]+) (?P<url>[^ ]+) (?P<status>\d+) (?P<resContentLength>\d+|-) - (?P<responseTime>\d+(?:\.\d+)?|-) ms$`)
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

func parseAccessLog(event *cloudwatch.LogEvent) (*schema.Event, error) {
	var accessLog AccessLog
	if err := json.Unmarshal([]byte(event.Message), &accessLog); err != nil {
		return nil, errUnhandledEvent
	}

	if accessLog.Stream != `stdout` || accessLog.Log == `` {
		return nil, errUnhandledEvent
	}

	match := morganTinyRegex.FindStringSubmatch(accessLog.Log)
	if match == nil {
		return nil, errUnhandledEvent
	}

	names := morganTinyRegex.SubexpNames()
	if len(names) != 6 {
		panic(`unexpected number of subexp names`)
	}
	var (
		methodIndex           int
		urlIndex              int
		statusIndex           int
		resContentLengthIndex int
		responseTimeIndex     int
	)
	for i, name := range names {
		if i == 0 {
			continue
		}
		switch name {
		case `method`:
			methodIndex = i
		case `url`:
			urlIndex = i
		case `status`:
			statusIndex = i
		case `resContentLength`:
			resContentLengthIndex = i
		case `responseTime`:
			responseTimeIndex = i
		default:
			panic(`unexpected subexp name`)
		}
	}

	var data schema.ApiAccessLog
	result := schema.Event{
		Timestamp: timestamppb.New(event.Timestamp.Value),
		Data:      &schema.Event_Api{Api: &schema.ApiEvent{Data: &schema.ApiEvent_AccessLog{AccessLog: &data}}},
	}

	switch match[methodIndex] {
	case `GET`:
		if v, err := parseGetBarsURL(match[urlIndex]); err == nil {
			data.Data = &schema.ApiAccessLog_GetBars_{GetBars: v}
			break
		} else if !errors.Is(err, errUnhandledEvent) {
			return nil, err
		}

		fallthrough

	default:
		return nil, errUnhandledEvent
	}

	if v, err := strconv.ParseInt(match[statusIndex], 10, 32); err != nil {
		return nil, err
	} else if v < 100 || v >= 600 {
		return nil, fmt.Errorf(`invalid status code: %d`, v)
	} else {
		data.StatusCode = int32(v)
	}

	if match[resContentLengthIndex] != `-` {
		v, err := strconv.ParseInt(match[resContentLengthIndex], 10, 32)
		if err != nil {
			return nil, err
		}
		data.ContentLength = v
	}

	if match[responseTimeIndex] != `-` {
		f, ok := new(big.Float).SetString(match[responseTimeIndex])
		if !ok || f.IsInf() {
			return nil, fmt.Errorf(`invalid response time: %s`, match[responseTimeIndex])
		}
		f.Mul(f, new(big.Float).SetInt64(int64(time.Millisecond)))
		v, _ := f.Int64()
		data.Duration = durationpb.New(time.Duration(v))
	}

	return &result, nil
}

func parseGetBarsURL(URLStr string) (*schema.ApiAccessLog_GetBars, error) {
	const prefix = `/charts/v3/getBars/`

	if !strings.HasPrefix(URLStr, prefix) {
		return nil, errUnhandledEvent
	}

	URL, err := url.ParseRequestURI(URLStr)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(URL.Path, prefix) {
		return nil, fmt.Errorf(`unexpected url: %s`, URL)
	}

	parts := strings.Split(URL.Path[len(prefix):], `/`)
	if len(parts) != 3 {
		return nil, errUnhandledEvent
	}
	for _, part := range parts {
		if part == `` {
			return nil, errUnhandledEvent
		}
	}

	var side schema.MarketSide
	switch parts[2] {
	case `ask`:
		side = schema.MarketSide_ASK
	case `bid`:
		side = schema.MarketSide_BID
	default:
		return nil, errUnhandledEvent
	}

	query := URL.Query()

	var startTime *timestamppb.Timestamp
	if v := query.Get(`timeStart`); v != `` {
		var epoch msgutil.Epoch
		if err := json.Unmarshal([]byte(v), &epoch); err == nil && epoch.Valid {
			startTime = timestamppb.New(epoch.Value)
		}
	}

	var endTime *timestamppb.Timestamp
	if v := query.Get(`timeEnd`); v != `` {
		var epoch msgutil.Epoch
		if err := json.Unmarshal([]byte(v), &epoch); err == nil && epoch.Valid {
			endTime = timestamppb.New(epoch.Value)
		}
	}

	var resolution *durationpb.Duration
	if v := query.Get(`resolution`); v != `` {
		if v == `1d` {
			resolution = durationpb.New(time.Hour * 24)
		} else if v, err := time.ParseDuration(v); err == nil {
			resolution = durationpb.New(v)
		}
	}

	return &schema.ApiAccessLog_GetBars{
		PrimaryAsset:   parts[0],
		SecondaryAsset: parts[1],
		MarketSide:     side,
		StartTime:      startTime,
		EndTime:        endTime,
		Resolution:     resolution,
	}, nil
}
