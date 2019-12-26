package main

import (
	"fmt"
	generate_indices "github.com/nomkhonwaan/myblog/cmd/generate-indices"
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
		Short:   "Personal blog website written in Go with Angular 2+",
		Version: fmt.Sprintf("%s %s\n", version, revision),
	}
	cmd.AddCommand(generate_indices.Command, serve.Command)

	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("myblog: %s", err)
	}
}
