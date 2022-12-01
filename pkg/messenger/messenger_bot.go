package messenger

import "binance-bot/logger"

type Message map[string]interface{}

var l = logger.GetLogger()

type MessengerBot interface {
	Send(Message) error
	Listen(chan interface{}) (chan Message, chan interface{})
	GetTokenApi() string
}
