package cmd

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dliakhov/bloxroutelabs/client-server-app/server"
	"github.com/dliakhov/bloxroutelabs/client-server-app/server/repository"
	"github.com/dliakhov/bloxroutelabs/client-server-app/server/service"
	"github.com/iamolegga/enviper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Server application",
		Run: func(cmd *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)

			configuration, err := getServerConfiguration()
			if err != nil {
				log.Errorf("Cannot read configuration: %v", err)
				return
			}

			startServerApp(configuration, err)
		},
	}

}

func getServerConfiguration() (server.Configurations, error) {
	e := enviper.New(viper.New())

	var pwd string
	var err error
	if pwd, err = os.Getwd(); err != nil {
		log.Fatal("unable to get current working directory: ", err)
	}

	e.AddConfigPath(pwd)
	e.SetConfigName(".config.server")

	// enable viper to handle env values for nested structs
	e.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// defaults to ENV variable values
	e.AutomaticEnv()

	var configuration server.Configurations
	if err := e.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal("Error reading config file: ", err)
		}
	}

	err = e.Unmarshal(&configuration)
	if err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
		return server.Configurations{}, err
	}
	return configuration, nil
}

func startServerApp(configuration server.Configurations, err error) {
	repo := repository.New()
	itemService := service.New(repo)

	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)

	app := server.NewApp(configuration, itemService)

	go func() {
		err = app.Init()
		if err != nil {
			log.Errorf("Cannot init server app: %v", err)
			return
		}

		err := app.Start()
		if err != nil {
			log.Errorf("Cannot start server app: %v", err)
			return
		}
	}()

	<-terminate

	log.Info("Terminating application")
	err = app.Cleanup()
	if err != nil {
		log.Errorf("Cannot clean up server app: %v", err)
		return
	}
}
