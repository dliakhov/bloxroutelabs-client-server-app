package cmd

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dliakhov/bloxroutelabs/client-server-app/client"
	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/iamolegga/enviper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func clientCmd() *cobra.Command {
	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Client application",
		Run: func(cmd *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)

			configuration, err := getClientConfiguration()
			if err != nil {
				log.Errorf("Cannot read configuration: %v", err)
				return
			}

			startClientApp(configuration)
		},
	}

	return clientCmd
}

func getClientConfiguration() (client.Configurations, error) {
	e := enviper.New(viper.New())

	var pwd string
	var err error
	if pwd, err = os.Getwd(); err != nil {
		log.Fatal("unable to get current working directory: ", err)
	}

	e.AddConfigPath(pwd)
	e.SetConfigName(".config.client")

	// enable viper to handle env values for nested structs
	e.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// defaults to ENV variable values
	e.AutomaticEnv()

	var configuration client.Configurations
	if err := e.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal("Error reading config file: ", err)
		}
	}

	err = e.Unmarshal(&configuration)
	if err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
		return client.Configurations{}, err
	}
	return configuration, nil
}

func startClientApp(configuration client.Configurations) {
	c := client.New(configuration)
	app := client.NewApp(c)

	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)

	commandType, ok := models.CommandType_value[configuration.CommandType]
	if !ok {
		log.Errorf("Command not found: %s", configuration.CommandType)
		return
	}
	go func() {
		err := c.InitClient()
		if err != nil {
			log.Errorf("Cannot initialize client: %v", err)
			return
		}

		err = app.Start(models.CommandType(commandType))
		if err != nil {
			log.Errorf("Error happened for client: %v", err)
		}
	}()

	<-terminate
	log.Info("Terminating application")

	err := c.Cleanup()
	if err != nil {
		log.Errorf("Error happened when cleaning up client: %v", err)
		return
	}
}
