package main

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"path"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "conf", "apps.json", "./procmain -conf apps.json")
}

func main() {
	flag.Parse()

	apps := loadApps(configFile)

	//register signal handler
	registerHandler(syscall.SIGINT, func() error {
		KillAll(apps)
		return nil
	})
	registerHandler(syscall.SIGTERM, func() error {
		KillAll(apps)
		return nil
	})

	handleSignals()

	// Run and wait
	RunAll(apps)
}

func loadApps(configFile string) []*App {
	dir, fileName := path.Split(configFile)
	if fileName == "" {
		fileName = "apps.json"
	}
	if dir == "" {
		dir = "."
	}

	log.Printf("Loading config file %s from directory '%s'.\n", fileName, dir)

	configName := fileName[0 : len(fileName)-len(path.Ext(fileName))]

	viper.SetConfigType("json")
	viper.SetConfigName(configName)
	viper.AddConfigPath(dir)

	//viper.AddConfigPath("")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error %s\n", err)
	}

	if !viper.InConfig("apps") {
		log.Fatalln("Field 'apps' is missing in config file!")
	}

	apps := make([]*App, 0)
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		log.Fatalf("'Decoding config failed: %s\n", err)
	}

	return apps
}
