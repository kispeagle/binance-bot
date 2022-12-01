package analyzers

import (
	"binance-bot/pkg/indicators"
	repo "binance-bot/pkg/repos"
	"binance-bot/pkg/stream"
	"fmt"
	"time"
)

var (
	streamingPrice stream.Stream
	timeframes     = []string{"1m", "5m", "15m", "1h", "4h", "1d"}
)

type Analysis struct {
	Params                   []interface{}
	MaxPercentageChangedRate float64
	ChangedRate              float64
	CurrentLowPrice          float64
	CurrentPercentage        float64
	// PreviousPercentage       float64
	Ticker *time.Ticker
}

type Options struct {
	strategies map[string]func(stream.Result, chan interface{}) chan Result
}

func initParams(exchangeRate string) map[string]Analysis {
	var analysis map[string]Analysis
	streamingPrice = stream.KlineStream{}

	analysis = make(map[string]Analysis)

	for _, t := range timeframes {
		params := []interface{}{exchangeRate, t}

		a := fetchDataPerTimeframe(params)

		analysis[t] = a
	}
	return analysis
}

func DefaultAnalyzing(exr string, done chan interface{}, r *repo.MongoRepo) chan Result {
	var analysis map[string]Analysis
	analysis = initParams(exr)
	outStream := make(chan Result)

	params1s := []interface{}{exr, "1s"}
	priceStream := streamingPrice.FetchData(done, params1s)

	PercentageChangedRateStream := make(chan stream.KlineModel)
	AnalyzePercentageChangedRate(PercentageChangedRateStream, analysis, outStream)

	go func() {
		for {
			select {
			case <-done:
				return
			case data := <-priceStream:
				if data.Error == stream.TimeoutError {
					l.Error(data.Error)
					return
				}
				d := data.Data.(stream.KlineModel)
				PercentageChangedRateStream <- d

			}
		}
	}()
	return outStream
}

func fetchDataPerTimeframe(params []interface{}) Analysis {
	a := Analysis{}

	a.Params = params
	done := make(chan interface{})
	for {
		data := <-streamingPrice.FetchData(done, params)
		if data.Error == stream.TimeoutError {
			continue
		}
		d := data.Data.(stream.KlineModel)
		a.CurrentLowPrice = d.K.L
		duration := d.K.CloseTime/1000 - time.Now().Unix() + 1
		if duration <= 0 {
			continue
		}
		a.Ticker = time.NewTicker(time.Duration(duration) * time.Second)
		close(done)
		break
	}
	return a
}

func AnalyzePercentageChangedRate(inStream chan stream.KlineModel, analysis map[string]Analysis, outStream chan Result) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Error(r)
			}
		}()

		for {
			select {
			case <-analysis["1m"].Ticker.C:
				analysis["1m"] = fetchDataPerTimeframe(analysis["1m"].Params)
			case <-analysis["5m"].Ticker.C:
				analysis["5m"] = fetchDataPerTimeframe(analysis["5m"].Params)
			case <-analysis["15m"].Ticker.C:
				analysis["15m"] = fetchDataPerTimeframe(analysis["15m"].Params)
			case <-analysis["1h"].Ticker.C:
				analysis["1h"] = fetchDataPerTimeframe(analysis["1h"].Params)
			case <-analysis["4h"].Ticker.C:
				analysis["4h"] = fetchDataPerTimeframe(analysis["4h"].Params)
			case <-analysis["1d"].Ticker.C:
				analysis["1d"] = fetchDataPerTimeframe(analysis["1d"].Params)
			case data := <-inStream:
				for _, t := range timeframes {
					p := indicators.ChangedPercentage(data.K.C, analysis[t].CurrentLowPrice)

					l.Debugf("%s %f", t, p)
					if p > analysis[t].CurrentPercentage+1 {

						a := analysis[t]
						if p > a.MaxPercentageChangedRate {
							a.MaxPercentageChangedRate = p
						}
						a.CurrentPercentage = p
						analysis[t] = a
						msg := fmt.Sprintf("%s percentage Changed rate: %f", data.Symbol, p)
						l.Debug(msg)
						outStream <- Result{Message: msg}

					}
				}
			}
		}
	}()
}
