package binancebot

import (
	"binance-bot/logger"
	"binance-bot/pkg/analyzers"
	"binance-bot/pkg/messenger"
	"fmt"
)

var (
	l = logger.GetLogger()
)

type Config map[string]interface{}

// func NewConfig(key string, value interface{}) Config{
// 	return
// }

type BinanceBot struct {
	publisher   chan analyzers.Result
	messengers  []messenger.MessengerBot
	subscribers map[string][]Config
	//  subscribers map[chan messenger.Message] map[string]interface{}
}

func NewBinanceBot(pub chan analyzers.Result, subs ...messenger.MessengerBot) *BinanceBot {
	subscribers := make(map[string][]Config)
	return &BinanceBot{publisher: pub, messengers: subs, subscribers: subscribers}
}

func (b *BinanceBot) Run() chan interface{} {
	terminated := make(chan interface{})

	outs := make([]chan messenger.Message, 0)
	for _, mb := range b.messengers {
		o, _ := mb.Listen(nil)
		outs = append(outs, o)
	}

	go func() {
		for {
			select {
			case data := <-b.publisher:
				for _, s := range b.subscribers[b.messengers[0].GetTokenApi()] {
					go func() {
						msg := make(messenger.Message)
						msg["roomId"] = s["roomId"].(int64)
						msg["text"] = data.Message
						err := b.messengers[0].Send(msg)
						if err != nil {
							l.Error(err)
						}
					}()
				}
			case command := <-outs[0]:
				fmt.Println(command)

			}
		}
	}()
	return terminated
}

func (b *BinanceBot) Subscribe(apiToken string, cfg Config) {
	b.subscribers[apiToken] = append(b.subscribers[apiToken], cfg)
}

// func (b *BinanceBot) AddMessenger(m chan messenger.Message) {
// 	b.subscribers = append(b.subscribers, m)
// }

// func (b *BinanceBot) listen() {

// 	for {
// 		for
// 	}
// }
