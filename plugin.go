package ynabler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

const Version = "0.0.3"

type PluginIf interface {
	Applicable(lines []string) bool
	Process(log *zap.SugaredLogger, lines []string) (string, error)
	Name() string
}

type Plugin struct {
	name string
}

func FindPlugin(log *zap.SugaredLogger, fname string) (PluginIf, error) {
	lines, err := ReadFile(fname)
	if err != nil {
		return nil, err
	}

	plugins := []PluginIf{NewCiti(), NewBofa(), NewBofaCC(), NewAmex(), NewChase()}

	log.Debugw("Input read", "Lines", lines)

	for _, plugin := range plugins {
		log.Debugw("Trying", "plugin", plugin.Name())
		if plugin.Applicable(lines) {
			return plugin, nil
		}
	}

	return nil, errors.New("No applicable plugin found")
}

func ReadFile(fname string) ([]string, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func asCSV(rows [][]string) (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	err := w.Write([]string{"Date", "Payee", "Memo", "Outflow", "Inflow"})
	if err != nil {
		return "", err
	}

	for _, row := range rows {
		err := w.Write(row)
		if err != nil {
			return "", err
		}
	}

	w.Flush()

	return buf.String(), nil
}

func toAbsDollars(amount string) string {
	if len(amount) == 0 {
		return ""
	}

	if strings.Contains(amount, "-") {
		return "$" + amount[1:]
	}
	return "$" + amount
}

func toOutIn(amount string) (string, string) {
	outFlow := ""
	inFlow := ""
	if strings.Contains(amount, "-") {
		inFlow = "$" + amount[1:]
	} else {
		outFlow = "$" + amount
	}
	return outFlow, inFlow
}
