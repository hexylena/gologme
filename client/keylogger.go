package main

import (
	"fmt"
	"github.com/MarinX/keylogger"
	"github.com/erasche/gologme"
	"log"
	"strings"
	"time"
)

func findKeyboard() *keylogger.KeyLogger {
	devs, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for i, val := range devs {
		fmt.Println("Id->", val.Id, "Device->", val.Name, " kb?->", strings.Contains(val.Name, "eyboard"))
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
				Time: time.Now(),
				Name: i.KeyString(),
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
	return nil
}
