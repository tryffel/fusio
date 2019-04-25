package main

import (
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/err"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path/filepath"
)

func main() {

	logFormat := &prefixed.TextFormatter{
		ForceColors:    true,
		FullTimestamp:  true,
		QuoteCharacter: "'",
	}
	logFormat.ForceFormatting = true
	logrus.SetFormatter(logFormat)

	parser := argparse.NewParser("fusio", "Fusio server, IoT server written in Go. For more help see --help. "+
		"For basic usage, run fusio -c <config-file>.")
	configFile := parser.String("c", "Config", &argparse.Options{Required: false, Help: "Configuration file"})
	createConfig := parser.Flag("n", "new", &argparse.Options{Required: false, Help: "Create new configuration file"})
	loadDemo := parser.Flag("d", "demo", &argparse.Options{Required: false, Help: "Load demo data"})

	migrator := NewMigrator(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		logrus.Error("Failed to parse input flags")
		os.Exit(1)
	}

	if *configFile == "" {
		*configFile = "/etc/fusio/config.yaml"
	}

	dir, file := filepath.Split(*configFile)
	conf := config.NewConfig(dir, file)

	if *loadDemo == true {
		conf.LoadDemoData()
	}

	if *createConfig == true {
		conf.CreateFile()
		os.Exit(0)
	}

	conf.LoadConfig()
	service, err := fusio.NewService(&conf)
	if err != nil {
		e := Err.Wrap(&err, "failed to start fusio server")
		Err.Log(e)
		os.Exit(1)
	}

	err = migrator.RunMigrations(service.Store.GetDb())
	if err != nil {
		Err.Log(err)
		os.Exit(1)
	}
	service.Start()

}
