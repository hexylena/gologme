package client

import (
	"github.com/MarinX/keylogger"
	gologme "github.com/erasche/gologme/util"
	"log"
	"strings"
	"time"
)

func findKeyboard() *keylogger.KeyLogger {
	devs, err := keylogger.NewDevices()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	for i, val := range devs {
		//fmt.Println("Id->", val.Id, "Device->", val.Name, " kb?->", strings.Contains(val.Name, "eyboard"))
		if strings.Contains(val.Name, "eyboard") {
			dev := devs[i]
			return keylogger.NewKeyLogger(dev)
		}
	}
	log.Fatal("Could not find keyboard.")
	// TODO
	return nil
}

func logKeys(c chan *gologme.KeyLogs) {
	rd := findKeyboard()
	if rd == nil {
		log.Fatal("Could not find keyboard")
	}
	//our keyboard..on your system, it will be diffrent
	in, err := rd.Read()
	if err != nil {
		log.Fatal(err)
	}
	for {
		i := <-in
		if i.Value == keylogger.EV_KEY || i.Value == keylogger.EV_REL {
			c <- &gologme.KeyLogs{
				Time:  time.Now(),
				Count: 1,
				//Name: i.KeyString(),
			}
		}
	}
}

func binLogKeys(c chan *gologme.KeyLogs, keyLoggingGranularity int) {
	intermediate := make(chan *gologme.KeyLogs, 1000)
	go logKeys(intermediate)

	clock := time.Tick(time.Duration(keyLoggingGranularity) * time.Millisecond)
	for {
		// Each clock tick
		<-clock
		data := logKeyList(intermediate)
		if len(data) > 0 {
			c <- &gologme.KeyLogs{
				Time:  data[0].Time,
				Count: len(data),
			}
		}
	}
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
