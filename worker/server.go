package worker

import (
	"github.com/RichardKnop/machinery/v2"
	amqpbackend "github.com/RichardKnop/machinery/v2/backends/amqp"
	amqpbroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/TencentBlueKing/beego-runtime/conf"
)

func NewServer() (*machinery.Server, error) {
	cnf := conf.MachineryCnf()
	broker := amqpbroker.New(cnf)
	backend := amqpbackend.New(cnf)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	// Register tasks
	tasksMap := map[string]interface{}{
		"HandlePollTask": HandlePollTask,
	}

	return server, server.RegisterTasks(tasksMap)
}
