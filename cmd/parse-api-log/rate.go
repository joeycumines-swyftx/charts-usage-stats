package main

import (
	"encoding/json"
	"fmt"
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"strings"
)

func parseRateURL(URLStr string) (*schema.ApiAccessLog_Rate, error) {
	const prefix = `/charts/v3/rate/`

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

	var timestamp *timestamppb.Timestamp
	{
		var epoch msgutil.Epoch
		if err := json.Unmarshal([]byte(parts[2]), &epoch); err == nil && epoch.Valid {
			timestamp = timestamppb.New(epoch.Value)
		}
	}

	return &schema.ApiAccessLog_Rate{
		PrimaryAsset:   parts[0],
		SecondaryAsset: parts[1],
		Timestamp:      timestamp,
	}, nil
}
