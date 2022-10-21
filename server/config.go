package server

type Configurations struct {
	RabbitMQConfig RabbitMQConfig
}

type RabbitMQConfig struct {
	URL       string
	User      string
	Password  string
	QueueName string
}
