package server

import (
	"context"
	"fmt"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/dliakhov/bloxroutelabs/client-server-app/server/service"
	"github.com/dliakhov/bloxroutelabs/client-server-app/server/workerpool"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	traceIDKey   = "X-Trace-ID"
	numOfWorkers = 5
)

type App struct {
	config      Configurations
	conn        *amqp.Connection
	itemService service.ItemService
	workerPool  *workerpool.WorkerPool
}

func NewApp(config Configurations, itemService service.ItemService) *App {
	return &App{
		config:      config,
		itemService: itemService,
		workerPool:  workerpool.NewWorkerPool(numOfWorkers),
	}
}

func (a *App) Init() error {
	var err error
	a.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", a.config.RabbitMQConfig.User, a.config.RabbitMQConfig.Password, a.config.RabbitMQConfig.URL))
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Start() error {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		a.config.RabbitMQConfig.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	a.workerPool.Start()

	log.Info("Application is started")
	for d := range msgs {
		a.workerPool.SubmitTask(func() {
			if err := a.ProcessMessage(d); err != nil {
				log.Errorf("Cannot process message: %v", err)
			}
		})
	}

	return nil
}

func (a *App) ProcessMessage(d amqp.Delivery) error {
	var traceID string
	defer func() {
		if err := recover(); err != nil {
			log.WithField(traceIDKey, traceID).Errorf("Panic occured: %v", err)
		}
	}()

	command := new(models.Command)
	err := proto.Unmarshal(d.Body, command)
	if err != nil {
		log.Errorf("Cannot unmarshal message: %v", err)
		return err
	}

	if d.Headers != nil {
		traceIDVal, ok := d.Headers[traceIDKey]
		if !ok {
			log.Warning("Trace id is not set")
		}

		traceID, ok = traceIDVal.(string)
		if !ok {
			log.Warningf("Trace id has no correct type: %v", traceIDVal)
		}
	}

	ctx := context.WithValue(context.Background(), traceIDKey, traceID)

	err = a.itemService.ProcessItemCommand(ctx, command)
	if err != nil {
		log.WithField(traceIDKey, traceID).Errorf("Cannot process message: %v", err)
		return err
	}

	log.Infof("Message ID: %s processed successfully\n", d.MessageId)
	return nil
}

func (a *App) Cleanup() error {
	a.workerPool.Quit()
	return a.conn.Close()
}
