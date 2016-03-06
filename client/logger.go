package client

import (
	"os"
	"os/signal"
	"syscall"

	gologme "github.com/erasche/gologme/types"
)

func Golog(logbuffer int, windowLogGranularity int, keyLogGranularity int, standalone bool, serverAddr string) {
	window_titles := make(chan *gologme.WindowLogs)
	keypresses := make(chan *gologme.KeyLogs, 1000)
	receiver := &receiver{
		ServerAddress: serverAddr,
	}

	// Start logging
	go logWindows(window_titles, windowLogGranularity)
	go binLogKeys(keypresses, keyLogGranularity, "11")

	wl := make([]gologme.WindowLogs, logbuffer)
	wi := 0
	first := true

	// Trap the exit signal
	exit_chan := make(chan os.Signal, 1)
	signal.Notify(exit_chan, os.Interrupt)
	signal.Notify(exit_chan, syscall.SIGTERM)

	// And push remaining entries to the server
	go func() {
		// In cleanup, we need to
		<-exit_chan
		kl := logKeyList(keypresses)
		receiver.send(wl, wi, kl)
		os.Exit(0)
	}()

	// Until then, loop
	for {
		// If we've hit our buffer, reset
		if wi >= logbuffer {
			wi = 0
			first = false
		}
		// and send
		if wi == 0 && !first {
			kl := logKeyList(keypresses)
			receiver.send(wl, len(wl), kl)
		}

		//Stick in next log position
		wl[wi] = *(<-window_titles)
		wi++
	}
}
