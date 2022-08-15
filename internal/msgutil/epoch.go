package msgutil

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"
)

const (
	milliToNano = 1_000_000
)

type (
	Epoch struct {
		Value time.Time
		Valid bool
	}
)

func (x *Epoch) UnmarshalJSON(b []byte) error {
	var value *json.Number
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	if value == nil {
		x.Value = time.Time{}
		x.Valid = false
		return nil
	}

	if v, err := value.Int64(); err == nil {
		if v >= math.MinInt64/milliToNano && v <= math.MaxInt64/milliToNano {
			x.Value = time.Unix(0, v*milliToNano)
			x.Valid = true
			return nil
		}
	} else if v, ok := new(big.Float).SetString(string(*value)); ok && !v.IsInf() {
		v.Mul(v, new(big.Float).SetInt64(milliToNano))
		if v.Cmp(new(big.Float).SetInt64(math.MinInt64)) >= 0 && v.Cmp(new(big.Float).SetInt64(math.MaxInt64)) <= 0 {
			v, _ := v.Int64()
			x.Value = time.Unix(0, v)
			x.Valid = true
			return nil
		}
	}

	return fmt.Errorf(`invalid ms epoch: %s`, value)
}

func (x Epoch) MarshalJSON() ([]byte, error) {
	if x.Valid {
		return json.Marshal(x.Value.UnixNano() / milliToNano)
	}
	return []byte(`null`), nil
}
