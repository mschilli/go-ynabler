package ynabler

import (
	"encoding/csv"
	"go.uber.org/zap"
	"strings"
)

type ChasePlugin Plugin

func NewChase() ChasePlugin {
	return ChasePlugin{name: "citi"}
}

func (p ChasePlugin) Applicable(lines []string) bool {
	return strings.Contains(lines[0],
		"Transaction Date,Post Date,Description,Category,Type,Amount,Memo")
}

func (p ChasePlugin) Process(log *zap.SugaredLogger, lines []string) (string, error) {
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

		// 03/02/2024,03/03/2024,AMZN Mktp US*RN7EN1EE0,Shopping,Sale,-10.85,
		inrow, err := csv.NewReader(strings.NewReader(line)).Read()

		if err != nil {
			return "", err
		}

		// outflow is marked negative, reverse
		inFlow, outFlow := toOutIn(inrow[5])

		row := []string{inrow[1], inrow[2], "", outFlow, inFlow}
		rows = append(rows, row)
	}

	return asCSV(rows)
}

func (p ChasePlugin) Name() string {
	return p.name
}
