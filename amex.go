package ynabler

import (
	"encoding/csv"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

type AmexPlugin Plugin

func NewAmex() AmexPlugin {
	return AmexPlugin{name: "citi"}
}

func (p AmexPlugin) Applicable(lines []string) bool {
	return strings.Contains(lines[0],
		"Date,Description,Card Member,Account #,Amount")
}

func (p AmexPlugin) Process(log *zap.SugaredLogger, lines []string) (string, error) {
	rows := [][]string{}

	inCSV := false

	spacesRe := regexp.MustCompile(`\s+`)

	for _, line := range lines {
		log.Debugw("Processing", "line", line)
		if !inCSV { // skip header line
			inCSV = true
			log.Debugw("Skipped (header)")
			continue
		}

		if len(line) == 0 {
			continue
		}

		inrow, err := csv.NewReader(strings.NewReader(line)).Read()

		if err != nil {
			return "", err
		}

		outFlow, inFlow := toOutIn(inrow[4])

		inrow[1] = spacesRe.ReplaceAllString(inrow[1], " ")

		row := []string{inrow[0], inrow[1], "", outFlow, inFlow}
		rows = append(rows, row)
	}

	return asCSV(rows)
}

func (p AmexPlugin) Name() string {
	return p.name
}
