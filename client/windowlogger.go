package client

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	gologme "github.com/erasche/gologme/types"
	"os/exec"
	"strings"
	"time"
)

func isScreenSaverRunning() bool {
	cmd := exec.Command("/usr/bin/xscreensaver-command", "-time")
	stdout, _ := cmd.Output()
	return !strings.Contains(string(stdout), "non-blanked")
}

func logWindows(c chan *gologme.WindowLogs, windowLogGranularity int) {
	X, err := xgb.NewConn()
	if err != nil {
		//log.Fatal(err)
	}

	var lastTitle string

	ticker := time.Tick(time.Duration(windowLogGranularity) * time.Millisecond)
	for {
		<-ticker

		if isScreenSaverRunning() {
			// Locked
			c <- &gologme.WindowLogs{Name: gologme.LOCKED_SCREEN, Time: time.Now()}
			lastTitle = gologme.LOCKED_SCREEN
		} else {

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
