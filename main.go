package main

import (
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"gitlab-version/runner"
	"log"
)

func main() {
	opt := &Runner.Runner{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("Get Gitlab version number automatically")
	flagSet.BoolVar(&opt.Silent, "silent", false, "show silent output")
	flagSet.IntVar(&opt.Timeout, "timeout", 20, "timeout")
	flagSet.IntVar(&opt.Concurrency, "concurrency", 100, "concurrency")
	flagSet.StringVar(&opt.Output, "output", "result.txt", "output file")
	flagSet.BoolVar(&opt.ShowCVE, "show-cve", false, "show cve")
	flagSet.BoolVar(&opt.Debug, "debug", false, "debug")
	flagSet.StringSliceVarP(&opt.UrlList, "inputs", "i", nil, "list of inputs (file,comma-separated)", goflags.FileCommaSeparatedStringSliceOptions)
	if err := flagSet.Parse(); err != nil {
		log.Fatalf("Could not parse flags: %s\n", err)
	}
	if len(opt.UrlList) == 0 {
		log.Fatalf("No input provided")
	}
	if opt.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	} else if opt.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	} else {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelInfo)
	}
	if err := opt.Run(); err != nil {
		log.Fatalf("Could not run gitlab version: %s\n", err)
	}
}
