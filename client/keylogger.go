package client

// TODO logging

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	gologme "github.com/erasche/gologme/types"
)

func binLogKeys(c chan *gologme.KeyLogs, keyLoggingGranularity int, keyboardDeviceId string) {
	// Figured out watching procs from
	// http://www.darrencoxall.com/golang/executing-commands-in-go/
	cmd := exec.Command("xinput", "test", keyboardDeviceId)
	procBytes := &bytes.Buffer{}
	cmd.Stdout = procBytes

	go func(rb *bytes.Buffer) {
		clock := time.Tick(time.Duration(keyLoggingGranularity) * time.Millisecond)
		sharedBuf := &bytes.Buffer{}
		for {
			<-clock

			// Find the current length
			expectedBytes := len(procBytes.Bytes())
			// Extract those bytes into a new byte array
			newBytes := procBytes.Next(expectedBytes)
			// And write those to our "sharedBuf" -- sharedBuf now
			// contains all the latest keylogs
			n, err := sharedBuf.Write(newBytes)
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
				line, err := sharedBuf.ReadBytes(10)
				// If we hit EOF before '\n', then we get a non-nil error
				if err != nil {
					// So we write back the remaining data to be recylced
					// into the next run
					sharedBuf.Write(line)
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
			c <- &gologme.KeyLogs{
				Time:  time.Now(),
				Count: keyCount,
			}
		}
	}(procBytes)

	// Start command asynchronously
	_ = cmd.Start()

	// Only proceed once the process has finished
	cmd.Wait()
}

func logKeyList(c chan *gologme.KeyLogs) []gologme.KeyLogs {
	slice := make([]gologme.KeyLogs, 0)
	for {
		select {
		case e := <-c:
			slice = append(slice, *e)
		default:
			return slice
		}
	}
}
