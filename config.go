package main

import (
	"github.com/cespare/flagconf"
	//"log"
	"os"
)

type Config struct {
	Http   string `desc:"the IP address and port for the DocsJotter server. Defaults to 127.0.0.1:8989"`
	Root   string `desc:"which directory to serve"`
	Resdir string `desc:"path to the DocsJotter resource files. Defaults to html"`
	Mode   string `desc:"can take values production, fstest or mocktest. In production mode, Root must be speciried. In fstest, root defaults to testdata. MockTest uses a mock backend"`
}

func LoadConfig(configfile string) Config {
	config := Config{
		Http:   "127.0.0.1:8989",
		Root:   "",
		Resdir: "html",
		Mode:   "fstest",
	}
	/*
		if err := flagconf.Parse(configfile, &config); err != nil {
			log.Printf("LoadConfig: Could not find config file %s. Will use default values/flags", configfile)
		}*/
	flagconf.ParseStrings(os.Args, configfile, &config, true)

	switch config.Mode {
	case "fstest":
		config.Root = "testdata"
	}
	return config
}
