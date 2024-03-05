package ynabler

import (
	"encoding/csv"
	"go.uber.org/zap"
	"strings"
)

// NOTE: Use "Microsoft Excel Format" for download, which gives you a .csv file

type BofaCCPlugin Plugin

func NewBofaCC() BofaCCPlugin {
	return BofaCCPlugin{name: "bofacc"}
}

func (p BofaCCPlugin) Applicable(lines []string) bool {
	return strings.Contains(lines[0], "Posted Date,Reference Number,Payee,Address,Amount")
}

func (p BofaCCPlugin) Process(log *zap.SugaredLogger, lines []string) (string, error) {
	rows := [][]string{}

	inCSV := false

	for _, line := range lines {
		log.Debugw("Processing", "line", line)
		if !inCSV { // skip header line
			log.Debugw("Skipped (header)")
			inCSV = true
			continue
		}

		if len(line) == 0 {
			log.Debugw("Skipped (empty)")
			continue
		}

		inrow, err := csv.NewReader(strings.NewReader(line)).Read()
		if err != nil {
			return "", err
		}

		// bofa shows outflows as negative, reverse
		inFlow, outFlow := toOutIn(inrow[4])

		row := []string{inrow[0], inrow[2], "", outFlow, inFlow}
		rows = append(rows, row)
	}

	return asCSV(rows)
}

func (p BofaCCPlugin) Name() string {
	return p.name
}
