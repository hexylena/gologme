package store

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gologme "github.com/erasche/gologme/types"
)

// ErrUserNotFound produced when a user cannot be found
var ErrUserNotFound = errors.New("User not found")

// ErrFailedToConnect produced when cannot connect
var ErrFailedToConnect = errors.New("Could not connect to database")

// DataStore interface definition, all DBs must implement this
type DataStore interface {
	LogToDb(
		uid int,
		windowlogs []*gologme.WindowLogs,
		keylogs []*gologme.KeyLogs,
	)
	CreateNote(
		uid int,
		date time.Time,
		message string,
	)
	CreateBlog(
		uid int,
		date time.Time,
		message string,
	)
	CheckAuth(
		user string,
		key string,
	) (int, error)
	Name() string
	FindUserNameByID(id int) (string, error)
	SetupDb()
	ExportWindowLogsByRange(t0, t1 int64) []*gologme.SEvent
	ExportEventsByDate(tm time.Time) *gologme.EventLog
	MinDate() int
	MaxDate() int
}

// DataStoreFactory config
type DataStoreFactory func(conf map[string]string) (DataStore, error)

var datastoreFactories = make(map[string]DataStoreFactory)

// Register a data store factory
func Register(name string, factory DataStoreFactory) {
	if factory == nil {
		log.Panicf("Datastore factory %s does not exist.", name)
	}
	_, registered := datastoreFactories[name]
	if registered {
		log.Fatal("Datastore factory %s already registered. Ignoring.", name)
	}
	datastoreFactories[name] = factory
}

func init() {
	Register("postgres", NewPostgreSQLDataStore)
	Register("sqlite3", NewSqliteSQLDataStore)
	//Register("file", NewFileDataStore)
}

// CreateDataStore creates a new database connection
func CreateDataStore(conf map[string]string) (DataStore, error) {
	// Query configuration for datastore defaulting to "memory".
	var engineName string
	if val, ok := conf["DATASTORE"]; ok {
		engineName = val
	} else {
		engineName = "sqlite3"
	}

	engineFactory, ok := datastoreFactories[engineName]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableDatastores := make([]string, len(datastoreFactories))
		for k := range datastoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, fmt.Errorf(
			"Invalid Datastore name. Must be one of: %s",
			strings.Join(availableDatastores, ", "),
		)
	}

	// Run the factory with the configuration.
	ef, err := engineFactory(conf)
	if err != nil {
		log.Fatal(err)
	}
	ef.SetupDb()
	return ef, err
}
