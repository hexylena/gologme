// Package definition and import the required stdlib packages.
package store

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gologme "github.com/erasche/gologme/types"
)

var UserNotFoundError = errors.New("User not found")
var FailedToConnect = errors.New("Could not connect to database")

type DataStore interface {
	LogToDb(
		uid int,
		windowlogs []*gologme.WindowLogs,
		keylogs []*gologme.KeyLogs,
	)
	CheckAuth(
		user string,
		key string,
	) (int, error)
	Name() string
	FindUserNameById(id int) (string, error)
	SetupDb()
	ExportEventsByDate(tm time.Time) *gologme.EventLog
	MinDate() int
	MaxDate() int
}

type DataStoreFactory func(conf map[string]string) (DataStore, error)

var datastoreFactories = make(map[string]DataStoreFactory)

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
}

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
		for k, _ := range datastoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid Datastore name. Must be one of: %s", strings.Join(availableDatastores, ", ")))
	}

	// Run the factory with the configuration.
	ef, err := engineFactory(conf)
	if err != nil {
		log.Fatal(err)
	}
	ef.SetupDb()
	return ef, err
}
