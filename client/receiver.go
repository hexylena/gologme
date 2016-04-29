package client

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erasche/gologme/loggers"
	gologme "github.com/erasche/gologme/types"
)

type lgr struct {
	KeyLogger    loggers.LogGenerator
	WindowLogger loggers.LogGenerator
	Receiver     *Receiver

	winLogs []*gologme.WindowLogs
	keyLogs []*gologme.KeyLogs
}

func (l *lgr) setupLoggers() {
	keyLogger, err := loggers.CreateLogGenerator(map[string]string{
		"LOGGER":        "keys",
		"X11_DEVICE_ID": "11",
	})
	// Must have a key Logger available. Maybe no panic later
	if err != nil {
		panic(err)
	}

	windowLogger, err := loggers.CreateLogGenerator(map[string]string{
		"LOGGER": "windows",
	})
	// Must have a window Logger available. Maybe no panic later
	if err != nil {
		panic(err)
	}

	l.KeyLogger = keyLogger
	l.WindowLogger = windowLogger
	l.winLogs = make([]*gologme.WindowLogs, 0)
	l.keyLogs = make([]*gologme.KeyLogs, 0)
}

func (l *lgr) Updater(windowLogGranularity int, keyLogGranularity int) {
	go func() {
		// Fetch freshest logs
		c := time.Tick(
			time.Millisecond * time.Duration(windowLogGranularity),
		)
		for _ = range c {
			newWLogs := l.WindowLogger.GetFreshestTxtLogs()
			if newWLogs != nil {
				l.winLogs = append(
					l.winLogs,
					newWLogs,
				)
			}
		}
	}()

	go func() {
		// Fetch freshest logs
		c := time.Tick(
			time.Millisecond * time.Duration(keyLogGranularity),
		)
		for _ = range c {
			newKLogs := l.KeyLogger.GetFreshestNumLogs()
			if newKLogs != nil {
				l.keyLogs = append(
					l.keyLogs,
					newKLogs,
				)
			}
		}
	}()
}

func (l *lgr) SendLogs() {
	widx := len(l.winLogs)
	kidx := len(l.keyLogs)

	l.Receiver.Send(l.winLogs[:widx], l.keyLogs[:kidx])

	// This seems like we need a sync on it, but I'm not smart enough for that
	// stuff just yet.
	l.winLogs = l.winLogs[widx:]
	l.keyLogs = l.keyLogs[kidx:]
}

// Golog function which starts local loggers
func Golog(windowLogGranularity int, keyLogGranularity int, standalone bool, serverAddr string) {
	// Setup our receiver
	receiver := &Receiver{
		ServerAddress: serverAddr,
	}

	// Setup our loggers
	l := &lgr{
		Receiver: receiver,
	}
	l.setupLoggers()

	// Trigger our updater in the background
	l.Updater(windowLogGranularity, keyLogGranularity)

	// Trap the exit signal
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	// And push remaining entries to the server
	go func() {
		// In cleanup, we need to
		<-exitChan
		l.SendLogs()
		os.Exit(0)
	}()

	// Every 10 seconds, send logs
	c := time.Tick(10 * time.Second)
	for _ = range c {
		l.SendLogs()
	}
}
