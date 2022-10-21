package main

import (
	"os"

	"github.com/dliakhov/bloxroutelabs/client-server-app/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	cli := cmd.NewCLI()
	if err := cli.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
