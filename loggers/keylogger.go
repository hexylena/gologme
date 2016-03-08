package loggers

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"

	gologme "github.com/erasche/gologme/types"
)

type KeyLogger struct {
	X11DeviceID string

	Cmd       *exec.Cmd
	ProcBytes *bytes.Buffer
	SharedBuf *bytes.Buffer
}

func (logger *KeyLogger) Setup() {
	logger.ProcBytes = &bytes.Buffer{}
	logger.SharedBuf = &bytes.Buffer{}

	logger.Cmd = exec.Command("xinput", "test", logger.X11DeviceID)
	logger.Cmd.Stdout = logger.ProcBytes

	go func() {
		// Start command asynchronously
		_ = logger.Cmd.Start()
		// Only proceed once the process has finished
		logger.Cmd.Wait()
	}()
}

func NewKeyLogger(conf map[string]string) (LogGenerator, error) {
	var devId string
	if val, ok := conf["X11_DEVICE_ID"]; ok {
		devId = val
	} else {
		return &KeyLogger{}, errors.New("X11_DEVICE_ID is required for the key logger")
	}

	return &KeyLogger{
		X11DeviceID: devId,
	}, nil
}

func (logger *KeyLogger) GetFreshestTxtLogs() *gologme.WindowLogs {
	// not implemented
	return nil
}

func (logger *KeyLogger) GetFreshestNumLogs() *gologme.KeyLogs {
	// Find the current length
	expectedBytes := len(logger.ProcBytes.Bytes())
	// Extract those bytes into a new byte array
	newBytes := logger.ProcBytes.Next(expectedBytes)
	// And write those to our "sharedBuf" -- sharedBuf now
	// contains all the latest keylogs
	n, err := logger.SharedBuf.Write(newBytes)
	// Handle posisble errors
	if n != expectedBytes {
		fmt.Printf("Unexpected bytes copied %s vs %s\n", expectedBytes, n)
	}
	if err != nil {
		fmt.Print(err)
	}

	// Number of pressed keys
	keyCount := 0

	// Here we copy data from the latest procBytes into sharedBuf
	// and attempt to process it.
	for {
		// Read as much data as possible
		line, err := logger.SharedBuf.ReadBytes(10)
		// If we hit EOF before '\n', then we get a non-nil error
		if err != nil {
			// So we write back the remaining data to be recylced
			// into the next run
			logger.SharedBuf.Write(line)
			break
		}

		// Now we have a line of data (and there's a '\n' after
		// it), we must
		if bytes.Contains(line, []byte("key press")) {
			keyCount++
		}
	}

	// We're done processing recent data, so we need to send back
	// a message
	return &gologme.KeyLogs{
		Time:  time.Now(),
		Count: keyCount,
	}
}
