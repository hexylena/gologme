package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	gologme "github.com/erasche/gologme/types"
)

// FileDataStore struct
type FileDataStore struct {
	DSN      string
	KeyFreqs *os.File
	Windows  *os.File
	Blog     *os.File
	Notes    *os.File
}

// SetupDb runs any migrations needed
func (ds *FileDataStore) SetupDb() {
	return
}

// LogToDb logs a set of windowLogs and keyLogs to the DB
func (ds *FileDataStore) LogToDb(uid int, windowlogs []*gologme.WindowLogs, keylogs []*gologme.KeyLogs) {

	wll := len(windowlogs)
	log.Printf("%d window logs %d key logs from [%d]\n", wll, len(keylogs), uid)

	for _, w := range keylogs {
		keyfreqText := fmt.Sprintf("%d %d", w.Time.Unix(), w.Count)
		if _, err := ds.KeyFreqs.WriteString(keyfreqText); err != nil {
			log.Fatal(err)
		}
	}

	for i, w := range windowlogs {
		windowText := fmt.Sprintf("%d %s", w.Time.Unix(), w.Name)
		if _, err := ds.Windows.WriteString(windowText); err != nil {
			log.Fatal(err)
		}

		// Dunno why windowLogs comes through two too big, so whatever.
		if i >= wll-1 {
			break
		}
	}
}

// CheckAuth of the user+key, returning -1 or the user's ID
func (ds *FileDataStore) CheckAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	return 0, nil
}

// Name of the DS implementatino
func (ds *FileDataStore) Name() string {
	return "FileDataStore"
}

// MaxDate returns latest log entry
func (ds *FileDataStore) MaxDate() int {
	// not implemented
	return 0
}

// MinDate returns first log entry
func (ds *FileDataStore) MinDate() int {
	// not implemented
	return 0
}

// FindUserNameByID returns a username for a given ID
func (ds *FileDataStore) FindUserNameByID(id int) (string, error) {
	return "admin", nil
}

func (ds *FileDataStore) ExportWindowLogsByRange(t0 int64, t1 int64) []*gologme.SEvent {
	return nil
}

func (ds *FileDataStore) exportKeyLogsByRange(t0 int64, t1 int64) []*gologme.IEvent {
	return nil
}

func (ds *FileDataStore) exportBlog(t0 int64, t1 int64) []*gologme.SEvent {
	return nil
}

func (ds *FileDataStore) exportNotes(t0 int64, t1 int64) []*gologme.SEvent {
	return nil
}

// ExportEventsByDate extracts events for a given day
func (ds *FileDataStore) ExportEventsByDate(tm time.Time) *gologme.EventLog {
	return nil
}

// NewFileDataStore builds a new sqlite3 DS
func NewFileDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_URL"]; ok {
		dsn = val
	} else {
		return nil, errors.New(fmt.Sprintf("%s is required for the sqlite datastore", "DATASTORE_URL"))
	}

	keyFilePath := path.Join(dsn, "keyfreq_date.txt")
	keyFileHandle, err := os.OpenFile(keyFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
	}
	winFilePath := path.Join(dsn, "window_date.txt")
	winFileHandle, err := os.OpenFile(winFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
	}
	blogFilePath := path.Join(dsn, "blog_date.txt")
	blogFileHandle, err := os.OpenFile(blogFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
	}
	notesFilePath := path.Join(dsn, "notes_date.txt")
	notesFileHandle, err := os.OpenFile(notesFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
	}

	return &FileDataStore{
		DSN:      dsn,
		KeyFreqs: keyFileHandle,
		Windows:  winFileHandle,
		Blog:     blogFileHandle,
		Notes:    notesFileHandle,
	}, nil
}
