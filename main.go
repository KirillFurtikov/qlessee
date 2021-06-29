package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	qless "github.com/KirillFurtikov/qlessee/pkg/qless"
	"github.com/caarlos0/env"
	"github.com/jessevdk/go-flags"
)

var (
	qlessClients = map[string]*qless.Client{}
	logger       *log.Logger

	options struct {
		Debug bool `short:"d" long:"debug" description:"Debug: Save log into ./log"`
	}
)

type config struct {
	RedisList []string `env:"QLESS_GO_REDIS_LIST,required" envSeparator:","`
}

func init() {
	parseArgs()
	newLogger()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	for _, value := range cfg.RedisList {
		a := strings.Split(value, "=")
		qlessClient := qless.NewClient(&qless.Options{Name: a[0], URL: a[1]})
		qlessClients[a[0]] = qlessClient
		qlessClient.LoadQueues()
	}
}

func run() int {

	gui := NewGui()
	logger.Println("Gui initialized")

	if err := gui.Start(); err != nil {
		return 1
	}

	return 0

}

func main() {
	os.Exit(run())
}

// NewLogger create logger instance
func newLogger() *log.Logger {
	output := ioutil.Discard
	if options.Debug == true {
		err := error(nil)

		output, err = os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(output, err)
		}
	}
	logger = log.New(output, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func parseArgs() {
	_, err := flags.Parse(&options)
	if err != nil {
		fmt.Printf("%+v\n", options)
		os.Exit(1)
	}
}
