package main

import (
	"fmt"
	"github.com/nomkhonwaan/myblog/cmd/serve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// version refers to the latest Git tag
	version = "v0.0.1"

	// revision refers to the latest Git commit hash
	revision = "development"
)

func main() {
	cmd := &cobra.Command{
		Version: fmt.Sprintf("%s %s\n", version, revision),
	}
	cmd.AddCommand(serve.Cmd)

	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("myblog: %s", err)
	}
}
