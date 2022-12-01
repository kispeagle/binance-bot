package analyzers

import (
	repo "binance-bot/pkg/repos"
	"fmt"
	"runtime/debug"
	"sync"
)

type Analyzers struct {
	coinList    []string
	coinConfigs map[string]Config
	wg          *sync.WaitGroup
	cond        *sync.Cond
	outStream   chan Result
	repo        *repo.Repository
}

func New(
	coinList []string,
	coinConfigs map[string]Config,
	outStream chan Result,
	repo *repo.Repository,
) *Analyzers {

	return &Analyzers{
		coinList:    coinList,
		coinConfigs: coinConfigs,
		outStream:   outStream,
		wg:          &sync.WaitGroup{},
		cond:        sync.NewCond(&sync.Mutex{}),
		repo:        repo,
	}
}

func (a *Analyzers) Run() {

	go func() {
		for _, c := range a.coinList {
			a.wg.Add(1)

			a.analyze(a.coinConfigs[c].name, a.coinConfigs[c].exrs[0])

		}

	}()

}

func (a *Analyzers) ShutDown() {
	a.cond.Broadcast()
}

func (a *Analyzers) analyze(name, exr string) chan interface{} {
	streamingDone := make(chan interface{})
	done := make(chan interface{})

	outAnalyzingStream := DefaultAnalyzing(exr, streamingDone, nil)

	go func() {
		for {
			select {
			case <-done:
				close(streamingDone)
				a.EliminateCoin(name)
				return
			case data := <-outAnalyzingStream:
				// data.Message = fmt.Sprintf("%s %s", name, data.Message)
				data.EXR = exr
				a.outStream <- data
			}
		}
	}()

	return nil

}

func (a *Analyzers) AddCoin(name string, cfg Config) {

	_, ok := a.coinConfigs[name]
	if ok {

		return
	}
	a.coinList = append(a.coinList, name)
	a.coinConfigs[name] = cfg
	a.analyze(name, cfg.exrs[0])
}

func (a *Analyzers) EliminateCoin(name string) (err error) {
	err = a.stop(name)
	if err == AlreadyEliminatedError {
	} else if err != nil {
		l.Error(err.Error())
		return
	}

	delete(a.coinConfigs, name)

	for i, n := range a.coinList {
		if n == name {
			a.coinList[0], a.coinList[i] = a.coinList[i], a.coinList[0]
		}
	}
	a.coinList = a.coinList[1:]
	return
}

func (a *Analyzers) stop(name string) (err error) {
	defer func() {
		if panicInfo := recover(); err != nil {
			errMsg := fmt.Sprintf("%v, %s", panicInfo, string(debug.Stack()))
			l.Error(errMsg)
			err = AlreadyEliminatedError
		}
	}()
	close(a.coinConfigs[name].done)
	return nil
}
