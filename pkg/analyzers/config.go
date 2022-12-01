package analyzers

import (
	"binance-bot/logger"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	l                      = logger.GetLogger()
	AlreadyEliminatedError = errors.New("Streaming Coin was eliminated")
)

type Result struct {
	EXR     string
	Message string
	Error   error
}

type Config struct {
	name string
	exrs []string
	done chan interface{}
}

func NewConfig(name string, exrs []string) Config {
	done := make(chan interface{})
	return Config{
		name: name,
		exrs: exrs,
		done: done,
	}
}

func GetCoinList() ([]string, []string, map[string][]string) {
	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)
	if err != nil {
		l.Error(err.Error())
		return nil, nil, nil
	}

	pairs, _ := ioutil.ReadAll(resp.Body)

	pairSymbolPattern, _ := regexp.Compile("{\"symbol\":\"([a-zA-Z]+)\",")
	quoteAssetPattern, _ := regexp.Compile("\"quoteAsset\":\"([a-zA-Z]+)\"")
	baseAssetPattern, _ := regexp.Compile("\"baseAsset\":\"([a-zA-Z]+)\"")

	pairCoinSet := make(map[string][]string)
	quoteList := []string{}
	baseList := []string{}
	baseSet := make(map[string]struct{})
	quoteSet := make(map[string]struct{})

	quoteResults := quoteAssetPattern.FindAllSubmatch(pairs, -1)
	baseResults := baseAssetPattern.FindAllSubmatch(pairs, -1)
	pairResults := pairSymbolPattern.FindAllSubmatch(pairs, -1)

	for _, e := range quoteResults {
		quoteSet[string(e[1])] = struct{}{}
	}
	for key, _ := range quoteSet {
		quoteList = append(quoteList, key)
	}

	for _, e := range baseResults {
		key := string(e[1])
		baseSet[key] = struct{}{}
		ps := []string{}
		for _, ee := range pairResults {
			pattern := key + ".+"
			p := string(ee[1])
			if ok, _ := regexp.MatchString(pattern, p); ok {
				for _, q := range quoteList {
					if p[len(key):] == q {
						ps = append(ps, string(ee[1]))
					}

				}
			}
		}
		pairCoinSet[key] = ps
		ps = []string{}
	}

	for key, _ := range baseSet {
		baseList = append(baseList, key)
	}

	return quoteList, baseList, pairCoinSet
}
