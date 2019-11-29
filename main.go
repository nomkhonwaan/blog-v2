package main

import (
	"github.com/nomkhonwaan/myblog/cmd/myblog"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := myblog.Execute(); err != nil {
		logrus.Fatalf("myblog: %s", err)
	}
}
