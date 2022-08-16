package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cyx/streampb"
	"github.com/erkkah/margaid"
	"github.com/joeycumines-swyftx/charts-usage-stats/schema"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	aspectRatioWidth  = 4
	aspectRatioHeight = 3
	colourScheme      = 90
	diagramWidth      = 1200
	diagramInset      = 150
	diagramPadding    = 10
	diagramMarker     = "filled-circle"
)

type (
	Config struct {
		Reader    io.Reader
		OutputDir string
		Scale     *big.Rat
	}

	eventHandler interface {
		Handle(event *schema.Event) error
		Flush() error
	}

	eventHandlerFuncs struct {
		handle func(event *schema.Event) error
		flush  func() error
	}
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func run() error {
	cfg, err := configure()
	if err != nil {
		return err
	}
	return process(cfg)
}

func configure() (*Config, error) {
	cfg := Config{
		Reader: bufio.NewReader(os.Stdin),
		// assume 1 day period if scale is 1
		// bigger scale == bigger period (for aggregations)
		Scale: big.NewRat(1, 1),
	}
	{
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		cfg.OutputDir = wd
	}
	return &cfg, nil
}

func process(cfg *Config) error {
	if cfg == nil ||
		cfg.Reader == nil ||
		cfg.OutputDir == `` ||
		cfg.Scale == nil {
		return errors.New(`invalid config`)
	}

	handlers, err := cfg.eventHandlers()
	if err != nil {
		return err
	}

	for dec := streampb.NewDecoder(cfg.Reader); ; {
		var event schema.Event
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if err := callHandlers(handlers, &event); err != nil {
			return err
		}
	}

	if err := flushHandlers(handlers); err != nil {
		return err
	}

	return nil
}

func (x *Config) eventHandlers() (handlers []eventHandler, _ error) {
	for _, factory := range [...]func() (eventHandler, error){
		x.avgRequestDuration,
		x.avgContentLength,
		x.numStatusCodes,
	} {
		handler, err := factory()
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, handler)
	}
	return
}

func (x *Config) avgRequestDuration() (eventHandler, error) {
	period, err := x.scaleDuration(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	seriesOptions := []margaid.SeriesOption{
		margaid.AggregatedBy(margaid.Avg, period),
	}
	getBars := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`get-bars`),
	)...)
	lastKnownPrice := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`last-known-price`),
	)...)
	rate := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`rate`),
	)...)

	var charts []func() error
	type Config struct {
		Title    string
		Filename string
		Range    *margaid.Series
		Series   []*margaid.Series
	}
	addChart := func(cfg Config) {
		charts = append(charts, func() error {
			const (
				diagramTitle    = `Avg Request Duration`
				diagramFilename = `avg-request-duration`
				xAxisTitle      = `Time (seconds)`
				yAxisTitle      = `Duration (milliseconds)`
			)

			if cfg.Range == nil {
				cfg.Range = mergeSeriesMinMax(cfg.Series...)
			}

			diagram := margaid.New(
				diagramWidth,
				calculateHeight(diagramWidth),
				margaid.WithAutorange(margaid.XAxis, cfg.Range),
				margaid.WithAutorange(margaid.YAxis, cfg.Range),
				margaid.WithColorScheme(colourScheme),
				margaid.WithInset(diagramInset),
				margaid.WithPadding(diagramPadding),
			)
			for _, series := range cfg.Series {
				diagram.Line(
					series,
					margaid.UsingAxes(margaid.XAxis, margaid.YAxis),
					margaid.UsingMarker(diagramMarker),
				)
			}
			diagram.Axis(cfg.Range, margaid.XAxis, diagram.ValueTicker('f', 0, 10), false, xAxisTitle)
			diagram.Axis(cfg.Range, margaid.YAxis, diagram.ValueTicker('f', 0, 2), true, yAxisTitle)
			diagram.Title(fmt.Sprintf(`%s (%s) - %s`, diagramTitle, period, cfg.Title))
			if len(cfg.Series) > 1 {
				diagram.Legend(margaid.BottomLeft)
			}

			return x.writeOutputFile(fmt.Sprintf(`%s-%s.svg`, diagramFilename, cfg.Filename), func(w io.Writer) error {
				buf := bufio.NewWriter(w)
				if err := diagram.Render(buf); err != nil {
					return err
				}
				return buf.Flush()
			})
		})
	}
	addChart(Config{
		Title:    `combined`,
		Filename: `combined`,
		Series:   []*margaid.Series{getBars, lastKnownPrice, rate},
	})
	addChart(Config{
		Title:    `getBars`,
		Filename: `get-bars`,
		Range:    getBars,
		Series:   []*margaid.Series{getBars},
	})
	addChart(Config{
		Title:    `lastKnownPrice`,
		Filename: `last-known-price`,
		Range:    lastKnownPrice,
		Series:   []*margaid.Series{lastKnownPrice},
	})
	addChart(Config{
		Title:    `rate`,
		Filename: `rate`,
		Range:    rate,
		Series:   []*margaid.Series{rate},
	})

	starts := make(map[*margaid.Series]time.Time)
	addValue := func(series *margaid.Series, t time.Time, value float64) {
		if starts[series] == (time.Time{}) {
			starts[series] = t
		}
		// WARNING: the aggregators depend on the X axis being in seconds
		series.Add(margaid.MakeValue(t.Sub(starts[series]).Seconds(), value))
	}

	var h eventHandlerFuncs
	h.handle = func(event *schema.Event) error {
		d := event.GetApi().GetAccessLog().GetDuration().AsDuration()
		if d < 0 {
			return nil
		}
		var target *margaid.Series
		switch event.GetApi().GetAccessLog().GetData().(type) {
		case *schema.ApiAccessLog_GetBars_:
			target = getBars
		case *schema.ApiAccessLog_LastKnownPrice_:
			target = lastKnownPrice
		case *schema.ApiAccessLog_Rate_:
			target = rate
		default:
			return nil
		}
		t := event.GetTimestamp().AsTime()
		value := d.Seconds() * 1_000
		addValue(target, t, value)
		return nil
	}
	h.flush = func() error { return callFuncs(charts) }
	return &h, nil
}

