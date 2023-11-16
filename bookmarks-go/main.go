package main

import (
	"flag"
	"github.com/sivaprasadreddy/bookmarks/cmd"
	"github.com/sivaprasadreddy/bookmarks/internal/config"
	"log"
)

func main() {
	cfg := getConfig()
	app := cmd.NewApp(cfg)
	app.Run()
}

func getConfig() config.AppConfig {
	var confFile string
	flag.StringVar(&confFile, "conf", "config.json", "config path, eg: -conf dev.json")
	flag.Parse()
	cfg, err := config.GetConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
