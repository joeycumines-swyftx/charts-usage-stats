package main

import (
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_parseGetBarsURL_success(t *testing.T) {
	p := func(s string) *schema.ApiAccessLog_GetBars {
		var msg schema.ApiAccessLog_GetBars
		if err := protojson.Unmarshal([]byte(s), &msg); err != nil {
			t.Fatal(err)
		}
		return &msg
	}
	for _, tc := range [...]struct {
		Name string
		URL  string
		Msg  *schema.ApiAccessLog_GetBars
	}{
		{
			Name: `ask`,
			URL:  "/charts/v3/getBars/AUD/SOL/ask?resolution=1h\u0026timeEnd=1660522748765\u0026timeStart=1660436348765",
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"SOL", "marketSide":"ASK", "startTime":"2022-08-14T00:19:08.765Z", "endTime":"2022-08-15T00:19:08.765Z", "resolution":"3600s"}`),
		},
		{
			Name: `bid`,
			URL:  "/charts/v3/getBars/AUD/SOL/bid?resolution=1h\u0026timeEnd=1660522748765\u0026timeStart=1660436348765",
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"SOL", "marketSide":"BID", "startTime":"2022-08-14T00:19:08.765Z", "endTime":"2022-08-15T00:19:08.765Z", "resolution":"3600s"}`),
		},
		{
			Name: `trailing slash`,
			URL:  "/charts/v3/getBars/AUD/SOL/ask/?resolution=1h\u0026timeEnd=1660522748765\u0026timeStart=1660436348765",
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"SOL", "marketSide":"ASK", "startTime":"2022-08-14T00:19:08.765Z", "endTime":"2022-08-15T00:19:08.765Z", "resolution":"3600s"}`),
		},
		{
			Name: `no query params`,
			URL:  "/charts/v3/getBars/AUD/SOL/ask",
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"SOL", "marketSide":"ASK"}`),
		},
		{
			Name: `one day resolution`,
			URL:  "/charts/v3/getBars/BTC/AUD/bid?resolution=1d",
			Msg:  p(`{"primaryAsset":"BTC","secondaryAsset":"AUD","marketSide":"BID","resolution":"86400s"}`),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			msg, err := parseGetBarsURL(tc.URL)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(msg, tc.Msg) {
				b, _ := protojson.Marshal(msg)
				t.Errorf("unexpected msg: %q\n%s", b, b)
			}
		})
	}
}