func (x *Config) avgContentLength() (eventHandler, error) {
	period, err := x.scaleDuration(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	seriesOptions := []margaid.SeriesOption{
		margaid.AggregatedBy(margaid.Avg, period),
	}
	getBars := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`get-bars`),
	)...)
	lastKnownPrice := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`last-known-price`),
	)...)
	rate := margaid.NewSeries(append(append([]margaid.SeriesOption(nil), seriesOptions...),
		margaid.Titled(`rate`),
	)...)

	var charts []func() error
	type Config struct {
		Title    string
		Filename string
		Range    *margaid.Series
		Series   []*margaid.Series
	}
	addChart := func(cfg Config) {
		charts = append(charts, func() error {
			const (
				diagramTitle    = `Avg Content Length`
				diagramFilename = `avg-content-length`
				xAxisTitle      = `Time (seconds)`
				yAxisTitle      = `Size (bytes)`
			)

			if cfg.Range == nil {
				cfg.Range = mergeSeriesMinMax(cfg.Series...)
			}

			diagram := margaid.New(
				diagramWidth,
				calculateHeight(diagramWidth),
				margaid.WithAutorange(margaid.XAxis, cfg.Range),
				margaid.WithAutorange(margaid.YAxis, cfg.Range),
				margaid.WithColorScheme(colourScheme),
				margaid.WithInset(diagramInset),
				margaid.WithPadding(diagramPadding),
			)
			for _, series := range cfg.Series {
				diagram.Line(
					series,
					margaid.UsingAxes(margaid.XAxis, margaid.YAxis),
					margaid.UsingMarker(diagramMarker),
				)
			}
			diagram.Axis(cfg.Range, margaid.XAxis, diagram.ValueTicker('f', 0, 10), false, xAxisTitle)
			diagram.Axis(cfg.Range, margaid.YAxis, diagram.ValueTicker('f', 0, 2), true, yAxisTitle)
			diagram.Title(fmt.Sprintf(`%s (%s) - %s`, diagramTitle, period, cfg.Title))
			if len(cfg.Series) > 1 {
				diagram.Legend(margaid.BottomLeft)
			}

			return x.writeOutputFile(fmt.Sprintf(`%s-%s.svg`, diagramFilename, cfg.Filename), func(w io.Writer) error {
				buf := bufio.NewWriter(w)
				if err := diagram.Render(buf); err != nil {
					return err
				}
				return buf.Flush()
			})
		})
	}
	addChart(Config{
		Title:    `combined`,
		Filename: `combined`,
		Series:   []*margaid.Series{getBars, lastKnownPrice, rate},
	})
	addChart(Config{
		Title:    `getBars`,
		Filename: `get-bars`,
		Range:    getBars,
		Series:   []*margaid.Series{getBars},
	})
	addChart(Config{
		Title:    `lastKnownPrice`,
		Filename: `last-known-price`,
		Range:    lastKnownPrice,
		Series:   []*margaid.Series{lastKnownPrice},
	})
	addChart(Config{
		Title:    `rate`,
		Filename: `rate`,
		Range:    rate,
		Series:   []*margaid.Series{rate},
	})

	starts := make(map[*margaid.Series]time.Time)
	addValue := func(series *margaid.Series, t time.Time, value float64) {
		if starts[series] == (time.Time{}) {
			starts[series] = t
		}
		// WARNING: the aggregators depend on the X axis being in seconds
		series.Add(margaid.MakeValue(t.Sub(starts[series]).Seconds(), value))
	}

	var h eventHandlerFuncs
	h.handle = func(event *schema.Event) error {
		size := event.GetApi().GetAccessLog().GetContentLength()
		if size < 0 {
			return nil
		}
		var target *margaid.Series
		switch event.GetApi().GetAccessLog().GetData().(type) {
		case *schema.ApiAccessLog_GetBars_:
			target = getBars
		case *schema.ApiAccessLog_LastKnownPrice_:
			target = lastKnownPrice
		case *schema.ApiAccessLog_Rate_:
			target = rate
		default:
			return nil
		}
		t := event.GetTimestamp().AsTime()
		value := float64(size)
		addValue(target, t, value)
		return nil
	}
	h.flush = func() error { return callFuncs(charts) }
	return &h, nil
}

