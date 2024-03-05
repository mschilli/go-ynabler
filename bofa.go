package ynabler

import (
	"encoding/csv"
	"go.uber.org/zap"
	"strings"
)

// NOTE: Use "Microsoft Excel Format" for download, which gives you a .csv file

type BofaPlugin Plugin

func NewBofa() BofaPlugin {
	return BofaPlugin{name: "bofa"}
}

func (p BofaPlugin) Applicable(lines []string) bool {
	return strings.Contains(lines[0], "Description,,Summary Amt.")
}

func (p BofaPlugin) Process(log *zap.SugaredLogger, lines []string) (string, error) {
	rows := [][]string{}

	inCSV := false

	for _, line := range lines {
		log.Debugw("Processing", "line", line)
		if !inCSV { // skip header lines
			if strings.HasPrefix(line, "Date,Description") {
				inCSV = true
			}
			log.Debugw("Skipped (header)")
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

		// 02/14/2024,Beginning balance as of 02/14/2024,,"12,345.67"
		if strings.HasPrefix(inrow[1], "Beginning balance as of") {
			continue
		}

		// bofa shows outflows as negative, reverse
		inFlow, outFlow := toOutIn(inrow[2])

		row := []string{inrow[0], inrow[1], "", outFlow, inFlow}
		rows = append(rows, row)
	}

	return asCSV(rows)
}

func (p BofaPlugin) Name() string {
	return p.name
}
