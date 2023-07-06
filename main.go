package main

import (
	"flag"
	"fmt"

	"github.com/jass-walia/bank_ops/config"
	"github.com/golang/glog"
)

// configPath path to config file
var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./.env", "path to config file")
}

func main() {
	// Parse flag arguments.
	flag.Parse()

	// initalize config
	config.Initialize(configPath)

	fmt.Println("Welcome to our banking app!")

	// Start app.
	a := app{}
	a.run()

	// Flush flushes all pending log I/O.
	glog.Flush()
}
