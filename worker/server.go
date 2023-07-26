package worker

import (
	"github.com/RichardKnop/machinery/v2"
	nullbackend "github.com/RichardKnop/machinery/v2/backends/null"
	amqpbroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/TencentBlueKing/beego-runtime/conf"
)

var server *machinery.Server

func NewServer() (*machinery.Server, error) {

	if server != nil {
		return server, nil
	}

	cnf := conf.MachineryCnf()
	broker := amqpbroker.New(cnf)
	backend := nullbackend.New()
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	// Register tasks
	tasksMap := map[string]interface{}{
		"HandlePollTask": HandlePollTask,
	}

	return server, server.RegisterTasks(tasksMap)
}
