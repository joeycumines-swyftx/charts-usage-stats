package main

import (
	"fmt"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"net/url"
	"strings"
)

func parseLastKnownPriceURL(URLStr string) (*schema.ApiAccessLog_LastKnownPrice, error) {
	const prefix = `/charts/v3/lastKnownPrice/`

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
	case 2:
	case 3:
		if parts[2] == `` {
			parts = parts[:2]
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

	return &schema.ApiAccessLog_LastKnownPrice{
		PrimaryAsset:   parts[0],
		SecondaryAsset: parts[1],
	}, nil
}
