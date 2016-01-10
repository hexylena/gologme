// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/erasche/gologme"
)

var logbuffer int = 4

func main() {
	window_titles := make(chan *gologme.WindowLogs)
	keypresses := make(chan *gologme.KeyLogs, 1000)

	// Start logging
	go logWindows(window_titles)
	go logKeys(keypresses)

	wl := make([]gologme.WindowLogs, logbuffer)
	wi := 0
	first := true

	// Trap the exit signal
	exit_chan := make(chan os.Signal, 1)
	signal.Notify(exit_chan, os.Interrupt)
	signal.Notify(exit_chan, syscall.SIGTERM)

	// And send remaining entries to the server
	go func() {
		// In cleanup, we need to
		<-exit_chan
		send(wl, wi, logKeyList(keypresses))
		os.Exit(1)
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
			send(wl, len(wl), logKeyList(keypresses))
		}

		//Stick in next log position
		wl[wi] = *(<-window_titles)
		wi++
	}
}
