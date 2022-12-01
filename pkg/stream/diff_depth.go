package stream

import (
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type DiffDepthModel struct {
	Type      string `json:"e"`
	Timestamp int64  `json:"E"`
	Symbol    string `json:"s"`
	FirstId   string `json:"U"`
	LastId    string `json:"u"`
	Bid       string `json:"b"`
	Ask       string `json:"a"`
}

type DiffDepthStream struct {
}

func (d DiffDepthStream) FetchData(done chan interface{}, wg *sync.WaitGroup, params []interface{}) chan Result {
	dataStream := make(chan Result)
	url := wrapUrl(DiffDepthTag, params...)

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			l.Error(err.Error())
			dataStream <- Result{Error: err}
			close(dataStream)
			return
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
				data := ToMap(msg)
				dataStream <- Result{Data: data}
			}
		}
	}()
	return dataStream
}

func map2DiffDepth(data MapData) DiffDepthModel {
	var d DiffDepthModel
	d.Type = data["e"].(string)
	d.Timestamp, _ = strconv.ParseInt(data["E"].(string), 10, 64)
	d.Symbol = data["s"].(string)
	d.FirstId = data["f"].(string)
	d.LastId = data["l"].(string)
	d.Ask = data["a"].(string)
	d.Bid = data["b"].(string)
	return d
}