func (x *Config) numStatusCodes() (eventHandler, error) {
	period, err := x.scaleDuration(time.Hour)
	if err != nil {
		return nil, err
	}

	// StatusCodes are server-side success / failure splits (ignores 4XX etc)
	type StatusCodes struct {
		Success *margaid.Series
		Failure *margaid.Series
	}
	allSeries := func(v *StatusCodes) []*margaid.Series {
		return []*margaid.Series{
			v.Success,
			v.Failure,
		}
	}
	newSeries := func(name string) *StatusCodes {
		agg := margaid.AggregatedBy(margaid.Sum, period)
		return &StatusCodes{
			Success: margaid.NewSeries(agg, margaid.Titled(name+`-success`)),
			Failure: margaid.NewSeries(agg, margaid.Titled(name+`-failure`)),
		}
	}
	getBars := newSeries(`get-bars`)
	lastKnownPrice := newSeries(`last-known-price`)
	rate := newSeries(`rate`)

	var charts []func() error
	type Config struct {
		Title    string
		Filename string
		Range    *StatusCodes
		Series   []*StatusCodes
	}
	addChart := func(cfg Config) {
		charts = append(charts, func() error {
			const (
				diagramTitle    = `Num Status Codes`
				diagramFilename = `num-status-codes`
				xAxisTitle      = `Time (seconds)`
				yAxisTitle      = `Count (requests)`
			)

			var rnge *margaid.Series
			if cfg.Range != nil {
				rnge = mergeSeriesMinMax(allSeries(cfg.Range)...)
			} else {
				var series []*margaid.Series
				for _, v := range cfg.Series {
					series = append(series, allSeries(v)...)
				}
				rnge = mergeSeriesMinMax(series...)
			}

			diagram := margaid.New(
				diagramWidth,
				calculateHeight(diagramWidth),
				margaid.WithAutorange(margaid.XAxis, rnge),
				margaid.WithAutorange(margaid.YAxis, rnge),
				margaid.WithColorScheme(colourScheme),
				margaid.WithInset(diagramInset),
				margaid.WithPadding(diagramPadding),
			)
			for _, series := range cfg.Series {
				diagram.Bar(
					allSeries(series),
					margaid.UsingAxes(margaid.XAxis, margaid.YAxis),
				)
			}
			diagram.Axis(rnge, margaid.XAxis, diagram.ValueTicker('f', 0, 10), false, xAxisTitle)
			diagram.Axis(rnge, margaid.YAxis, diagram.ValueTicker('f', 0, 2), true, yAxisTitle)
			diagram.Title(fmt.Sprintf(`%s (%s) - %s`, diagramTitle, period, cfg.Title))
			diagram.Legend(margaid.BottomLeft)

			return x.writeOutputFile(fmt.Sprintf(`%s-%s.svg`, diagramFilename, cfg.Filename), func(w io.Writer) error {
				buf := bufio.NewWriter(w)
				if err := diagram.Render(buf); err != nil {
					return err
				}
				return buf.Flush()
			})
		})
	}
	addChart(Config{
		Title:    `combined`,
		Filename: `combined`,
		Series:   []*StatusCodes{getBars, lastKnownPrice, rate},
	})
	addChart(Config{
		Title:    `getBars`,
		Filename: `get-bars`,
		Range:    getBars,
		Series:   []*StatusCodes{getBars},
	})
	addChart(Config{
		Title:    `lastKnownPrice`,
		Filename: `last-known-price`,
		Range:    lastKnownPrice,
		Series:   []*StatusCodes{lastKnownPrice},
	})
	addChart(Config{
		Title:    `rate`,
		Filename: `rate`,
		Range:    rate,
		Series:   []*StatusCodes{rate},
	})

	starts := make(map[*margaid.Series]time.Time)
	addValue := func(series *margaid.Series, t time.Time, value float64) {
		if starts[series] == (time.Time{}) {
			starts[series] = t
		}
		// WARNING: the aggregators depend on the X axis being in seconds
		series.Add(margaid.MakeValue(t.Sub(starts[series]).Seconds(), value))
	}

	var h eventHandlerFuncs
	h.handle = func(event *schema.Event) error {
		var selectSeries func(v *StatusCodes) *margaid.Series
		statusCode := event.GetApi().GetAccessLog().GetStatusCode()
		switch {
		case statusCode >= 200 && statusCode < 300:
			selectSeries = func(v *StatusCodes) *margaid.Series { return v.Success }
		case statusCode >= 500 && statusCode < 600:
			selectSeries = func(v *StatusCodes) *margaid.Series { return v.Failure }
		default:
			return nil
		}
		var target *StatusCodes
		switch event.GetApi().GetAccessLog().GetData().(type) {
		case *schema.ApiAccessLog_GetBars_:
			target = getBars
		case *schema.ApiAccessLog_LastKnownPrice_:
			target = lastKnownPrice
		case *schema.ApiAccessLog_Rate_:
			target = rate
		default:
			return nil
		}
		addValue(selectSeries(target), event.GetTimestamp().AsTime(), 1)
		return nil
	}
	h.flush = func() error { return callFuncs(charts) }
	return &h, nil
}

