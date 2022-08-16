package main

import (
	"encoding/json"
	"fmt"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"strings"
	"time"
)

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
	switch len(parts) {
	case 3:
	case 4:
		if parts[3] == `` {
			parts = parts[:3]
			break
		}
		fallthrough
	default:
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
