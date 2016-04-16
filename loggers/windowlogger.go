package loggers

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	gologme "github.com/erasche/gologme/types"
)

type WindowLogger struct {
	X11Connection *xgb.Conn
	lastText      string
}

func (logger *WindowLogger) Setup() {
}

func NewWindowLogger(conf map[string]string) (LogGenerator, error) {
	x, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	return &WindowLogger{
		X11Connection: x,
	}, nil
}

func (logger *WindowLogger) GetFreshestTxtLogs() *gologme.WindowLogs {
	if logger.isScreenSaverRunning() {
		// Locked
		return &gologme.WindowLogs{
			Name: gologme.LOCKED_SCREEN,
			Time: time.Now(),
		}
	} else {
		title, err := logger.getCurWindowTitle()
		if title == logger.lastText {
			return nil
		} else {
			logger.lastText = title
		}

		if err != nil {
			// Ignore errors
			log.Fatal(err)
			return nil
		} else {
			return &gologme.WindowLogs{
				Name: title,
				Time: time.Now(),
			}
		}
	}
}

func (logger *WindowLogger) GetFreshestNumLogs() *gologme.KeyLogs {
	// not implemented
	return nil
}

func (logger *WindowLogger) getCurWindowTitle() (name string, err error) {
	// Get the window id of the root window.
	setup := xproto.Setup(logger.X11Connection)
	root := setup.DefaultScreen(logger.X11Connection).Root

	// Get the atom id (i.e., intern an atom) of "_NET_ACTIVE_WINDOW".
	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(logger.X11Connection, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", err
	}

	// Get the atom id (i.e., intern an atom) of "_NET_WM_NAME".
	aname = "_NET_WM_NAME"
	nameAtom, err := xproto.InternAtom(logger.X11Connection, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", err
	}

	// Get the actual value of _NET_ACTIVE_WINDOW.
	// Note that 'reply.Value' is just a slice of bytes, so we use an
	// XGB helper function, 'Get32', to pull an unsigned 32-bit integer out
	// of the byte slice. We then convert it to an X resource id so it can
	// be used to get the name of the window in the next GetProperty request.
	reply, err := xproto.GetProperty(logger.X11Connection, false, root, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", err
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))

	// Now get the value of _NET_WM_NAME for the active window.
	// Note that this time, we simply convert the resulting byte slice,
	// reply.Value, to a string.
	reply, err = xproto.GetProperty(logger.X11Connection, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", err
	}
	return string(reply.Value), nil
}

func (logger *WindowLogger) isScreenSaverRunning() bool {
	cmd := exec.Command("/usr/bin/xscreensaver-command", "-time")
	stdout, _ := cmd.Output()
	return !strings.Contains(string(stdout), "non-blanked")
}
