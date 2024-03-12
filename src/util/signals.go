package util

import (
	"os"
	"os/signal"
)

func HandleSignals(signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		GetLogger("signal").Debug("received %s signal", sig)
		//time.Sleep(time.Second)
		done <- true
	}()
	<-done
}
