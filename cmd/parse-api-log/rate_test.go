package main

import (
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_parseRateURL_success(t *testing.T) {
	p := func(s string) *schema.ApiAccessLog_Rate {
		var msg schema.ApiAccessLog_Rate
		if err := protojson.Unmarshal([]byte(s), &msg); err != nil {
			t.Fatal(err)
		}
		return &msg
	}
	for _, tc := range [...]struct {
		Name string
		URL  string
		Msg  *schema.ApiAccessLog_Rate
	}{
		{
			Name: `integer ids`,
			URL:  `/charts/v3/rate/95/107/1660617904567`,
			Msg:  p(`{"primaryAsset":"95", "secondaryAsset":"107", "timestamp":"2022-08-16T02:45:04.567Z"}`),
		},
		{
			Name: `integer ids trailing slash`,
			URL:  `/charts/v3/rate/95/107/1660617904567/`,
			Msg:  p(`{"primaryAsset":"95", "secondaryAsset":"107", "timestamp":"2022-08-16T02:45:04.567Z"}`),
		},
		{
			Name: `string ids`,
			URL:  `/charts/v3/rate/AUD/LEVER/1660617904567`,
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"LEVER", "timestamp":"2022-08-16T02:45:04.567Z"}`),
		},
		{
			Name: `string ids trailing slash`,
			URL:  `/charts/v3/rate/AUD/LEVER/1660617904567/`,
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"LEVER", "timestamp":"2022-08-16T02:45:04.567Z"}`),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			msg, err := parseRateURL(tc.URL)
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
