package main

import (
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_parseLastKnownPriceURL_success(t *testing.T) {
	p := func(s string) *schema.ApiAccessLog_LastKnownPrice {
		var msg schema.ApiAccessLog_LastKnownPrice
		if err := protojson.Unmarshal([]byte(s), &msg); err != nil {
			t.Fatal(err)
		}
		return &msg
	}
	for _, tc := range [...]struct {
		Name string
		URL  string
		Msg  *schema.ApiAccessLog_LastKnownPrice
	}{
		{
			Name: `integer ids`,
			URL:  `/charts/v3/lastKnownPrice/95/107`,
			Msg:  p(`{"primaryAsset":"95", "secondaryAsset":"107"}`),
		},
		{
			Name: `integer ids trailing slash`,
			URL:  `/charts/v3/lastKnownPrice/95/107/`,
			Msg:  p(`{"primaryAsset":"95", "secondaryAsset":"107"}`),
		},
		{
			Name: `string ids`,
			URL:  `/charts/v3/lastKnownPrice/AUD/LEVER`,
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"LEVER"}`),
		},
		{
			Name: `string ids trailing slash`,
			URL:  `/charts/v3/lastKnownPrice/AUD/LEVER/`,
			Msg:  p(`{"primaryAsset":"AUD", "secondaryAsset":"LEVER"}`),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			msg, err := parseLastKnownPriceURL(tc.URL)
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
