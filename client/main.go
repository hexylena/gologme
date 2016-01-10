// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erasche/gologme"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

var logbuffer int = 32

func getCurWindowTitle(X *xgb.Conn) (name string, err error) {
	// Get the window id of the root window.
	setup := xproto.Setup(X)
	root := setup.DefaultScreen(X).Root

	// Get the atom id (i.e., intern an atom) of "_NET_ACTIVE_WINDOW".
	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(X, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", err
	}

	// Get the atom id (i.e., intern an atom) of "_NET_WM_NAME".
	aname = "_NET_WM_NAME"
	nameAtom, err := xproto.InternAtom(X, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", err
	}

	// Get the actual value of _NET_ACTIVE_WINDOW.
	// Note that 'reply.Value' is just a slice of bytes, so we use an
	// XGB helper function, 'Get32', to pull an unsigned 32-bit integer out
	// of the byte slice. We then convert it to an X resource id so it can
	// be used to get the name of the window in the next GetProperty request.
	reply, err := xproto.GetProperty(X, false, root, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", err
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))

	// Now get the value of _NET_WM_NAME for the active window.
	// Note that this time, we simply convert the resulting byte slice,
	// reply.Value, to a string.
	reply, err = xproto.GetProperty(X, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", err
	}
	return string(reply.Value), nil
}

func logWindows(c chan *gologme.WindowLogs) {
	X, err := xgb.NewConn()
	if err != nil {
		//log.Fatal(err)
	}

	var lastTitle string
	for {
		title, err := getCurWindowTitle(X)
		if err != nil {
			// Ignore errors
			//log.Fatal(err)
		} else {
			if title != lastTitle {
				c <- &gologme.WindowLogs{Name: title, Time: time.Now()}
				lastTitle = title
			}
		}
	}

}

func send(wl []gologme.WindowLogs, wi int) {
	client, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatal("Error in dialing", err)
	}
	args := &gologme.RpcArgs{
		User:    0,
		Windows: wl,
		Length:  wi,
	}
	var result gologme.Result
	err = client.Call("Golog.Log", args, &result)
	if err != nil {
		log.Fatal("Error calling RPC method", err)
	}
}

func main() {
	window_titles := make(chan *gologme.WindowLogs)

	// Start logging
	go logWindows(window_titles)
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
		send(wl, wi)
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
			send(wl, len(wl))
		}

		//Stick in next log position
		wl[wi] = *(<-window_titles)
		wi++
	}
}
