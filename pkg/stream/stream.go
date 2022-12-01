package stream

import (
	"binance-bot/logger"
	"errors"
	"fmt"
)

var (
	l = logger.GetLogger()
)

var (
	TimeoutError = errors.New("i/o timeout when reading from websocket")
)

type Stream interface {
	FetchData(chan interface{}, []interface{}) chan Result
}

type MapData map[string]interface{}

type Result struct {
	Data  interface{}
	Error error
}

const (
	WebsocketURL = "wss://stream.binance.com:9443/ws/"
)

const (
	AggTradeTag                             = "%s@aggTrade"
	TradeTag                                = "%s@trade"
	KlineTag                                = "%s@kline_%s"
	MiniTickerTag                           = "%s@miniTicker"
	AllMiniTickerTag                        = "!miniTicker@arr"
	IndividualSymbolTickerTag               = "%s@ticker"
	AllTickerTag                            = "!ticker@arr"
	IndividualSymbolRollingWindowStaticsTag = "%s@ticker_%s"
	AllMarketRollingWindowStatisticTag      = "!ticker_%s@arr"
	IndividualSymbolBookTickerTag           = "%s@bookTicker"
	AllBookTickerTag                        = "!bookTicker"
	PartialBookDepthTag                     = "%s@depth%s@%s"
	DiffDepthTag                            = "%s@depth@%s"
)

func ToMap(data []byte) MapData {
	md := make(MapData)

	var parseJson func(data []byte) (MapData, int)
	parseJson = func(data []byte) (MapData, int) {
		m := make(MapData)
		i := 0
		if data[i] == '{' {
			i++
		}
		key := []byte{}
		for i < len(data) {
			if data[i] == '"' || data[i] == '}' {
				i++
				continue
			}
			if data[i] == ':' {
				value := []byte{}
				j := i + 1
				count := 0
				for j < len(data) {
					if data[j] == '{' {
						v, k := parseJson(data[j:])
						m[string(key)] = v
						key = []byte{}
						j += k + 1
						continue
					}
					if data[j] == ',' && count == 0 {
						i = j
						break
					}
					if data[j] == '}' {
						m[string(key)] = string(value)
						return m, j + 1
					}
					if data[j] == '"' {
						j++
						continue
					}
					if data[j] == '[' {
						count++
					}
					if data[j] == ']' {
						count--
					}
					value = append(value, data[j])
					j++
				}
				if string(key) != "" {
					m[string(key)] = string(value)
					key = []byte{}
				}
				i = j + 1
				continue
			}
			key = append(key, data[i])
			i++
		}
		return m, i
	}

	md, _ = parseJson(data)
	return md
}

func wrapUrl(tag string, params ...interface{}) string {
	url := WebsocketURL + tag
	return fmt.Sprintf(url, params...)
}
