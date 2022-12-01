package stream

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type KlineType int

const (
	Kline1s KlineType = iota
	Kline1m
	Kline3m
	Kline15m
	Kline30m
	Kline1h
	Kline2h
	Kline4h
	Kline6h
	Kline8h
	Kline12h
	Kline1d
	Kline3d
	Kline1w
	Kline1M
)

func (t KlineType) String() string {
	names := [...]string{
		"kline1s",
		"kline1m",
		"kline3m",
		"kline15m",
		"kline30m",
		"kline1h",
		"kline2h",
		"kline4h",
		"kline6h",
		"kline8h",
		"kline12h",
		"kline1d",
		"kline3d",
		"kline1w",
		"kline1M",
	}
	return names[t]
}

type KlineModel struct {
	Type      string `json:"e"`
	Timestamp int64  `json:"E"`
	Symbol    string `json:"s"`
	K         Candle `json:"k"`
}

type Candle struct {
	StartTime        int64   `json:"t"`
	CloseTime        int64   `json:"T"`
	Symbol           string  `json:"s"`
	Interval         string  `json:"i"`
	FirstTradeId     string  `json:"f"`
	LastTradeId      string  `json:"L"`
	O                float64 `json:"o"`
	C                float64 `json:"c"`
	H                float64 `json:"h"`
	L                float64 `json:"l"`
	V                float64 `json:"v"`
	NumOfTrade       int64   `json:"n"`
	IsClose          bool    `json:"x"`
	QuoteVolume      string  `json:"q"`
	BaseAssetVolume  string  `json:"V"`
	QuoteAssetVolume string  `json:"Q"`
	Ignore           string  `json:"B"`
}

type KlineStream struct{}

func (k KlineStream) FetchData(done chan interface{}, params []interface{}) chan Result {

	dataStream := make(chan Result, 1000)
	url := wrapUrl(KlineTag, params...)
	strReadingTimeout := os.Getenv("reading_socket_timeout")
	readingTimeout, err := strconv.ParseInt(strReadingTimeout, 10, 64)
	if err != nil {
		l.Error(err)
		return nil
	}

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		var result = Result{}
		if err != nil {
			l.Error(err.Error())
			dataStream <- Result{Error: err}
			return
		}
		for {

			select {
			case <-done:
				close(dataStream)
				return
			default:
				conn.SetReadDeadline(time.Now().Add(time.Duration(readingTimeout) * time.Second))
				_, msg, err := conn.ReadMessage()
				if err != nil {
					l.Debug(url)
					l.Error(err.Error())
					if strings.Contains(err.Error(), "timeout") {
						result.Error = TimeoutError
					} else {
						result.Error = err
					}
					dataStream <- result
					close(dataStream)
					return
				}
				data := map2Kline(ToMap(msg))
				result.Data = data
				dataStream <- result
				result.Data = nil
			}
		}
	}()
	return dataStream
}

func map2Kline(data MapData) KlineModel {

	var k KlineModel
	k.Type = data["e"].(string)
	k.Timestamp, _ = strconv.ParseInt(data["E"].(string), 10, 64)
	k.Symbol = data["s"].(string)
	kline := data["k"].(MapData)
	k.K.StartTime, _ = strconv.ParseInt(kline["t"].(string), 10, 64)
	k.K.CloseTime, _ = strconv.ParseInt(kline["T"].(string), 10, 64)
	k.K.Symbol = kline["s"].(string)
	k.K.Interval = kline["i"].(string)
	k.K.FirstTradeId = kline["f"].(string)
	k.K.FirstTradeId = kline["f"].(string)
	k.K.LastTradeId = kline["L"].(string)
	k.K.H, _ = strconv.ParseFloat(kline["h"].(string), 64)
	k.K.O, _ = strconv.ParseFloat(kline["o"].(string), 64)
	k.K.C, _ = strconv.ParseFloat(kline["c"].(string), 64)
	k.K.L, _ = strconv.ParseFloat(kline["l"].(string), 64)
	k.K.V, _ = strconv.ParseFloat(kline["v"].(string), 64)
	k.K.NumOfTrade, _ = strconv.ParseInt(kline["n"].(string), 10, 64)
	k.K.IsClose, _ = strconv.ParseBool(kline["x"].(string))
	k.K.QuoteAssetVolume = kline["q"].(string)
	k.K.BaseAssetVolume = kline["V"].(string)
	k.K.QuoteAssetVolume = kline["Q"].(string)
	k.K.Ignore = kline["B"].(string)

	return k
}
