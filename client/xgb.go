package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/erasche/gologme"
	"time"
)

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
