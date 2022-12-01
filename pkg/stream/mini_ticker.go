package stream

import (
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type MiniTickerModel struct {
	Type      string  `json:"e"`
	Timestamp int64   `json:"E"`
	Symbol    string  `json:"s"`
	C         float64 `json:"c"`
	O         float64 `json:"o"`
	H         float64 `json:"h"`
	L         float64 `json:"l"`
	V         float64 `json:"v"`
	Q         float64 `json:"q"`
}

type MiniTicker struct{}

func (m MiniTicker) FetchData(done chan interface{}, wg *sync.WaitGroup, params []interface{}) chan Result {

	dataStream := make(chan Result, 1000)
	url := wrapUrl(MiniTickerTag, params...)

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			l.Error(err.Error())
			dataStream <- Result{Error: err}
		}
		for {
			select {
			case <-done:
				wg.Done()
				close(dataStream)
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					l.Error(err.Error())
					dataStream <- Result{Error: err}
					close(dataStream)
					return
				}

				data := map2MiniTicker(ToMap(msg))
				dataStream <- Result{Data: data}
			}
		}
	}()
	return dataStream
}

func map2MiniTicker(data MapData) MiniTickerModel {
	var mt MiniTickerModel
	mt.Type = data["e"].(string)
	mt.Timestamp, _ = strconv.ParseInt(data["E"].(string), 10, 64)
	mt.Symbol = data["s"].(string)
	mt.C, _ = strconv.ParseFloat(data["c"].(string), 64)
	mt.H, _ = strconv.ParseFloat(data["h"].(string), 64)
	mt.L, _ = strconv.ParseFloat(data["l"].(string), 64)
	mt.O, _ = strconv.ParseFloat(data["o"].(string), 64)
	mt.Q, _ = strconv.ParseFloat(data["q"].(string), 64)
	mt.V, _ = strconv.ParseFloat(data["v"].(string), 64)

	return mt
}
