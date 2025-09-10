package templates

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"strconv"
	"time"

	"github.com/aity-cloud/monty/pkg/util"
	"github.com/prometheus/common/model"
)

func floatToTime(v float64) (*time.Time, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return nil, errNaNOrInf
	}
	timestamp := v * 1e9
	if timestamp > math.MaxInt64 || timestamp < math.MinInt64 {
		return nil, fmt.Errorf("%v cannot be represented as a nanoseconds timestamp since it overflows int64", v)
	}
	t := model.TimeFromUnixNano(int64(timestamp)).Time().UTC()
	return &t, nil
}

func convertToFloat(i any) (float64, error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("can't convert %T to float", v)
	}
}

var errNaNOrInf = errors.New("value is NaN or Inf")

var DefaultTemplateFuncs = template.FuncMap{
	"humanize":     util.Humanize,
	"humanize1024": util.Humanize1024,
	"humanizePercentage": func(i any) (string, error) {
		v, err := convertToFloat(i)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%.4g%%", v*100), nil
	},
	"humanizeTimestamp": func(i any) (string, error) {
		v, err := convertToFloat(i)
		if err != nil {
			return "", err
		}

		tm, err := floatToTime(v)
		switch {
		case errors.Is(err, errNaNOrInf):
			return fmt.Sprintf("%.4g", v), nil
		case err != nil:
			return "", err
		}

		return fmt.Sprint(tm), nil
	},
	"formatTime": func(i any) (string, error) {
		ts, ok := i.(time.Time)
		if !ok {
			return "", fmt.Errorf("formatTime: expected time.Time, got %T", i)
		}
		return fmt.Sprint(ts.Format(time.RFC822)), nil
	},
	"toTime": func(i any) (*time.Time, error) {
		v, err := convertToFloat(i)
		if err != nil {
			return nil, err
		}

		return floatToTime(v)
	},
}

func RegisterNewAlertManagerDefaults[T, U ~map[string]any](amTmplMap T, newTmplMap U) {
	for key := range newTmplMap {
		if err := RegisterTemplateMap(amTmplMap, key, newTmplMap[key]); err != nil {
			panic(err)
		}
	}
}

func RegisterTemplateMap[T ~map[string]any](templateMap T, key string, templateFunc any) error {
	if _, ok := templateMap[key]; ok {
		return fmt.Errorf("key error : template function %s already exists", key)
	}
	templateMap[key] = templateFunc
	return nil
}