func (x *Config) scaleDuration(d time.Duration) (time.Duration, error) {
	r := big.NewRat(int64(d), 1)
	r.Mul(r, x.Scale)
	v, err := strconv.ParseInt(r.FloatString(0), 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(v), nil
}

func (x *Config) writeOutputFile(name string, fn func(w io.Writer) error) error {
	if !validFileName(name) {
		return fmt.Errorf(`invalid file name: %s`, name)
	}

	path := filepath.Join(x.OutputDir, name)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	var success bool
	defer func() {
		if !success {
			// attempt to remove the file
			_ = os.Remove(path)
		}
	}()
	defer f.Close()

	if err := fn(f); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	success = true
	return nil
}

func (h *eventHandlerFuncs) Handle(event *schema.Event) error {
	return h.handle(event)
}

func (h *eventHandlerFuncs) Flush() error {
	return h.flush()
}

func validFileName(name string) bool {
	if name == `.` || name == string([]rune{filepath.Separator}) {
		return false
	}
	return filepath.Base(name) == name
}

func parallelOperation[E any](values []E, fn func(v E) error) error {
	if len(values) == 0 {
		return nil
	}

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(values))
	for _, value := range values {
		value := value
		go func() {
			defer wg.Done()
			if err := fn(value); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func callHandlers(handlers []eventHandler, event *schema.Event) error {
	return parallelOperation(handlers, func(h eventHandler) error {
		return h.Handle(event)
	})
}

func flushHandlers(handlers []eventHandler) error {
	return parallelOperation(handlers, func(h eventHandler) error {
		return h.Flush()
	})
}

func callFuncs(funcs []func() error) error {
	return parallelOperation(funcs, func(fn func() error) error { return fn() })
}

func calculateHeight(w int) int {
	return calculateHeightWithRatio(aspectRatioWidth, aspectRatioHeight, w)
}

func calculateHeightWithRatio(wr, hr int, w int) int {
	r := big.NewRat(int64(hr), int64(wr))
	r.Mul(r, big.NewRat(int64(w), 1))
	h, err := strconv.Atoi(r.FloatString(0))
	if err != nil {
		panic(err)
	}
	return h
}

func mergeSeriesMinMax(series ...*margaid.Series) *margaid.Series {
	r := margaid.NewSeries()
	for _, s := range series {
		minX := s.MinX()
		maxX := s.MaxX()
		minY := s.MinY()
		maxY := s.MaxY()
		r.Add(margaid.MakeValue(minX, minY))
		r.Add(margaid.MakeValue(maxX, maxY))
		r.Add(margaid.MakeValue(minX, maxY))
		r.Add(margaid.MakeValue(maxX, minY))
	}
	return r
}
