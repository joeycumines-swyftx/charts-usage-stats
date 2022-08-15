package cloudwatch

import (
	"github.com/joeycumines-swyftx/charts-usage-stats/internal/msgutil"
)

type (
	LogEvent struct {
		Timestamp     msgutil.Epoch `json:"timestamp"`
		IngestionTime msgutil.Epoch `json:"ingestionTime"`
		Message       string        `json:"message"`
	}
)
