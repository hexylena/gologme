package loggers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	gologme "github.com/erasche/gologme/types"
)

//var UserNotFoundError = errors.New("User not found")
//var FailedToConnect = errors.New("Could not connect to database")

type LogGenerator interface {
	Setup()
	NumLog(
		chan *gologme.KeyLogs,
	)
	TextLog(
		c chan *gologme.WindowLogs,
	)
}

type LogGeneratingFactory func(conf map[string]string) (LogGenerator, error)

var logGeneratingFactories = make(map[string]LogGeneratingFactory)

func Register(name string, factory LogGeneratingFactory) {
	if factory == nil {
		log.Panicf("Log generating factory %s does not exist.", name)
	}
	_, registered := logGeneratingFactories[name]
	if registered {
		log.Fatal("Log generating factory %s already registered. Ignoring.", name)
	}
	logGeneratingFactories[name] = factory
}

func init() {
	Register("keys", NewKeyLogger)
	Register("windows", NewWindowLogger)
}

//loggers/logger.go:41: cannot use NewKeyLogger    (type func(map[string]string) (*KeyLogger, error))   as type LogGeneratingFactory in argument to Register
//loggers/logger.go:42: cannot use NewWindowLogger (type func(map[string]string) (WindowLogger, error)) as type LogGeneratingFactory in argument to Register

func AvailableLoggers() []string {
	availableLogGenerators := make([]string, len(logGeneratingFactories))
	for k, _ := range logGeneratingFactories {
		availableLogGenerators := append(availableLogGenerators, k)
	}
	return availableLogGenerators
}

func CreateLogGenerator(conf map[string]string) (LogGenerator, error) {
	// Query configuration for datastore defaulting to "memory".
	var engineName string
	if val, ok := conf["LOGGER"]; ok {
		engineName = val
	} else {
		engineName = "keys"
	}

	//loggers/logger.go:62: type LogGeneratingFactory is not an expression
	logFactory, ok := logGeneratingFactories[engineName]
	if !ok {
		availableLogGenerators := AvailableLoggers()
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		return nil, errors.New(fmt.Sprintf("Invalid log generator name. Must be one of: %s", strings.Join(availableLogGenerators, ", ")))
	}

	// Run the factory with the configuration.
	ef, err := logFactory(conf)
	if err != nil {
		log.Fatal(err)
	}
	ef.Setup()
	return ef, err
}
