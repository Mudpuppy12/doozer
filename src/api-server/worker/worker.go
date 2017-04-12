package main

import (
	"api-server/tasks"
	"log"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"github.com/spf13/viper"
)

var (
	broker        string
	resultBackend string
	exchange      string
	exchangeType  string
	defaultQueue  string
	bindingKey    string

	cnf    config.Config
	server *machinery.Server
	worker *machinery.Worker
)

func init() {
	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	broker = viper.GetString("dozer.broker")
	resultBackend = viper.GetString("dozer.result_backend")
	exchange = viper.GetString("dozer.exchange")
	exchangeType = viper.GetString("dozer.exchange_type")
	defaultQueue = viper.GetString("dozer.default_queue")
	bindingKey = viper.GetString("dozer.binding_key")

	cnf = config.Config{
		Broker:        broker,
		ResultBackend: resultBackend,
		Exchange:      exchange,
		ExchangeType:  exchangeType,
		DefaultQueue:  defaultQueue,
		BindingKey:    bindingKey,
	}

	server, err := machinery.NewServer(&cnf)
	errors.Fail(err, "Could not initialize server")

	// Register tasks
	tasks := map[string]interface{}{
		"add":        exampletasks.Add,
		"multiply":   exampletasks.Multiply,
		"panic_task": exampletasks.PanicTask,
	}
	server.RegisterTasks(tasks)

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker = server.NewWorker("machinery_worker")
}

func main() {
	err := worker.Launch()
	errors.Fail(err, "Could not launch worker")
}
