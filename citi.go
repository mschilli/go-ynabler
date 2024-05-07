package ynabler

import (
	"encoding/csv"
	"go.uber.org/zap"
	"strings"
)

type CitiPlugin Plugin

func NewCiti() CitiPlugin {
	return CitiPlugin{name: "citi"}
}

func (p CitiPlugin) Applicable(lines []string) bool {
	return strings.Contains(lines[0],
		"Status,Date,Description,Debit,Credit,Member Name")
}

func (p CitiPlugin) Process(log *zap.SugaredLogger, lines []string) (string, error) {
	rows := [][]string{}

	inCSV := false

	for _, line := range lines {
		log.Debugw("Processing", "line", line)
		if !inCSV { // skip header line
			inCSV = true
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

		outFlow := toAbsDollars(inrow[3])
		inFlow := toAbsDollars(inrow[4])

		row := []string{inrow[1], inrow[2], "", outFlow, inFlow}
		rows = append(rows, row)
	}

	return asCSV(rows)
}

func (p CitiPlugin) Name() string {
	return p.name
}
