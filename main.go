package main

import (
	"binance-bot/logger"
	"binance-bot/pkg/analyzers"
	binancebot "binance-bot/pkg/binance_bot"
	"binance-bot/pkg/messenger"
	"binance-bot/pkg/stream"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var l *zap.SugaredLogger

func init() {
	l = logger.GetLogger()
}
func main() {

	// testAnalyzer(wg)
	// go testPriceStream(wg)
	terminated := testBinanceBot()
	// terminated := testPriceStreaming()

	// wg.Wait()
	// analyzers.GetAllPairs()
	<-terminated
}

func testBinanceBot() chan interface{} {

	// Create coin config
	name := "omg"
	exrs := []string{"omgusdt"}
	conf := analyzers.NewConfig(name, exrs)

	name2 := "eth"
	exrs2 := []string{"ethusdt", "ethbusd"}
	conf2 := analyzers.NewConfig(name2, exrs2)

	name3 := "btc"
	exrs3 := []string{"btcusdt", "btcbusd"}
	conf3 := analyzers.NewConfig(name3, exrs3)

	coinConfigs := make(map[string]analyzers.Config)
	coinConfigs[name] = conf
	coinConfigs[name2] = conf2
	coinConfigs[name3] = conf3
	coinList := []string{"btc", "omg", "eth"}

	outAnalyzerStream := make(chan analyzers.Result)
	a := analyzers.New(coinList, coinConfigs, outAnalyzerStream, nil)
	a.Run()

	// Initiate messenger boot
	api_token := os.Getenv("telegram_api_token")
	tlgbot := messenger.NewTelegramBot(api_token)

	// initiate binance bot
	bot := binancebot.NewBinanceBot(outAnalyzerStream, tlgbot)

	// subscribe to binance bot
	cfg := make(binancebot.Config)
	cfg["roomId"] = int64(5758976328)
	bot.Subscribe(api_token, cfg)

	terminated := bot.Run()

	return terminated
}

func testAnalyzer() chan interface{} {

	terminated := make(chan interface{})
	name := "omg"
	exrs := []string{"omgusdt"}
	conf := analyzers.NewConfig(name, exrs)

	name2 := "eth"
	exrs2 := []string{"ethusdt", "ethbusd"}
	conf2 := analyzers.NewConfig(name2, exrs2)

	name3 := "btc"
	exrs3 := []string{"btcusdt", "btcbusd"}
	conf3 := analyzers.NewConfig(name3, exrs3)

	coinConfigs := make(map[string]analyzers.Config)
	coinConfigs[name] = conf
	coinConfigs[name2] = conf2
	coinConfigs[name3] = conf3
	coinList := []string{"btc", "omg", "eth"}
	// coinList := []string{"btc", "eth"}

	outAnalyzerStream := make(chan analyzers.Result)
	a := analyzers.New(coinList, coinConfigs, outAnalyzerStream, nil)
	a.Run()

	go func() {
		for {
			data := <-outAnalyzerStream
			fmt.Println(data)
		}
	}()

	// token_api := os.Getenv("telegram_api_token")
	// bot := messenger.NewTelegramBot(token_api)
	// msg := make(messenger.Message)
	// msg["roomId"] = int64(5758976328)
	// for data := range outAnalyzerStream {

	// 	msg["text"] = data.Message
	// 	bot.Send(msg)
	// }
	return terminated
}

func testPriceStreaming() chan interface{} {

	_, baseList, _ := analyzers.GetCoinList()
	fmt.Println(len(baseList))

	terminated := make(chan interface{})
	var streaming stream.Stream
	streaming = stream.KlineStream{}

	params := []interface{}{"icxusdt", "1s"}
	priceStream := streaming.FetchData(terminated, params)

	for data := range priceStream {
		fmt.Println(data)
	}

	return terminated
}
