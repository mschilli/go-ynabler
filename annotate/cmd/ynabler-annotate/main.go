package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/mschilli/go-ynabler/annotate"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
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

	if len(*orderFile) == 0 {
		log.Error("Provide a valid order file")
		flag.Usage()
		os.Exit(1)
	}

	orders := annotate.NewOrders()
	err := orders.ParseHistoryFile(*orderFile)
	if err != nil {
		log.Error("Can't parse", zap.String("file", *orderFile))
		return
	}

	if flag.NArg() == 0 {
		for _, o := range orders.Orders {
			fmt.Printf("%s,%.2f,%.60s\n",
				o.At.Format("2006-01-02"),
				float64(o.Total)/100,
				o.Item,
			)
		}
		return
	}

	csvFile := flag.Arg(0)

	file, err := os.Open(csvFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	log.Debug("Reading", zap.String("csv file", csvFile))
	recs, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(os.Stdout)
	writer.Write(recs[0])

	for _, rec := range recs[1:] {
		ts, err := time.Parse("01/02/2006", rec[0])
		if err != nil {
			panic(err)
		}
		if len(rec[3]) == 0 {
			log.Debug("Ignoring", zap.String("credit", rec[4]))
			continue
		}
		price, err := annotate.IntFromAmount(rec[3])
		if err != nil {
			panic(err)
		}

		o, err := orders.ExtractAt(price, ts.AddDate(0, 0, -3))
		if err == nil {
			maxLen := 40
			item := o.Item
			if len(item) > maxLen {
				item = item[:maxLen-1]
			}
			rec[1] = item + " " + rec[1]
		} else {
			log.Debug("No order found for", zap.String("transaction", strings.Join(rec, " ")))
		}

		writer.Write(rec)
		writer.Flush()
	}
}
