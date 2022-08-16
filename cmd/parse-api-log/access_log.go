package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/cloudwatch"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/big"
	"regexp"
	"strconv"
	"time"
)

var (
	// :method :url :status :res[content-length] - :response-time ms
	morganTinyRegex = regexp.MustCompile(`^(?P<method>[A-Z]+) (?P<url>[^ ]+) (?P<status>\d+) (?P<resContentLength>\d+|-) - (?P<responseTime>\d+(?:\.\d+)?|-) ms$`)
)

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

		if v, err := parseLastKnownPriceURL(match[urlIndex]); err == nil {
			data.Data = &schema.ApiAccessLog_LastKnownPrice_{LastKnownPrice: v}
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
