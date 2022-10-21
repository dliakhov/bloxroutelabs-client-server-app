package client

type Configurations struct {
	RabbitMQConfig RabbitMQConfig
	CommandType    string
}

type RabbitMQConfig struct {
	URL       string
	User      string
	Password  string
	QueueName string
}
