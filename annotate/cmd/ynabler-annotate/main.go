package main

import (
	"flag"
	"fmt"
	"github.com/mschilli/go-ynabler/annotate"
	"go.uber.org/zap"
	"os"
)

const Version = "0.0.1"

func main() {
	verbose := flag.Bool("verbose", false, "Verbose mode")
	version := flag.Bool("version", false, "Print release version")
	orderFile := flag.String("orders", "", "Order File")

	flag.Parse()
	flag.Usage = func() {
	    fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s --orders=orders.csv ynabler.csv\n", os.Args[0])
	    flag.PrintDefaults()
        }

	if *version {
		fmt.Println(Version, Version)
		return
	}

	log := zap.Must(zap.NewProduction())
	if *verbose {
		log = zap.NewExample()
	}

	if flag.NArg() != 1 {
	    log.Error("Provide a valid ynabler .csv file")
	    flag.Usage()
	    os.Exit(1)
	}

	orders := annotate.NewOrders()
	err := orders.ParseHistoryFile(*orderFile)
	if err != nil {
		log.Error("")
		return
	}

	for _, o := range orders.Orders {
		fmt.Printf("%s,%.2f,%.60s\n",
			o.At.Format("2006-01-02"),
			float64(o.Total)/100,
			o.Item,
		)
	}
}
