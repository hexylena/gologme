package main

import (
	"fmt"
	gologme "github.com/erasche/gologme/util"
	"log"
	"os"
)

func WriteNotes(notes []gologme.NoteEvent, dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var (
		f            *os.File
		lastUlogTime int64
		init         bool = false
	)
	for _, note := range notes {
		ult := gologme.Ulogme7amTime(note.RealTime)
		if !init || ult != lastUlogTime {
			_ = f.Close()
			var fn string
			if note.Type == gologme.BLOG_TYPE {
				fn = fmt.Sprintf("%s/blog_%d.txt", dir, ult)
			} else {
				fn = fmt.Sprintf("%s/notes_%d.txt", dir, ult)
			}
			fmt.Printf("Writing to %s\n", fn)
			f, err = os.Create(fn)
			if err != nil {
				log.Fatal(err)
			}
			lastUlogTime = ult
			init = true
		}

		if note.Type == gologme.BLOG_TYPE {
			f.WriteString(fmt.Sprintf("%s", note.Contents))
		} else {
			f.WriteString(fmt.Sprintf("%d %s\n", note.RealTime.Unix(), note.Contents))
		}
	}
}

func WriteIntFile(category string, windows []gologme.IEventT, dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var (
		f            *os.File
		lastUlogTime int64
		init         bool = false
	)

	for _, wl := range windows {
		ult := gologme.Ulogme7amTime(wl.RealTime)
		if !init || ult != lastUlogTime {
			_ = f.Close()
			fn := fmt.Sprintf("%s/%s_%d.txt", dir, category, ult)
			fmt.Printf("Writing to %s\n", fn)
			f, err = os.Create(fn)
			if err != nil {
				log.Fatal(err)
			}
			lastUlogTime = ult
			init = true
		}
		f.WriteString(fmt.Sprintf("%d %d\n", wl.RealTime.Unix(), wl.Count))
	}
}

func WriteStringFile(category string, windows []gologme.SEventT, dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var (
		f            *os.File
		lastUlogTime int64
		init         bool = false
	)

	for _, wl := range windows {
		ult := gologme.Ulogme7amTime(wl.RealTime)
		if !init || ult != lastUlogTime {
			_ = f.Close()
			fn := fmt.Sprintf("%s/%s_%d.txt", dir, category, ult)
			fmt.Printf("Writing to %s\n", fn)
			f, err = os.Create(fn)
			if err != nil {
				log.Fatal(err)
			}
			lastUlogTime = ult
			init = true
		}
		f.WriteString(fmt.Sprintf("%d %s\n", wl.RealTime.Unix(), wl.Title))
	}
}
