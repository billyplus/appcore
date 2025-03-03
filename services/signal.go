package services

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type SignalService chan os.Signal

func NewSignalService() SignalService {
	ch := make(SignalService, 1)
	return ch
}

func (ss SignalService) Init() error {
	switch runtime.GOOS {
	case `windows`:
		signal.Notify(ss, syscall.SIGINT, syscall.SIGTERM)
	default:
		signal.Notify(ss, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	}

	return nil
}

func (ss SignalService) Name() string {
	return "SignalService"
}

func (ss SignalService) Start() error {
	sig := <-ss
	if sig != nil {
		fmt.Printf("退出程序 signal=%v\n", sig)
	}
	signal.Stop(ss)
	return nil
}

func (ss SignalService) Shutdown() {
	close(ss)
}
