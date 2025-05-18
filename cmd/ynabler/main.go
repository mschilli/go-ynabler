package main

import (
	"flag"
	"fmt"
	"github.com/mschilli/go-ynabler"
	"go.uber.org/zap"
)

const Version = "0.0.5"

func main() {
	verbose := flag.Bool("verbose", false, "Verbose mode")
	version := flag.Bool("version", false, "Print release version")
	flag.Parse()

	if *version {
		fmt.Println(Version, ynabler.Version)
		return
	}

	if flag.NArg() != 1 {
		fmt.Println("Please provide a .csv file name.")
		return
	}

	blog := zap.Must(zap.NewProduction())
	if *verbose {
		blog = zap.NewExample()
	}
	log := blog.Sugar()

	filename := flag.Arg(0)

	plugin, err := ynabler.FindPlugin(log, filename)
	if err != nil {
		panic(err)
	}

	lines, err := ynabler.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	processedContent, err := plugin.Process(log, lines)
	if err != nil {
		panic(err)
	}

	fmt.Print(processedContent)
}
