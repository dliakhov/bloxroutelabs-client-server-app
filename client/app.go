package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const traceIDKey = "X-Trace-ID"

type App struct {
	client *Client
}

func NewApp(client *Client) *App {
	return &App{
		client: client,
	}
}

func (a *App) Start(commandType models.CommandType) error {
	rand.Seed(time.Now().Unix())

	for {
		traceID := uuid.New().String()
		ctx := context.WithValue(context.Background(), traceIDKey, traceID)

		command, err := createRandomCommand(commandType)
		if err != nil {
			return err
		}

		err = a.client.SendCommand(ctx, command)
		if err != nil {
			log.WithField(traceIDKey, ctx.Value(traceIDKey)).Errorf("Fail send command: %v", err)
		}

		// wait some time to send command again
		secondsWait := rand.Int63n(10)
		time.Sleep(time.Duration(secondsWait) * time.Second)
	}
}

func createRandomCommand(commandType models.CommandType) (*models.Command, error) {
	switch commandType {
	case models.CommandType_AddItem:
		itemId := rand.Int63()
		randomChar := rand.Intn(25) + 65 // random character from A to Z

		return &models.Command{
			Type:        commandType,
			ItemID:      itemId,
			ItemPayload: fmt.Sprintf("%c", randomChar),
		}, nil
	case models.CommandType_GetItem:
		itemId := rand.Int63()

		return &models.Command{
			Type:   commandType,
			ItemID: itemId,
		}, nil
	case models.CommandType_RemoveItem:
		itemId := rand.Int63()

		return &models.Command{
			Type:   commandType,
			ItemID: itemId,
		}, nil
	case models.CommandType_GetAllItems:
		return &models.Command{
			Type: commandType,
		}, nil
	default:
		return nil, errors.New("Command type is unknown")
	}
}
