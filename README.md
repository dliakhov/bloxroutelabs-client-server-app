# Client server application

This repo contains two applications:
* server
* client

Client infinitely sends messages to RabbitMQ queue with command which can be one of this:
* AddItem
* RemoveItem
* GetItem
* GetAllItems

Command type can be specified via environment variable `COMMANDTYPE` or in the config file `.config.client.yaml`. Payload for which command client generate randomly.

Server polls messages from the items and process them. It stores items in the memory (LinkedHashMap structure). Server can process several items simultaneously.

## Prerequisites

You have to have installed:
* Docker
* make

## How to run

To run all services you can run command:
> make run-demo-docker-compose

It will start locally:
* RabbitMQ
* Server application
* 4 client applications
  * one will send AddItem command
  * second will send GetItem command
  * third will send RemoveItem command
  * fourth will send GetAllItems command

Clients will send messages to queue in random period of time to simulate work.

Examples of client and server  configuration you can find in the next files accordingly:
* .config.default.client.yaml
* .config.default.server.yaml

### Run client locally
To run locally client without docker you have to:
* run `make run-rabbit-mq` to start RabbitMQ instance
* copy `.config.default.client.yaml` to `.config.client.yaml`
* update `.config.client.yaml` (rabbit config data or command type)
* run `make run-client`

### Run server locally
To run locally server without docker you have to:
* run `make run-rabbit-mq` to start RabbitMQ instance if it's not started. Skip this step if you started it already
* copy `.config.default.server.yaml` to `.config.server.yaml`
* update `.config.server.yaml` (rabbit config data or command type)
* run `make run-server`

### Run tests

To run tests run command:

> make run-tests
