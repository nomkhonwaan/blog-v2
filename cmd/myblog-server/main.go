package main

import (
	"fmt"
	"github.com/nomkhonwaan/myblog/cmd/myblog-server/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Version refers to the latest Git tag version.
	Version = "v0.0.1"

	// Revision refers to the latest Git commit hash.
	Revision = "development"
)

func main() {
	cmd := cobra.Command{Version: fmt.Sprintf("%s %s", Version, Revision)}
	cmd.AddCommand(app.Cmd)

	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("server: %s", err)
	}
}
