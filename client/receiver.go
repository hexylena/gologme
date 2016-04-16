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
	Receiver     *receiver

	wLogs []*gologme.WindowLogs
	kLogs []*gologme.KeyLogs
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
	l.wLogs = make([]*gologme.WindowLogs, 0)
	l.kLogs = make([]*gologme.KeyLogs, 0)
}

func (l *lgr) Updater(windowLogGranularity int, keyLogGranularity int) {
	go func() {
		// Fetch freshest logs
		c := time.Tick(
			time.Millisecond * time.Duration(windowLogGranularity),
		)
		for _ = range c {
			l.wLogs = append(
				l.wLogs,
				l.WindowLogger.GetFreshestTxtLogs(),
			)
		}
	}()

	go func() {
		// Fetch freshest logs
		c := time.Tick(
			time.Millisecond * time.Duration(keyLogGranularity),
		)
		for _ = range c {
			l.kLogs = append(
				l.kLogs,
				l.KeyLogger.GetFreshestNumLogs(),
			)
		}
	}()
}

func (l *lgr) SendLogs() {
	widx := len(l.wLogs)
	kidx := len(l.kLogs)

	l.Receiver.send(l.wLogs[:widx], l.kLogs[:kidx])

	// This seems like we need a sync on it, but I'm not smart enough for that
	// stuff just yet.
	l.wLogs = l.wLogs[widx:]
	l.kLogs = l.kLogs[kidx:]
}

func Golog(windowLogGranularity int, keyLogGranularity int, standalone bool, serverAddr string) {
	// Setup our receiver
	receiver := &receiver{
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
	exit_chan := make(chan os.Signal, 1)
	signal.Notify(exit_chan, os.Interrupt)
	signal.Notify(exit_chan, syscall.SIGTERM)

	// And push remaining entries to the server
	go func() {
		// In cleanup, we need to
		<-exit_chan
		l.SendLogs()
		os.Exit(0)
	}()

	// Every 10 seconds, send logs
	c := time.Tick(10 * time.Second)
	for _ = range c {
		l.SendLogs()
	}
}
